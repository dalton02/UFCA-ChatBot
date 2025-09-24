package interceptor

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	pt_translations "github.com/go-playground/validator/v10/translations/pt_BR"
)

var Validate = validator.New()
var trans ut.Translator

func InitValidator() {

	pt := pt_BR.New()
	uni := ut.New(pt, pt)
	trans, _ = uni.GetTranslator("pt_BR")
	err := pt_translations.RegisterDefaultTranslations(Validate, trans)
	if err != nil {
		panic("Erro ao registrar traduções em português: " + err.Error())
	}
}

func ValidateAndExtract[T any](ctx *gin.Context, body *T) (err error) {
	err = ctx.ShouldBindJSON(body)
	if err != nil {
		return err
	}
	err = Validate.Struct(body)
	return err
}

func ValidateAndExtractQuery[T any](ctx *gin.Context, body *T) (err error) {
	err = ctx.ShouldBindQuery(body)
	if err != nil {
		return err
	}
	err = Validate.Struct(body)
	return err
}
