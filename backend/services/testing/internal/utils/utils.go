package utils

import (
	"math/rand/v2"
)

func RandomString(length int, withSpecialChars bool) string {
	runes := make([]rune, length)
	for i := range length {
		runes[i] = RandomRune(withSpecialChars)
	}
	// fmt.Println("created word:", string(runes))
	return string(runes)
}

var runes = []rune("abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789" + "_-.")
var lettersOnly = []rune("abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomRune(withSpecialChars bool) rune {
	if withSpecialChars {
		return runes[rand.IntN(len(runes))]
	}
	return lettersOnly[rand.IntN(len(lettersOnly))]

}
