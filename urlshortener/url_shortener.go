package urlshortener

import "errors"

func GetFullUrl(short string) (string, error) {
	shortUrl := "blog"
	longUrl := "https://invisibleprogrammer.com"

	if short == shortUrl {
		return longUrl, nil
	}

	return "", errors.New("couldn't find short URL")
}
