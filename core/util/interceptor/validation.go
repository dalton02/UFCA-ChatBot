package interceptor

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()
var messagesErrors = map[string]string{
	"required":  "deve ser um campo obrigatório",
	"email":     "deve ser um endereço de email válido",
	"min":       "é um valor muito curto",
	"max":       "é um valor muito longo",
	"numeric":   "deve ser um número válido",
	"alphanum":  "deve conter apenas letras e números",
	"alpha":     "deve conter apenas letras",
	"url":       "deve ser uma URL válida",
	"uuid":      "deve ser uma UUID válido",
	"uuid4":     "deve ser uma UUID v4 válido",
	"datetime":  "deve ser uma data no formato correto",
	"eqfield":   "deve ser igual ao campo de referência",
	"gt":        "deve ser maior que o valor especificado",
	"gte":       "deve ser maior ou igual ao valor especificado",
	"lt":        "deve ser menor que o valor especificado",
	"lte":       "deve ser menor ou igual ao valor especificado",
	"unique":    "deve conter valores únicos",
	"hexcolor":  "deve ser uma cor hexadecimal válida",
	"rgb":       "deve ser uma cor RGB válida",
	"lowercase": "deve conter apenas letras minúsculas",
	"uppercase": "deve conter apenas letras maiúsculas",
}

func ValidateAndExtract[T any](ctx *gin.Context, body *T) (err error) {

	err = ctx.ShouldBindJSON(body)

	if err != nil {
		return parseError(err)
	}
	err = validate.Struct(body)
	if err != nil {
		return parseError(err)
	}
	return nil
}

func ValidateAndExtractQuery[T any](ctx *gin.Context, body *T) (err error) {

	err = ctx.ShouldBindQuery(body)
	if err != nil {
		return err
	}
	err = validate.Struct(body)
	if err != nil {
		return parseError(err)
	}
	return nil
}

func RegisterValidation(name string, fn func(fl validator.FieldLevel) bool, message string) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation(name, fn)
	}
	validate.RegisterValidation(name, fn)

	messagesErrors[name] = message
}

func getKeys(message string) (nestedKey string, key string, tag string, err error) {
	re := regexp.MustCompile(`'([^']*)'|"([^"]*)"`)
	matches := re.FindAllStringSubmatch(message, -1)

	var results []string
	for _, match := range matches {
		if match[1] != "" {
			results = append(results, match[1])
		} else if match[2] != "" {
			results = append(results, match[2])
		}
	}
	if len(results) < 3 {
		return nestedKey, key, tag, fmt.Errorf(message)
	}
	nestedKey = results[0]
	key = results[1]
	tag = results[2]
	return nestedKey, key, tag, nil
}

func parseError(err error) error {
	_, key, tag, errKeys := getKeys(err.Error())
	if errKeys != nil {
		return err
	}
	msg, custom := messagesErrors[tag]

	if custom {
		return fmt.Errorf("erro de validação em '%s': %s", key, msg)
	}

	return fmt.Errorf("erro de validação em '%s': '%s'", key, tag)

}
