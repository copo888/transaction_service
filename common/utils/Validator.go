package utils

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

var (
	MyValidator *validator.Validate
)

func init() {
	MyValidator = validator.New()
	// 自定义验证方法
	MyValidator.RegisterValidation("alphanumLength", checkAlphanumLength)
	MyValidator.RegisterValidation("length", checkLength)
	MyValidator.RegisterValidation("prec", checkPrecision)
	MyValidator.RegisterValidation("floatString", checkFloatString)
	// 翻譯
	//en := en.New()
	//uni = ut.New(en, en)
	//trans, _ := uni.GetTranslator("en")
}

func translateOverride(trans ut.Translator) {
	MyValidator.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})
}

func checkPrecision(fl validator.FieldLevel) bool {
	field := fmt.Sprintf("%v", fl.Field().Float())
	param := fl.Param()
	re := fmt.Sprintf("^\\d{1,}\\.?\\d{0,%s}$", param)
	r := regexp.MustCompile(re)
	return r.MatchString(field)
}

func checkAlphanumLength(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if len(field) == 0 {
		return true
	}
	param := strings.Split(fl.Param(), "/")
	re := fmt.Sprintf("^[a-zA-Z0-9_-]{%s,%s}$", param[0], param[1])
	r := regexp.MustCompile(re)
	return r.MatchString(field)
}

func checkLength(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	param := strings.Split(fl.Param(), "/")
	re := fmt.Sprintf("\\w{%s,%s}", param[0], param[1])
	r := regexp.MustCompile(re)
	return r.MatchString(field)
}

func checkFloatString(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	param := fl.Param()
	re := fmt.Sprintf("^(([1-9][0-9]*)|(([0]\\.\\d{0,%s}|[1-9][0-9]*\\.\\d{0,%s})))$", param, param)
	r := regexp.MustCompile(re)
	return r.MatchString(field)
}
