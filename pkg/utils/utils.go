package utils

import (
	"bytes"
	"encoding/gob"
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

// DeepCopy performs a deep copy of an object using gob encoding/decoding.
func DeepCopy(src, dst interface{}) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	dec := gob.NewDecoder(&b)

	if err := enc.Encode(src); err != nil {
		return err
	}
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

func IsMobile(userAgent string) bool {
	return strings.Contains(userAgent, "Mobile") ||
		strings.Contains(userAgent, "Android") ||
		strings.Contains(userAgent, "iPhone") ||
		strings.Contains(userAgent, "iPad")
}
