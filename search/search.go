package search

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"strings"
)

//go:embed dict.txt
var dict string

// LENGTH какой длины ищем слово
const LENGTH = 5

// IncExcMaxLength максимальное количество рун в Include и Exclude
const IncExcMaxLength = 25

// query параметры запроса для поиска
type query struct {
	Single           [LENGTH]rune
	SingleSearchable bool

	Include           []rune
	IncludeSearchable bool

	Exclude           []rune
	ExcludeSearchable bool
}

var str2rune = func(s string) rune {
	if s == DefaultString {
		return DefaultRune
	}

	return []rune(s)[0]
}

var rune2str = func(r rune) string {
	if r == DefaultRune {
		return DefaultString
	}

	return string(r)
}

var IncEnc2Runes = func(str string) []rune {
	if str == DefaultString {
		return []rune{}
	}

	var result []rune
	iterate := 0
	for _, r := range []rune(str) {
		if iterate >= IncExcMaxLength {
			break
		}
		iterate++

		result = append(result, r)
	}

	return result
}

func QueryConstructor(s0, s1, s2, s3, s4, inc, exc string) *query {
	q := query{
		Single: [LENGTH]rune{
			str2rune(s0),
			str2rune(s1),
			str2rune(s2),
			str2rune(s3),
			str2rune(s4),
		},
	}

	for _, rn := range q.Single {
		if rn != DefaultRune {
			q.SingleSearchable = true
			break
		}
	}

	q.Include = IncEnc2Runes(inc)
	q.IncludeSearchable = len(q.Include) > 0

	q.Exclude = IncEnc2Runes(exc)
	q.ExcludeSearchable = len(q.Exclude) > 0

	return &q
}

const DefaultRune = 0

const DefaultString = ""

func WebSearch(q *query) *[]string {
	return search(q)
}

func ConsoleSearch(q *query) {
	fmt.Printf("\nУсловия поиска\n1 - %s\n2 - %s\n3 - %s\n4 - %s\n5 - %s\nвключая: (%d) %s\nисключая: (%d) %s\n",
		rune2str(q.Single[0]),
		rune2str(q.Single[1]),
		rune2str(q.Single[2]),
		rune2str(q.Single[3]),
		rune2str(q.Single[4]),
		len(q.Include),
		sliceOfRunesToString(q.Include),
		len(q.Exclude),
		sliceOfRunesToString(q.Exclude),
	)

	words := search(q)

	if len(*words) > 0 {
		fmt.Printf("\nПодходят слова:\n%s\n\n", strings.Join(*words, "\n"))
	} else {
		fmt.Print("\nНичего не найдено\n\n")
	}

}

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

func search(q *query) *[]string {
	var words []string

	//timeStart := time.Now()

	scanner := bufio.NewScanner(strings.NewReader(dict))
	quantityRunesIncluded := len(q.Include)
	for scanner.Scan() {
		text := scanner.Text()
		word := []rune(text)
		if len(word) != LENGTH {
			log.Fatalf("В словаре найдено слово состоящее не из %d букв %s", LENGTH, text)
		}

		// был ли осуществлен поиск фактически
		searchState := false

		// сначала исключаем слово
		if q.ExcludeSearchable {
			searchState = true
			matchExcluded := false
			for _, r := range word {
				for _, er := range q.Exclude {
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

		// дальше ищем слова с угаданными буквами
		if q.SingleSearchable {
			searchState = true

			// количество букв для поиска
			found := 0
			// количество соответствующих букв
			equal := 0
			for idx, r := range q.Single {
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
		if q.IncludeSearchable {
			searchState = true
			FoundedRuns := FoundedRunsRegedit{}
			FoundedRuns.Clean()
			for _, r := range word {
				for _, ir := range q.Include {
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
			words = append(words, text)
		}
	}

	return &words

	//elapsed := time.Since(timeStart)
	//log.Printf("Binomial took %s", elapsed)
}

func sliceOfRunesToString(rs []rune) string {
	var out []string
	for _, r := range rs {
		out = append(out, string(r))
	}

	return strings.Join(out, ", ")
}
