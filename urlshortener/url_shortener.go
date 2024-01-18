package urlshortener

import (
	"errors"
	"log"
	"strings"
)

type ShortenedUrl struct {
	UserId      string
	UrlId       int
	FullUrl     string
	OriginalUrl string
	ShortUrl    string
}

var shortenedUrls = []ShortenedUrl{
	{
		UserId:   "1",
		UrlId:    1,
		FullUrl:  "https://invisibleprogrammer.com",
		ShortUrl: "blog",
	},
}

func GetAll() []ShortenedUrl {
	return shortenedUrls
}

func GetFullUrl(short string) (string, error) {
	for i := 0; i < len(shortenedUrls); i++ {
		if shortenedUrls[i].UserId == "" {
			break
		}

		if shortenedUrls[i].ShortUrl == short {
			return shortenedUrls[i].FullUrl, nil
		}

	}

	return "", errors.New("couldn't find short URL")
}

func MakeShortUrl(userId string, fullUrl string) (string, error) {

	nextUrlId := getnextUrlId(shortenedUrls)
	log.Default().Printf("Next id: %d", nextUrlId)
	encoded := base62Encode(nextUrlId)
	log.Default().Printf("encoded: %s", encoded)

	shortenedUrls = append(shortenedUrls, ShortenedUrl{
		UserId:      userId,
		UrlId:       nextUrlId,
		FullUrl:     fullUrl,
		OriginalUrl: fullUrl,
		ShortUrl:    encoded,
	})

	return encoded, nil
}

func base62Encode(id int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var encoded strings.Builder

	for id > 0 {
		encoded.WriteByte(alphabet[id%62])
		id = id / 62
	}

	return encoded.String()
}

func getnextUrlId(shortenedUrls []ShortenedUrl) int {
	maxId := 20_000_000 // We start from this value to make sure the shortened version's length at least 5

	for _, v := range shortenedUrls {
		if v.UrlId > maxId {
			maxId = v.UrlId
		}
	}

	return maxId + 1
}
