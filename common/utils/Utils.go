package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"math"
	"reflect"
	"strings"
	"time"
)

// PasswordHash 密码加密
func PasswordHash(plainpwd string) string {
	//谷歌的加密包
	hash, err := bcrypt.GenerateFromPassword([]byte(plainpwd), bcrypt.DefaultCost) //加密处理
	if err != nil {
		fmt.Println(err)
	}
	encodePWD := string(hash) // 保存在数据库的密码，虽然每次生成都不同，只需保存一份即可
	return encodePWD
}

// CheckPassword 密码校验
func CheckPassword(plainpwd, cryptedpwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(cryptedpwd), []byte(plainpwd)) //验证（对比）
	return err == nil
}

// ParseTime 時間隔式處理
func ParseTime(t string) string {
	timeString, _ := time.Parse(time.RFC3339, t)
	str := strings.Split(timeString.String(), " +")
	res := str[0]
	return res
}

// ParseIntTime int時間隔式處理
func ParseIntTime(t int64) string {
	return time.Unix(t, 0).UTC().Format("2006-01-02 15:04:05")
}

// Contain 判斷obj是否在target中，target支援的型別array,slice,map
func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// FloatMulC : 增加currency 判别决定小数位数(预设4位)
func FloatMulC(s float64, p float64, currency ...string) float64 {

	if len(currency) > 0 && (strings.EqualFold(currency[0], "BTC") || strings.EqualFold(currency[0], "ETH")) {
		return FloatMul(s, p, 7)
	} else {
		return FloatMul(s, p)
	}
}

// FloatMul 浮點數乘法 (precision=4)
func FloatMul(s float64, p float64, precisions ...int32) float64 {

	f1 := decimal.NewFromFloat(s)
	f2 := decimal.NewFromFloat(p)

	var precision int32
	if len(precisions) > 0 {
		precision = precisions[0]
	} else {
		precision = 4
	}

	res, _ := f1.Mul(f2).Truncate(precision).Float64()

	return res
}

// FloatMulC : 增加currency 判别决定小数位数
func FloatDivC(s float64, p float64, currency ...string) float64 {

	if len(currency) > 0 && (strings.EqualFold(currency[0], "BTC") || strings.EqualFold(currency[0], "ETH")) {
		return FloatDiv(s, p, 7)
	} else {
		return FloatDiv(s, p)
	}
}

// FloatDiv 浮點數除法 (precision=4)
func FloatDiv(s float64, p float64, precisions ...int32) float64 {

	f1 := decimal.NewFromFloat(s)
	f2 := decimal.NewFromFloat(p)

	var precision int32
	if len(precisions) > 0 {
		precision = precisions[0]
	} else {
		precision = 4
	}
	res, _ := f1.Div(f2).Truncate(precision).Float64()

	return res
}

// FloatSubC : 增加currency 判别决定小数位数
func FloatSubC(s float64, p float64, currency ...string) float64 {

	if len(currency) > 0 && (strings.EqualFold(currency[0], "BTC") || strings.EqualFold(currency[0], "ETH")) {
		return FloatSub(s, p, 7)
	} else {
		return FloatSub(s, p)
	}
}

// FloatSub 浮點數減法 (precision=4)
func FloatSub(s float64, p float64, precisions ...int32) float64 {

	f1 := decimal.NewFromFloat(s)
	f2 := decimal.NewFromFloat(p)

	var precision int32
	if len(precisions) > 0 {
		precision = precisions[0]
	} else {
		precision = 4
	}
	res, _ := f1.Sub(f2).Truncate(precision).Float64()

	return res
}

// FloatAddC : 增加currency 判别决定小数位数
func FloatAddC(s float64, p float64, currency ...string) float64 {

	if len(currency) > 0 && (strings.EqualFold(currency[0], "BTC") || strings.EqualFold(currency[0], "ETH")) {
		return FloatAdd(s, p, 7)
	} else {
		return FloatAdd(s, p)
	}
}

// FloatAdd 浮點數加法 (precision=4)
func FloatAdd(s float64, p float64, precisions ...int32) float64 {

	f1 := decimal.NewFromFloat(s)
	f2 := decimal.NewFromFloat(p)

	var precision int32
	if len(precisions) > 0 {
		precision = precisions[0]
	} else {
		precision = 4
	}
	res, _ := f1.Add(f2).Truncate(precision).Float64()

	return res
}

func TruncateToThreeDecimal(currency string, amount float64) float64 {

	if strings.EqualFold(currency, "BTC") || strings.EqualFold(currency, "ETH") {
		return amount
	}
	// 乘以 1000 以便保留前三位小数
	pow := math.Pow(10, 3)

	// 取整，舍去第四位及之后的位数
	truncated := math.Trunc(amount*pow) / pow

	return truncated
}
