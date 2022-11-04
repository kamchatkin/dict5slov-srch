package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const DictFile = "dict.txt"

const DefaultRune = 0

const DefaultString = ""

const LENGTH = 5

// RunesStruct руны поиска точного положения
type RunesStruct struct {
	Rune0, Rune1, Rune2, Rune3, Rune4 rune
	Values                            [LENGTH]rune
	Searchable                        bool
}

// RunesStructConstructor конструктор RunesStruct
func RunesStructConstructor(letter0, letter1, letter2, letter3, letter4 *string) *RunesStruct {
	r := RunesStruct{}

	strToRune := func(s string) rune {
		if s == DefaultString {
			return DefaultRune
		}

		return []rune(s)[0]
	}

	r.Rune0 = strToRune(*letter0)
	r.Rune1 = strToRune(*letter1)
	r.Rune2 = strToRune(*letter2)
	r.Rune3 = strToRune(*letter3)
	r.Rune4 = strToRune(*letter4)

	r.Values = [LENGTH]rune{
		r.Rune0,
		r.Rune1,
		r.Rune2,
		r.Rune3,
		r.Rune4}

	// будем искать отдельные буквы в словах?
	for _, rn := range r.Values {
		if rn != DefaultRune {
			r.Searchable = true
		}
	}

	return &r
}

var Runes *RunesStruct

var runesIncluded []rune
var runesExcluded []rune

func init() {
	letter0 := flag.String("1", DefaultString, "первая буква. Пример: б")
	letter1 := flag.String("2", DefaultString, "вторая буква. Пример: у")
	letter2 := flag.String("3", DefaultString, "третья буква. Пример: к")
	letter3 := flag.String("4", DefaultString, "четвертая буква. Пример: в")
	letter4 := flag.String("5", DefaultString, "пятая буква. Пример: а")

	included := flag.String("i", DefaultString, "буквы где-то должны быть. Через пробел. Пример: бук")
	excluded := flag.String("e", DefaultString, "буквы которые не должны встречаться. Через пробел. Пример: ва")

	flag.Parse()

	Runes = RunesStructConstructor(letter0, letter1, letter2, letter3, letter4)

	if *included != DefaultString {
		for _, letter := range strings.Split(*included, "") {
			runesIncluded = append(runesIncluded, []rune(letter)[0])
		}
	}

	if *excluded != DefaultString {
		for _, letter := range strings.Split(*excluded, "") {
			runesExcluded = append(runesExcluded, []rune(letter)[0])
		}
	}
}

var Words []string

type FoundedRunsRegedit struct {
	q int
	f map[rune]bool
}

func (fr *FoundedRunsRegedit) Quantity() int {
	return fr.q
}

func (fr *FoundedRunsRegedit) Clean() {
	fr.q = 0
	fr.f = map[rune]bool{}
}

func (fr *FoundedRunsRegedit) Found(r rune) {
	if _, ok := fr.f[r]; ok {
		return
	}

	fr.f[r] = true
	fr.q++
}

func main() {
	fmt.Printf("\nУсловия поиска\n1 - %s\n2 - %s\n3 - %s\n4 - %s\n5 - %s\nвключая: (%d) %s\nисключая: (%d) %s\n",
		RuneToString(Runes.Rune0),
		RuneToString(Runes.Rune1),
		RuneToString(Runes.Rune2),
		RuneToString(Runes.Rune3),
		RuneToString(Runes.Rune4),
		len(runesIncluded),
		SliceOfRunesToString(runesIncluded),
		len(runesExcluded),
		SliceOfRunesToString(runesExcluded),
	)

	dict, err := os.Open(DictFile)
	if err != nil {
		log.Panicf("не удалось прочитать файл словаря %s", DictFile)
	}

	scanner := bufio.NewScanner(dict)
	quantityRunesIncluded := len(runesIncluded)
	quantityRunesExcluded := len(runesExcluded)
	for scanner.Scan() {
		text := scanner.Text()
		word := []rune(text)
		if len(word) != LENGTH {
			log.Fatalf("В словае найдено слово состоящее не из 5 букв %s", word)
		}

		// был ли осуществляен поиск фактически
		searchState := false

		// сначала исключаем слово
		if quantityRunesExcluded > 0 {
			searchState = true
			matchExcluded := false
			for _, r := range word {
				for _, er := range runesExcluded {
					if r == er {
						matchExcluded = true
						break
					}
				}
				if matchExcluded {
					break
				}
			}
			if matchExcluded {
				continue
			}
		}

		// дальше ищем угаданное слово
		if Runes.Searchable {
			searchState = true

			// количество букв для поиска
			found := 0
			// количество соответствующих букв
			equal := 0
			for idx, r := range Runes.Values {
				if r != DefaultRune {
					found++
					if word[idx] == r {
						equal++
					}
				}
			}

			if found != equal {
				continue
			}
		}

		// проверяем наличие букв угаданных без точного положения в слове
		if quantityRunesIncluded > 0 {
			searchState = true
			FoundedRuns := FoundedRunsRegedit{}
			FoundedRuns.Clean()
			for _, r := range word {
				for _, ir := range runesIncluded {
					if r == ir {
						FoundedRuns.Found(r)
						break
					}
				}
			}
			if FoundedRuns.Quantity() != quantityRunesIncluded {
				continue
			}
		}

		if searchState {
			Words = append(Words, text)
		}
	}

	if len(Words) > 0 {
		fmt.Printf("\nПодходят слова:\n%s\n\n", strings.Join(Words, "\n"))
	} else {
		fmt.Print("\nНичего не найдено\n\n")
	}
}

func SliceOfRunesToString(rs []rune) string {
	var out []string
	for _, r := range rs {
		out = append(out, string(r))
	}

	return strings.Join(out, ", ")
}

func RuneToString(r rune) string {
	if r == DefaultRune {
		return ""
	}

	return string(r)
}
