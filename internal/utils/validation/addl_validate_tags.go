package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func init() {
	validate.RegisterValidation("alpha_space_dot", isAlphaSpaceDot)
	validate.RegisterValidation("contains_alphabet", containsAlphabet)
}

func isAlphaSpaceDot(fl validator.FieldLevel) bool {
	isValid := regexp.MustCompile(`^[a-zA-Z\s.]+$`).MatchString
	return isValid(fl.Field().String())
}

func containsAlphabet(fl validator.FieldLevel) bool {
	hasAlphabet := regexp.MustCompile(`[A-Za-z]`).MatchString
	return hasAlphabet(fl.Field().String())
}
