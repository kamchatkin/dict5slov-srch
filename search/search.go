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
	// Single зеленые буквы
	Single [LENGTH]rune
	// SingleSearchable искать по зеленым буквам
	SingleSearchable bool

	// Not0 желтые буквы, не первая буква в слове
	Not0 []rune
	// Not0Searchable искать по первой букве?
	Not0Searchable bool

	// Not1 желтые буквы, не вторая буква в слове
	Not1 []rune
	// Not1Searchable искать по второй букве?
	Not1Searchable bool

	// Not2 желтые буквы, не третья буква в слове
	Not2 []rune
	// Not2Searchable искать по третьей букве?
	Not2Searchable bool

	// Not3 желтые буквы, не четвертая буква в слове
	Not3 []rune
	// Not3Searchable искать по четвертой букве?
	Not3Searchable bool

	// Not4 желтые буквы, не пятая буква в слове
	Not4 []rune
	// Not4Searchable искать по пятой букве?
	Not4Searchable bool

	// Include буквы которые встречаются в слове, но неизвестно точное положение
	Include []rune
	// IncludeSearchable искать по встречающимся буквам?
	IncludeSearchable bool

	// Exclude буквы которые точно не встречаются в слове
	Exclude []rune
	// ExcludeSearchable искать по буквам которые точно не встречаются?
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

var str2runes = func(str string) []rune {
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

func QueryConstructor(s0, not0, s1, not1, s2, not2, s3, not3, s4, not4, inc, exc string) *query {
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

	q.Not0 = str2runes(not0)
	q.Not0Searchable = len(q.Not0) > 0

	q.Not1 = str2runes(not1)
	q.Not1Searchable = len(q.Not1) > 0

	q.Not2 = str2runes(not2)
	q.Not2Searchable = len(q.Not2) > 0

	q.Not3 = str2runes(not3)
	q.Not3Searchable = len(q.Not3) > 0

	q.Not4 = str2runes(not4)
	q.Not4Searchable = len(q.Not4) > 0

	q.Include = str2runes(inc)
	q.IncludeSearchable = len(q.Include) > 0

	q.Exclude = str2runes(exc)
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
		fmt.Printf("\nПодходят слова (%d):\n%s\n\n", len(*words), strings.Join(*words, "\n"))
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

		if q.Not0Searchable {
			match := false
			for _, r := range q.Not0 {
				if r == word[0] {
					match = true
				}
			}
			if match {
				continue
			}
		}

		if q.Not1Searchable {
			match := false
			for _, r := range q.Not1 {
				if r == word[1] {
					match = true
				}
			}
			if match {
				continue
			}
		}

		if q.Not2Searchable {
			match := false
			for _, r := range q.Not2 {
				if r == word[2] {
					match = true
				}
			}
			if match {
				continue
			}
		}

		if q.Not3Searchable {
			match := false
			for _, r := range q.Not3 {
				if r == word[3] {
					match = true
				}
			}
			if match {
				continue
			}
		}

		if q.Not4Searchable {
			match := false
			for _, r := range q.Not4 {
				if r == word[4] {
					match = true
				}
			}
			if match {
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
