package urlshortener

import (
	"errors"
	"strings"

	"github.com/jxskiss/base62"
)

type ShortenedUrl struct {
	UserId      string
	FullUrl     string
	OriginalUrl string
	ShortUrl    string
}

var shortenedUrls = []ShortenedUrl{
	{
		UserId:   "1",
		FullUrl:  "https://invisibleprogrammer.com",
		ShortUrl: "blog",
	},
}

func GetAll() []ShortenedUrl {
	return shortenedUrls
}

func GetFullUrl(short string) (string, error) {
	lowerCasedShort := strings.ToLower(short)

	for i := 0; i < len(shortenedUrls); i++ {
		if shortenedUrls[i].UserId == "" {
			break
		}

		if shortenedUrls[i].ShortUrl == lowerCasedShort {
			return shortenedUrls[i].FullUrl, nil
		}

	}

	return "", errors.New("couldn't find short URL")
}

func MakeShortUrl(userId string, fullUrl string) (string, error) {

	/*

		Note: My memory completely faded about this algorithm.
		We should base62 encode a full url, just the Id of the URL.
		The Id should be a big bigint that ensures the shortened base62 encoded string
		will have at minimum 7 characters length

	*/
	encoded := base62.EncodeToString([]byte(fullUrl))

	for i := 0; i < len(shortenedUrls); i++ {
		if shortenedUrls[i].ShortUrl == encoded {
			return encoded, nil // Todo: figure out what to do if other user has the same shortened url
		}
	}

	shortenedUrls = append(shortenedUrls, ShortenedUrl{
		UserId:      userId,
		FullUrl:     fullUrl,
		OriginalUrl: fullUrl,
		ShortUrl:    encoded,
	})

	return encoded, nil
}
