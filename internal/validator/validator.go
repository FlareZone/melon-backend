package validator

import (
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

var hexRegPattern = regexp.MustCompile(`^[0-9A-Fa-f]+$`)

func ValidatorRegister() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("hexString", func(fl validator.FieldLevel) bool {
			s := fl.Field().String()
			if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
				s = s[2:]
			}
			return hexRegPattern.MatchString(s)
		})

		_ = v.RegisterValidation("userExists", func(fl validator.FieldLevel) bool {
			userID := fl.Field().String()
			user := service.NewUser(components.DBEngine).FindUserByUuid(userID)
			return user.ID > 0
		})
	}
}
