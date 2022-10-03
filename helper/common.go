package helper

import (
	"strings"
)

func Purify(str *string) string {
	purifiedString := strings.Trim(*str, " ")
	return purifiedString
}
