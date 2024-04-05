package currency

import "strings"

func SeparateCurrency(currencies string) (string, string) {
	return strings.Split(currencies, "/")[0], strings.Split(currencies, "/")[1]
}
