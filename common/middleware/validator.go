package middleware

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

var hexRegPattern = regexp.MustCompile(`^[0-9A-Fa-f]+$`)

func GinValidatorRegister() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("hexString", func(fl validator.FieldLevel) bool {
			s := fl.Field().String()
			if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
				s = s[2:]
			}
			return hexRegPattern.MatchString(s)
		})
	}
}
