package random

import (
	"math/rand"
	"strings"
	"time"
)

type RandomType int8

const (
	ALL    RandomType = 0
	NUMBER RandomType = 1
	STRING RandomType = 2
)

type UppLowType int8

const (
	MIX   UppLowType = 0
	UPPER UppLowType = 1
	LOWER UppLowType = 2
)

//生成随机字符串
func GetRandomString(length int, randomType RandomType, uppLowType UppLowType) string {
	var str string

	switch randomType {
	case NUMBER:
		str = "0123456789"
	case STRING:
		str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	default:
		str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	switch uppLowType {
	case UPPER:
		str = strings.ToUpper(str)
	case LOWER:
		str = strings.ToLower(str)
	}

	return string(result)
}
