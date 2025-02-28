package link

import (
	"math/rand/v2"
	"url-shortener/internal/stat"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url   string      `json:"url"`
	Hash  string      `json:"hash" gorm:"uniqueIndex"`
	Stats []stat.Stat `json:"stats" gorm:"constraints:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewLink(url string) *Link {
	link := &Link{
		Url: url,
	}
	link.GenerateHash()

	return link
}

func (link *Link) GenerateHash() {
	link.Hash = randHash(10)
}

var letterRunes = []rune("abcefghijklmnopqrstuvwxyzABCDEFGHIJKLMONPQRSTUVQXYZ")

func randHash(n int) string {
	generatedHash := make([]rune, n)

	for i := range generatedHash {
		generatedHash[i] = letterRunes[rand.IntN(len(letterRunes))]
	}

	return string(generatedHash)
}
