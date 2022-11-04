package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const DictFile = "dict.txt"

const DefaultRune = 0

const DefaultString = ""

const LENGTH = 5

// RunesStruct руны поиска точного положения
type RunesStruct struct {
	Values     [LENGTH]rune
	Searchable bool
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

	r.Values = [LENGTH]rune{
		strToRune(*letter0),
		strToRune(*letter1),
		strToRune(*letter2),
		strToRune(*letter3),
		strToRune(*letter4)}

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
	quantity int
	found    map[rune]bool
}

func (fr *FoundedRunsRegedit) Clean() {
	fr.quantity = 0
	fr.found = map[rune]bool{}
}

func (fr *FoundedRunsRegedit) Found(r rune) {
	if _, ok := fr.found[r]; ok {
		return
	}

	fr.found[r] = true
	fr.quantity++
}

func main() {
	fmt.Printf("\nУсловия поиска\n1 - %s\n2 - %s\n3 - %s\n4 - %s\n5 - %s\nвключая: (%d) %s\nисключая: (%d) %s\n",
		RuneToString(Runes.Values[0]),
		RuneToString(Runes.Values[1]),
		RuneToString(Runes.Values[2]),
		RuneToString(Runes.Values[3]),
		RuneToString(Runes.Values[4]),
		len(runesIncluded),
		SliceOfRunesToString(runesIncluded),
		len(runesExcluded),
		SliceOfRunesToString(runesExcluded),
	)

	dict, err := os.Open(DictFile)
	if err != nil {
		log.Panicf("не удалось прочитать файл словаря %s", DictFile)
	}

	timeStart := time.Now()

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
			if FoundedRuns.quantity != quantityRunesIncluded {
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

	elapsed := time.Since(timeStart)
	log.Printf("Binomial took %s", elapsed)
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
