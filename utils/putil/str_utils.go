package putil

import "unicode"

func F(str string, maxLen int) string {
	lengthNow := 0
	for _, c := range str {
		if unicode.Is(unicode.Han, c) {
			lengthNow += 2
		} else {
			lengthNow += 1
		}
	}
	if lengthNow < maxLen {
		for i := 0; i < (maxLen - lengthNow); i++ {
			str += " "
		}
	}
	return str
}
