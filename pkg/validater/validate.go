package validater

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validater struct {
	vn *validator.Validate
	ut ut.Translator
}

func NewValidater() (*Validater, error) {

	enlocal := en.New()
	uni := ut.New(enlocal, enlocal)
	engtrans, _ := uni.GetTranslator("en")
	validater := validator.New()
	if err := en_translations.RegisterDefaultTranslations(validater, engtrans); err != nil {
		return nil, err
	}
	validater.RegisterValidation("phone", phoneValidation)
	validater.RegisterValidation("password", passwordValidation)
	validater.RegisterValidation("player-position", playerPositionValidator)
	validater.RegisterValidation("team-name", teamNameValidation)
	validater.RegisterValidation("player-status", playerStatusValidater)
	validater.RegisterValidation("staff-role", staffRoleValidater)
	validater.RegisterValidation("staff-desig", staffDesignationValidater)

	validater.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name != "-" {
			return name
		}

		return ""
	})

	validater.RegisterTranslation("eqfield", engtrans,
		func(ut ut.Translator) error {
			return ut.Add("eqfield", "{0} must match {1}", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("eqfield", fe.Field(), fe.Param())
			return t
		},
	)

	validater.RegisterTranslation("phone", engtrans,
		func(ut ut.Translator) error {
			return ut.Add("phone", "{0} must be valid phone number", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("phone", fe.Field())
			return t
		},
	)

	validater.RegisterTranslation("player-position", engtrans, func(ut ut.Translator) error {
		return ut.Add("player-position", "invalid player position.", true)
	},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("player-position", fe.Field())
			return t
		},
	)
	validater.RegisterTranslation("player-status", engtrans, func(ut ut.Translator) error {
		return ut.Add("player-status", "invalid player status.", true)
	},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("player-status", fe.Field())
			return t
		},
	)

	validater.RegisterTranslation("staff-role", engtrans, func(ut ut.Translator) error {
		return ut.Add("staff-role", "Invalid staff role. Enter valid staff role.", true)
	},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("staff-role", fe.Field())
			return t
		},
	)

	validater.RegisterTranslation("staff-desig", engtrans, func(ut ut.Translator) error {
		return ut.Add("staff-desig", "Invalid staff designation. Enter valid staff designation.", true)
	},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("staff-desig", fe.Field())
			return t
		},
	)

	validater.RegisterTranslation("password", engtrans,
		func(ut ut.Translator) error {
			return ut.Add("password", "password must contain uppercase, lowercase, number and special character", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("password", fe.Field())
			return t
		},
	)

	validater.RegisterTranslation("team-name", engtrans,
		func(ut ut.Translator) error {
			return ut.Add("team-name", "{0}: '{1}' is not a valid team name", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("team-name", fe.Field(), fmt.Sprintf("%v", fe.Value()))
			return t
		},
	)

	return &Validater{
		vn: validater,
		ut: engtrans,
	}, nil
}

func (v *Validater) Validation(input interface{}) validator.ValidationErrorsTranslations {

	err := v.vn.Struct(input)
	if err == nil {
		return nil
	}

	translated := err.(validator.ValidationErrors).Translate(v.ut)
	return translated
}
