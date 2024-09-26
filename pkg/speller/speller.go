package speller

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Result struct {
	Word       string
	Suggestion string
}

func decodeUnicode(s string) (string, error) {
	re := regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	decodedStr := re.ReplaceAllStringFunc(s, func(match string) string {
		// Получаем код символа без префикса '\u'
		code := match[2:]
		val, err := strconv.ParseInt(code, 16, 32)
		if err != nil {
			return match
			// Возвращаем оригинал в случае ошибки
		}
		return string(rune(val))
	})
	return decodedStr, nil
}

func correctError(s string, text []string) ([]string, error) {
	var strs []string

	results := []Result{}
	var currentWord string

	strs = strings.Split(s, ",")

	for _, str := range strs {
		if strings.Contains(str, `"word":`) {
			currentWord = str[8 : len(str)-1]
		}

		if strings.Contains(str, `"s":`) {
			results = append(results, Result{Word: currentWord, Suggestion: str[6 : len(str)-1]})
		}
	}

	for _, res := range results {
		for i, word := range text {
			if word == res.Word {
				text[i] = res.Suggestion
				break
			}
		}
	}
	return text, nil
}

func CheckText(text []string) ([]string, error) {
	response, err := http.PostForm("https://speller.yandex.net/services/spellservice.json/checkTexts", url.Values{
		"text": text,
	})
	if err != nil {
		log.Fatalf(err.Error())
		return text, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf(err.Error())
		return text, err
	}

	formatted, err := decodeUnicode(string(body))
	if err != nil {
		log.Fatalf(err.Error())
		return text, err
	}

	fixed, err := correctError(formatted, text)
	if err != nil {
		log.Fatalf(err.Error())
		return text, err
	}

	return fixed, nil
}
