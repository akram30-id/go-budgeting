package helpers

import "strings"

func Contains(list []string, search string) bool {

	for _, item := range list {
		if item == search {
			return true
		}
	}

	return false

}

func Explode(delimiter, text string) []string {
	return strings.Split(text, delimiter)
}
