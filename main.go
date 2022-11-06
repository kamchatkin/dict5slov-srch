package main

import (
	"flag"
	"fmt"
	"kamchatkin.ru/dict5slov-srch/search"
	"kamchatkin.ru/dict5slov-srch/web"
	"os"
)

// Flags консольные аргументы
var Flags = struct {
	Letter0, Letter1, Letter2, Letter3, Letter4 *string
	Included, Excluded                          *string
}{}

func init() {
	Flags.Letter0 = flag.String("1", search.DefaultString, "первая буква. Пример: б")
	Flags.Letter1 = flag.String("2", search.DefaultString, "вторая буква. Пример: у")
	Flags.Letter2 = flag.String("3", search.DefaultString, "третья буква. Пример: к")
	Flags.Letter3 = flag.String("4", search.DefaultString, "четвертая буква. Пример: в")
	Flags.Letter4 = flag.String("5", search.DefaultString, "пятая буква. Пример: а")

	Flags.Included = flag.String("i", search.DefaultString, "буквы где-то должны быть. Через пробел. Пример: бук")
	Flags.Excluded = flag.String("e", search.DefaultString, "буквы которые не должны встречаться. Через пробел. Пример: ва")

	flag.Parse()
}

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "web" {
			web.Web()
			return
		}

		fmt.Println("Не опознанный аргумент.\nИспользуйте web для запуска веб-сервера\nи без аргументов для использования в консоли.")
		os.Exit(0)
	}

	query := search.QueryConstructor(
		*Flags.Letter0,
		*Flags.Letter1,
		*Flags.Letter2,
		*Flags.Letter3,
		*Flags.Letter4,
		*Flags.Included,
		*Flags.Excluded)

	search.ConsoleSearch(query)
}
