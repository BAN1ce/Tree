package util

import "strings"

func SplitTopic(topic string) []string {
	tmp := strings.Split(strings.Trim(topic, "/"), "/")
	for _, v := range tmp {
		if v == "" {
			result := make([]string, 0)
			for _, v := range tmp {
				if v != "" {
					result = append(result, v)
				}
			}
			return result
		}
	}
	return tmp
}

func HasWildcard(topic string) bool {
	return strings.Contains(topic, "+") || strings.Contains(topic, "#")
}

func subQosMoreThan0(topics map[string]int32) bool {
	for _, v := range topics {
		if v > 0 {
			return true
		}
	}
	return false
}

func GenerateTopics(levels int, base []string, wildcards []string) []string {
	var result []string
	var generate func(current []string, depth int, hasHash bool)

	generate = func(current []string, depth int, hasHash bool) {
		if depth == levels {
			result = append(result, "/"+strings.Join(current, "/"))
			return
		}
		for _, b := range base {
			generate(append(current, b), depth+1, hasHash)
		}
		for _, w := range wildcards {
			if w == "#" {
				if !hasHash && depth == levels-1 {
					generate(append(current, w), depth+1, true)
				}
			} else {
				generate(append(current, w), depth+1, hasHash)
			}
		}
	}

	generate([]string{}, 0, false)
	return result
}
