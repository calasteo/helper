package helper

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

type RandomOption struct {
	LowerCase   bool
	UpperCase   bool
	Numeric     bool
	SpecialChar bool
}

var (
	LowerCase        = "abcdefghijklmnopqrstuvwxyz"
	UpperCase        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numeric          = "1234567890"
	SpecialCharacter = "!@#$%^&*()_+-."
)

func (a *RandomOption) GenerateRandomString(length int) (string, error) {
	if length == 0 {
		return "", errors.New("minimum 1 character's length")
	}
	minChar := 0
	rawSeed := ""
	result := ""

	if a.LowerCase {
		minChar++
		rawSeed += LowerCase
		rand.Seed(time.Now().UnixNano())
		result += string([]rune(LowerCase)[rand.Intn(len(LowerCase))])
	}
	if a.UpperCase {
		minChar++
		rawSeed += UpperCase
		rand.Seed(time.Now().UnixNano())
		result += string([]rune(UpperCase)[rand.Intn(len(UpperCase))])
	}
	if a.Numeric {
		minChar++
		rawSeed += Numeric
		rand.Seed(time.Now().UnixNano())
		result += string([]rune(Numeric)[rand.Intn(len(Numeric))])
	}
	if a.SpecialChar {
		minChar++
		rawSeed += SpecialCharacter
		rand.Seed(time.Now().UnixNano())
		result += string([]rune(SpecialCharacter)[rand.Intn(len(SpecialCharacter))])
	}
	if minChar == 0 {
		b := make([]rune, length)
		rand.Seed(time.Now().UnixNano())
		for i := range b {
			b[i] = []rune(SpecialCharacter)[rand.Intn(len(SpecialCharacter))]
		}
		return string(b), nil
	}

	if minChar > length {
		return "", errors.New("min char is " + strconv.Itoa(minChar))
	}
	for i := minChar; i < length; i++ {
		rand.Seed(time.Now().UnixNano())
		result += string([]rune(rawSeed)[rand.Intn(len(rawSeed))])
	}
	return string(ShuffleBytes([]byte(result))), nil

}

func ShuffleBytes(src []byte) []byte {
	final := make([]byte, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}
	return final
}
