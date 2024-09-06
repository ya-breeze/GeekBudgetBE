package utils

import (
	"fmt"
	"strings"
)

func StrToRef(s string) *string {
	return &s
}

//nolint:forbidigo // it's okay to use fmt in this function
func PrintInTwoColumns(str1, str2 string) {
	lines1 := strings.Split(str1, "\n")
	lines2 := strings.Split(str2, "\n")

	maxLines := len(lines1)
	if len(lines2) > maxLines {
		maxLines = len(lines2)
	}

	for i := 0; i < maxLines; i++ {
		var line1, line2 string
		if i < len(lines1) {
			line1 = lines1[i]
		}
		if i < len(lines2) {
			line2 = lines2[i]
		}
		fmt.Printf("%-60s | %s\n", line1, line2)
	}
}
