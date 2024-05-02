package util

import (
	"fmt"

	"github.com/kataras/i18n"
)

var translator *i18n.I18n

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
