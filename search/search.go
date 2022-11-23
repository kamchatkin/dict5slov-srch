package search

import (
	_ "embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
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

	// желтые буквы, не N буква в слове
	Not [][]rune

	// Include буквы которые встречаются в слове, но неизвестно точное положение
	Include []rune
	// IncludeSearchable искать по встречающимся буквам?
	IncludeSearchable bool

	// Exclude буквы которые точно не встречаются в слове
	Exclude []rune
	// ExcludeSearchable искать по буквам которые точно не встречаются?
	ExcludeSearchable bool

	// Sort сортировка
	Sort sortMethod
}

// sortMethod Типы
type sortMethod string

const SortAlphabet = sortMethod("alphabet")
const SortWeight = sortMethod("weight")

var str2rune = func(s string) rune {
	if s == DefaultString {
		return DefaultRune
	}

	return []rune(s)[0]
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

// QueryConstructor собираем параметры поиска в единый объект
func QueryConstructor(s0, s1, s2, s3, s4, not0, not1, not2, not3, not4, exc string) (*query, error) {
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

	q.Not = append(q.Not, str2runes(not0))
	q.Not = append(q.Not, str2runes(not1))
	q.Not = append(q.Not, str2runes(not2))
	q.Not = append(q.Not, str2runes(not3))
	q.Not = append(q.Not, str2runes(not4))

	// добавляем к q.Include уникальный список букв из q.Not
	includeMap := map[rune]bool{}
	for _, runesSlice := range q.Not {
		for _, r := range runesSlice {
			if _, ok := includeMap[r]; !ok {
				includeMap[r] = true
				q.Include = append(q.Include, r)
			}
		}
	}
	if len(q.Include) > 5 {
		return &query{}, errors.New("В поиске по угаданным буквам с неточным положением не может быть более 5 разных букв")
	}
	q.IncludeSearchable = len(q.Include) > 0

	q.Exclude = str2runes(exc)
	q.ExcludeSearchable = len(q.Exclude) > 0

	// поиск противоречащих параметров поиска
	if q.ExcludeSearchable && (q.IncludeSearchable || q.SingleSearchable) {
		var incErr []rune
		var singleErr []rune
		for _, re := range q.Exclude {
			for _, ri := range q.Include {
				if re == ri {
					incErr = append(incErr, re)
				}
			}
			for _, rs := range q.Single {
				if re == rs {
					singleErr = append(singleErr, re)
				}
			}
		}

		errText := ""
		if len(incErr) > 0 {
			errText = fmt.Sprintf("\n\tв исключении и в поиске по неточному положению встречаеися буква: %s", string(incErr))
		}

		if len(singleErr) > 0 {
			errText = fmt.Sprintf("\n\tв исключении и в поиске с точным положением встречаеися буква: %s", string(singleErr))
		}

		if errText != "" {
			return &query{}, errors.New(fmt.Sprintf("%s\n\n", errText))
		}
	}

	return &q, nil
}

// DefaultRune нулевая руна, для пропуска поиска
const DefaultRune = 0

// DefaultString пустая строка для игнорирования в поиске
const DefaultString = ""

// WebSearch результаты поиска в json
func WebSearch(q *query) *[]string {
	return search(q)
}

// ConsoleSearch результаты поиска в терминал
func ConsoleSearch(q *query) {
	var tpl string

	if q.SingleSearchable {
		for idx, r := range q.Single {
			if r != DefaultRune {
				tpl = fmt.Sprintf("%d буква: %s\n", idx, string(r))
			}
		}
	}

	for idx, runesSlice := range q.Not {
		length := len(runesSlice)
		if length == 1 {
			tpl += fmt.Sprintf("%d буква не: (%d) %s\n", idx+1, length, string(runesSlice))
		} else if length > 1 {
			tpl += fmt.Sprintf("%d буква не одна из: (%d) %s\n", idx+1, length, sliceOfRunesToString(runesSlice))
		}
	}

	if q.IncludeSearchable {
		tpl += fmt.Sprintf("Включая буквы: (%d) %s\n", len(q.Include), sliceOfRunesToString(q.Include))
	}
	if q.ExcludeSearchable {
		tpl += fmt.Sprintf("Исключая буквы: (%d) %s\n\n", len(q.Exclude), sliceOfRunesToString(q.Exclude))
	}

	fmt.Println(tpl)

	timeStart := time.Now()
	words := search(q)

	if len(*words) < 1 {
		fmt.Printf("\nНичего не найдено за %s\n\n", time.Since(timeStart))
		return
	}

	fmt.Printf("\nПодходят слова (%d):\n%s\n\nПотребовалось на поиск: %s\n\n",
		len(*words),
		strings.Join(*words, "\n"),
		time.Since(timeStart))
}

type FoundedRunsRegedit struct {
	quantity int
	found    map[rune]bool
}

func (fr *FoundedRunsRegedit) Found(r rune) {
	if _, ok := fr.found[r]; ok {
		return
	}

	fr.found[r] = true
	fr.quantity++
}

// WordOfDict слово и его вес
type WordOfDict struct {
	Word   string
	Weight int
}

// search поиск
func search(q *query) *[]string {
	var words []WordOfDict
	//var words []string

	csvReader := csv.NewReader(strings.NewReader(dict))
	csvReader.Comma = ';'
	//scanner := bufio.NewScanner(strings.NewReader(dict))
	quantityRunesIncluded := len(q.Include)
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Sprintf("Ошибка чтения словаря: %s", err))
		}

		text := row[0]
		word := []rune(text)

		weight, _ := strconv.Atoi(row[1])

		//fmt.Println(weight)

		//text := scanner.Text()
		//word := []rune(text)
		if len(word) != LENGTH {
			continue
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

		// поиск отрицаний. N буква не один из
		match := false
		for runeNum, runes := range q.Not {
			if len(runes) < 1 {
				continue
			}

			searchState = true
			for _, r := range runes {
				if r == word[runeNum] {
					match = true
					break
				}
			}

			if match {
				break
			}
		}
		if match {
			continue
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
			FoundedRuns := FoundedRunsRegedit{quantity: 0, found: make(map[rune]bool)}
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
			words = append(words, WordOfDict{Weight: weight, Word: text})
		}
	}

	if q.Sort == SortAlphabet {
		return sliceOfDictToSliceOfString(words)
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].Weight > words[j].Weight
	})

	// сортировка по весу
	return sliceOfDictToSliceOfString(words)
}

// sliceOfDictToSliceOfString
func sliceOfDictToSliceOfString(words []WordOfDict) *[]string {
	var out []string
	for _, w := range words {
		out = append(out, w.Word)
	}

	return &out
}

// sliceOfRunesToString
func sliceOfRunesToString(runes []rune) string {
	var out []string
	for _, r := range runes {
		out = append(out, string(r))
	}

	return strings.Join(out, ", ")
}
