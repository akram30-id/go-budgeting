package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

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

func StrPadLeft(s string, length int, padChar string) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(padChar, length-len(s)) + s
}

func StrPadRight(s string, length int, padChar string) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(padChar, length-len(s))
}

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
