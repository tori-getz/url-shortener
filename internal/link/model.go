package link

import (
	"math/rand/v2"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url  string `json:"url"`
	Hash string `json:"hash" gorm:"uniqueIndex"`
}

func NewLink(url string) *Link {
	return &Link{
		Url:  url,
		Hash: randHash(10),
	}
}

var letterRunes = []rune("abcefghijklmnopqrstuvwxyzABCDEFGHIJKLMONPQRSTUVQXYZ")

func randHash(n int) string {
	generatedHash := make([]rune, n)

	for i := range generatedHash {
		generatedHash[i] = letterRunes[rand.IntN(len(letterRunes))]
	}

	return string(generatedHash)
}
