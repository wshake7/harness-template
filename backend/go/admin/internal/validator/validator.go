package validator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func init() {
	validate.SetTagName("binding")
	_ = validate.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		field := fl.Field()
		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				return false
			}
			field = field.Elem()
		}
		if field.Kind() != reflect.String {
			return false
		}
		return strings.TrimSpace(field.String()) != ""
	})
}

func Struct[T any](t T) error {
	return ParseValidateErr(t, validate.Struct(t))
}

type StructValidator struct{}

func (v *StructValidator) Validate(out any) error {
	return ParseValidateErr(out, validate.Struct(out))
}

type Validator interface {
	Validate() error
}

func ParseValidateErr(req any, err error) error {
	if err == nil {
		return nil
	}

	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return err
	}

	t := reflect.TypeOf(req)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	for _, e := range validationErrs {
		field, ok := t.FieldByName(e.Field())
		if !ok {
			continue
		}
		msgTag, ok := field.Tag.Lookup("binding_msg")
		if !ok {
			continue
		}
		for item := range strings.SplitSeq(msgTag, ",") {
			parts := strings.SplitN(item, "=", 2)
			if len(parts) == 2 && parts[0] == e.Tag() {
				return errors.New(parts[1])
			}
		}
	}
	return err
}
