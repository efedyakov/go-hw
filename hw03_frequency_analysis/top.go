package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func split(input string) map[string]int {
	var cache = make(map[string]int, 1)
	for _, element := range strings.Fields(input) {
		// for _, d := range []string{".", ",", "!", ":", `"`} {
		// 	element = s.ReplaceAll(element, d, "")
		// }
		cache[element]++
	}
	return cache
}

func sortMapByValue(intput map[string]int) []string {
	var keys []string = make([]string, len(intput))
	index := 0
	for key := range intput {
		keys[index] = key
		index++
	}
	sort.Slice(keys, func(i, j int) bool {
		return intput[keys[i]] > intput[keys[j]] || (intput[keys[i]] == intput[keys[j]] && keys[i] < keys[j])
	})

	return keys
}

func Top10(input string) []string {
	a := sortMapByValue(split(input))
	count := 10
	if len(a) < count {
		count = len(a)
	}
	return a[0:count]
}
