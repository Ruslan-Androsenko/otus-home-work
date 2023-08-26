package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type FrequencyAnalysis struct {
	word  string
	count int
}

// Top5 Получить топ 5 высокочастотных слов.
func Top5(inputString string) []string {
	return getHighFrequencyWords(inputString, 5)
}

// Top7 Получить топ 7 высокочастотных слов.
func Top7(inputString string) []string {
	return getHighFrequencyWords(inputString, 7)
}

// Top10 Получить топ 10 высокочастотных слов.
func Top10(inputString string) []string {
	return getHighFrequencyWords(inputString, 10)
}

// Top15 Получить топ 15 высокочастотных слов.
func Top15(inputString string) []string {
	return getHighFrequencyWords(inputString, 15)
}

func getHighFrequencyWords(inputString string, count int) []string {
	words := strings.Fields(inputString)
	frequencies := make(map[string]int)

	for _, word := range words {
		frequencies[word]++
	}

	analysisItems := sortMapByFrequency(frequencies)

	return getWords(analysisItems, count)
}

// Получить необходимое количество значений высокочастотных слов.
func getWords(analysisItems []FrequencyAnalysis, count int) []string {
	words := make([]string, 0, count)

	for _, item := range analysisItems {
		words = append(words, item.word)

		if len(words) == count {
			break
		}
	}

	return words
}

// Сортировка словаря по количеству вхождений слов.
func sortMapByFrequency(frequencies map[string]int) []FrequencyAnalysis {
	analysisItems := make([]FrequencyAnalysis, 0, len(frequencies))

	for word, counter := range frequencies {
		item := FrequencyAnalysis{
			word:  word,
			count: counter,
		}
		analysisItems = append(analysisItems, item)
	}

	// Сортируем по частоте вхождения слов и если частота одинаковая, то по алфавиту
	sort.Slice(analysisItems, func(i, j int) bool {
		return analysisItems[i].count > analysisItems[j].count ||
			(analysisItems[i].count == analysisItems[j].count && analysisItems[i].word < analysisItems[j].word)
	})

	return analysisItems
}
