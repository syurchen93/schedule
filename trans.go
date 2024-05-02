package main

import (
	"fmt"
	"github.com/kataras/i18n"
)

func main() {
	I18n, err := i18n.New(i18n.Glob("./tgbot/translation/*/*"), "en", "ru")
	if err != nil {
		panic(err)
	}
	fmt.Println(I18n.Tr("en", "Greetings"))
	fmt.Println(I18n.Tr("ru", "Greetings"))
}
