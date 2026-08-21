package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	engine *validator.Validate
}

func New() *Validator {
	return &Validator{engine: validator.New()}
}

func (v *Validator) Struct(s any) error {
	if err := v.engine.Struct(s); err != nil {
		return formatErrors(err)
	}
	return nil
}

func (v *Validator) Var(field any, tag string) error {
	if err := v.engine.Var(field, tag); err != nil {
		return formatErrors(err)
	}
	return nil
}

func (v *Validator) Engine() *validator.Validate {
	return v.engine
}

func formatErrors(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	messages := make([]string, 0, len(validationErrors))
	for _, fe := range validationErrors {
		messages = append(messages, fmt.Sprintf("%s failed on '%s'", fe.Field(), fe.Tag()))
	}
	return fmt.Errorf("%s", strings.Join(messages, "; "))
}
