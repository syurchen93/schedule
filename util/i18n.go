package util

import (
	"fmt"

	"github.com/kataras/i18n"
)

type TranslatorInterface interface {
	Tr(locale, key string, args ...interface{}) string
}

var translator TranslatorInterface

func InitTranslator(path string, supportedLocales []string) {
	var err error
	translator, err = i18n.New(i18n.Glob(fmt.Sprintf("./%s/*/*", path)), supportedLocales...)
	if err != nil {
		panic(err)
	}
}

func Translate(locale, key string) string {
	return translator.Tr(locale, key)
}

func SetTranslator(t TranslatorInterface) {
	translator = t
}
