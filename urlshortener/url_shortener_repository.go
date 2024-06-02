package urlshortener

import (
	"log"

	"invisibleprogrammer.com/invisibleurl/db"
)

type UrlShortenerRepository struct {
	db *db.Repository
}

func NewUrlShortenerRepository(db *db.Repository) *UrlShortenerRepository {
	return &UrlShortenerRepository{
		db: db,
	}
}

func (repository *UrlShortenerRepository) GetAll() ([]ShortenedUrl, error) {

	selectStmnt := `select short_url_id, user_id, full_url, short_url from short_urls`

	rows, err := repository.db.Db.Queryx(selectStmnt)
	if err != nil {
		return nil, err
	}

	var shortUrl ShortenedUrl
	var shortUrls []ShortenedUrl

	for rows.Next() {
		err := rows.StructScan(&shortUrl)

		if err != nil {
			log.Fatalln(err)
			return nil, err
		}

		shortUrls = append(shortUrls, shortUrl)
	}

	return shortUrls, nil
}

func (repository *UrlShortenerRepository) Store(ShortenedUrl) ([]ShortenedUrl, error) {

	selectStmnt := `insert into short_urls (user_id, full_url) values (Todo: how named parameters are working?)`

	rows, err := repository.db.Db.Queryx(selectStmnt)
	if err != nil {
		return nil, err
	}

	var shortUrl ShortenedUrl
	var shortUrls []ShortenedUrl

	for rows.Next() {
		err := rows.StructScan(&shortUrl)

		if err != nil {
			log.Fatalln(err)
			return nil, err
		}

		shortUrls = append(shortUrls, shortUrl)
	}

	return shortUrls, nil
}
