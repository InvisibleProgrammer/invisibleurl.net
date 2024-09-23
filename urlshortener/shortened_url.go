package urlshortener

type ShortenedUrl struct {
	UserId   int64  `db:"user_id"`
	UrlId    int    `db:"short_url_id"`
	FullUrl  string `db:"full_url"`
	ShortUrl string `db:"short_url"`
}
