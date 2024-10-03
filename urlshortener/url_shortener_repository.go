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

func (repository *UrlShortenerRepository) Store(shortenedUrl ShortenedUrl) error {

	insertStmnt := `insert into short_urls (short_url_id, user_id, full_url, short_url) values (:shortUrlId, :userId, :fullUrl, :shortUrl)`

	parameters := map[string]interface{}{
		"shortUrlId": shortenedUrl.UrlId,
		"userId":     shortenedUrl.UserId,
		"fullUrl":    shortenedUrl.FullUrl,
		"shortUrl":   shortenedUrl.ShortUrl,
	}

	_, err := repository.db.Db.NamedExec(insertStmnt, parameters)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UrlShortenerRepository) GetNextUrlId() (int64, error) {
	selectStmnt := `select nextval('short_url_seq')`

	rows, err := repository.db.Db.Query(selectStmnt)
	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		return 0, err
	}

	var nextUrlId int64
	err = rows.Scan(&nextUrlId)
	if err != nil {
		return 0, err
	}

	return nextUrlId, nil
}

func (repository *UrlShortenerRepository) DeleteShortUrl(userId int64, shortUrl string) error {
	deleteStmnt := `delete from short_urls where user_id = :userId and short_url = :shortUrl`

	parameters := map[string]interface{}{
		"userId":   userId,
		"shortUrl": shortUrl,
	}

	_, err := repository.db.Db.NamedExec(deleteStmnt, parameters)
	if err != nil {
		return err
	}

	return nil
}
