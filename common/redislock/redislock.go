package redislock

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	red "github.com/redis/go-redis/v9"
)

const (
	lockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`

	letters         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomLen       = 16
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
)

// A RedisLock is a redis lock.
type RedisLock struct {
	redis   *red.Client
	seconds uint32
	key     string
	id      string
}

var tempContext = context.Background()

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewRedisLock returns a RedisLock.
func New(redis *red.Client, key string, prefix string) *RedisLock {
	return &RedisLock{
		redis:   redis,
		seconds: 3,
		key:     prefix + key,
		id:      randomStr(randomLen),
	}
}

// Acquire acquires the lock.
func (rl *RedisLock) Acquire() (bool, error) {
	seconds := atomic.LoadUint32(&rl.seconds)
	resp, err := rl.redis.Eval(tempContext, lockCommand, []string{rl.key}, []string{
		rl.id, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
	}).Result()

	if err == red.Nil {
		return false, nil
	} else if err != nil {
		_ = fmt.Errorf("error on acquiring lock for %s, %s", rl.key, err.Error())
		return false, err
	} else if resp == nil {
		return false, nil
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return true, nil
	}

	_ = fmt.Errorf("unknown reply when acquiring lock for %s: %v", rl.key, resp)
	return false, nil
}

func (rl *RedisLock) TryLockTimeout(timeOutSeconds float64) (bool, error) {
	startTime := time.Now()
	for {
		if elapseTime := time.Since(startTime).Seconds(); elapseTime < timeOutSeconds {
			if ok, err := rl.Acquire(); !ok || err != nil {
				fmt.Printf("key:%s, id:%s Locked, retry %03f\n", rl.key, rl.id, elapseTime)
			} else {
				return true, nil
			}
		} else {
			break
		}
		time.Sleep(70 * time.Millisecond)
	}
	return false, errors.New(fmt.Sprintf("Cann't acquiring lock within %03fs", timeOutSeconds))
}

// Release releases the lock.
func (rl *RedisLock) Release() (bool, error) {
	resp, err := rl.redis.Eval(tempContext, delCommand, []string{rl.key}, []string{rl.id}).Result()
	if err != nil {
		return false, err
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}

	return reply == 1, nil
}

// SetExpire sets the expire.
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}

func randomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
