package config

import "strings"

func filterJson(content string) (newContent string) {
	lines := strings.Split(content, "\n")

	var result []string
	for _, v := range lines {
		if !strings.Contains(v, "// ") {
			result = append(result, v)
		}
	}

	return strings.Join(result, "\n")
}
