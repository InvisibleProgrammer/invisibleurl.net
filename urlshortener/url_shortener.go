package urlshortener

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type UrlShortener struct {
	urlShortenerRepository *UrlShortenerRepository
}

func NewUrlShortener(urlShortenerRepository *UrlShortenerRepository) *UrlShortener {
	return &UrlShortener{
		urlShortenerRepository: urlShortenerRepository,
	}
}

func (urlShortener *UrlShortener) GetFullUrl(short string) (string, error) {

	allUrls, err := urlShortener.urlShortenerRepository.GetAll()
	if err != nil {
		return "", fmt.Errorf("couldn't get the shortened urls: %v", err)
	}

	for i := 0; i < len(allUrls); i++ {
		if allUrls[i].ShortUrl == short {
			return allUrls[i].FullUrl, nil
		}
	}

	return "", errors.New("couldn't find short URL")
}

func (urlShortener *UrlShortener) MakeShortUrl(userId string, fullUrl string) (string, error) {
	return "", nil
	// encoded := base62Encode(nextUrlId)
	// log.Default().Printf("encoded: %s", encoded)

	// shortenedUrls = append(shortenedUrls, ShortenedUrl{
	// 	UserId:      userId,
	// 	UrlId:       nextUrlId,
	// 	FullUrl:     fullUrl,
	// 	OriginalUrl: fullUrl,
	// 	ShortUrl:    encoded,
	// })

	// return encoded, nil
}

func DeleteShortUrl(userId string, shortUrl string) error {
	// panic("not implemented yet")
	log.Println("not implemented yet")

	return nil
	// for i := 0; i < len(shortenedUrls); i++ {
	// 	if shortenedUrls[i].UserId == "" {
	// 		break
	// 	}

	// 	if shortenedUrls[i].ShortUrl == shortUrl {
	// 		if shortenedUrls[i].UserId != userId {
	// 			return fmt.Errorf("cannot delete %s: user has no such short URL", shortUrl)
	// 		}

	// 		shortenedUrls = append(shortenedUrls[:i], shortenedUrls[i+1:]...)
	// 		return nil
	// 	}

	// }

	// return fmt.Errorf("couldn't find short URL: %s", shortUrl)
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
