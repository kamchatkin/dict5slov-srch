package main

import (
	"flag"
	"fmt"
	"kamchatkin.ru/wordle-hack/search"
	"kamchatkin.ru/wordle-hack/web"
	"os"
)

// Flags консольные аргументы
var Flags = struct {
	Letter0, Letter1, Letter2, Letter3, Letter4 *string
	Not0, Not1, Not2, Not3, Not4                *string
	Included, Excluded                          *string
}{}

func init() {
	Flags.Letter0 = flag.String("1", search.DefaultString, "первая буква. Пример: б")
	Flags.Not0 = flag.String("n1", search.DefaultString, "буквы не встречаются на первом месте. Пример: ук")

	Flags.Letter1 = flag.String("2", search.DefaultString, "вторая буква. Пример: у")
	Flags.Not1 = flag.String("n2", search.DefaultString, "буквы не встречаются на втором месте. Пример: бк")

	Flags.Letter2 = flag.String("3", search.DefaultString, "третья буква. Пример: к")
	Flags.Not2 = flag.String("n3", search.DefaultString, "буквы не встречаются на третьем месте. Пример: ув")

	Flags.Letter3 = flag.String("4", search.DefaultString, "четвертая буква. Пример: в")
	Flags.Not3 = flag.String("n4", search.DefaultString, "буквы не встречаются на четвертом месте. Пример: ка")

	Flags.Letter4 = flag.String("5", search.DefaultString, "пятая буква. Пример: а")
	Flags.Not4 = flag.String("n5", search.DefaultString, "буквы не встречаются на пятом месте. Пример: кв")

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
		*Flags.Not0,
		*Flags.Letter1,
		*Flags.Not1,
		*Flags.Letter2,
		*Flags.Not2,
		*Flags.Letter3,
		*Flags.Not3,
		*Flags.Letter4,
		*Flags.Not4,
		*Flags.Included,
		*Flags.Excluded)

	search.ConsoleSearch(query)
}
