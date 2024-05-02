package util

import (
	"fmt"

	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func InitLocale(localeDir string, locales []string) {
	if bundle == nil {
		bundle := i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
		for _, locale := range locales {
			bundle.MustLoadMessageFile(fmt.Sprintf("%s/active.%s.json", localeDir, locale))
		}
	}
}

func GetLocalizer(locale string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, locale)
}
