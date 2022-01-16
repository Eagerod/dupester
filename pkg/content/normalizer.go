package content

import (
	"regexp"
	"strings"
)

func NormalizeContents(arr []string) string {
	re1 := regexp.MustCompile("[\\.$-/:-?{-~!\"^_`\\[\\]\t]")
	re2 := regexp.MustCompile(" +")

	var sb strings.Builder
	for _, s := range arr {
		s = re1.ReplaceAllString(s, " ")

		s = strings.TrimSpace(s)

		s = re2.ReplaceAllString(s, " ")

		if len(s) != 0 {
			sb.WriteString(s)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
