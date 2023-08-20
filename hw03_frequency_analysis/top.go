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
	analysisItems := make([]FrequencyAnalysis, 0, len(words))

	for _, word := range words {
		itemIndex := searchItem(word, analysisItems)

		if itemIndex != -1 {
			analysisItems[itemIndex].count++
		} else {
			item := FrequencyAnalysis{
				word:  word,
				count: 1,
			}
			analysisItems = append(analysisItems, item)
		}
	}

	// Сортируем по частоте вхождения слов
	sort.Slice(analysisItems, func(i, j int) bool {
		return analysisItems[i].count > analysisItems[j].count
	})
	analysisItems = sortByGroups(analysisItems)

	return getWords(analysisItems, count)
}

// Поиск индекса элемента в срезе.
func searchItem(word string, analysisItems []FrequencyAnalysis) int {
	itemIndex := -1

	for index, item := range analysisItems {
		if strings.Compare(word, item.word) == 0 {
			itemIndex = index
			break
		}
	}

	return itemIndex
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

// Сформировать группы по ключам количества вхождений.
func makeGroupOfKeys(analysisSources []FrequencyAnalysis) (map[int][]FrequencyAnalysis, []int) {
	groups := make(map[int][]FrequencyAnalysis)
	var keys []int

	for _, item := range analysisSources {
		if _, ok := groups[item.count]; !ok {
			keys = append(keys, item.count)
		}

		groups[item.count] = append(groups[item.count], item)
	}

	return groups, keys
}

// Сортировка по группам с одинаковым количеством вхождений.
func sortByGroups(analysisSources []FrequencyAnalysis) []FrequencyAnalysis {
	analysisDestinations := make([]FrequencyAnalysis, 0, len(analysisSources))
	groups, keys := makeGroupOfKeys(analysisSources)

	for _, key := range keys {
		if items, ok := groups[key]; ok {
			sort.Slice(items, func(i, j int) bool {
				return items[i].word < items[j].word
			})
			analysisDestinations = append(analysisDestinations, items...)
		}
	}

	return analysisDestinations
}
