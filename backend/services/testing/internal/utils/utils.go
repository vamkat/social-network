package utils

import (
	"fmt"
	"math/rand/v2"
)

func RandomString(length int) string {
	runes := make([]rune, length)
	for i := range length {
		runes[i] = RandomRune()
	}
	fmt.Println("created word:", string(runes))
	return string(runes)
}

var runes = []rune("abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789" + "_-.")

func RandomRune() rune {
	return runes[rand.IntN(len(runes))]
}
