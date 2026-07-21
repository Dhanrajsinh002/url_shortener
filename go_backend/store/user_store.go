package store

import "time"

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

func GetUserByUsername(username string) (*User, error) {
	var u User
	err := pgPool.QueryRow(ctx,
		`SELECT id, username, password_hash FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash)
	if err != nil {
		return nil, ErrNotFound
	}
	return &u, nil
}

func CreateUser(username, passwordHash string) error {
	_, err := pgPool.Exec(ctx,
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)`,
		username, passwordHash,
	)
	return err
}

type UrlRecord struct {
	ShortCode  string `json:"short_code"`
	LongUrl    string `json:"long_url"`
	ClickCount int64  `json:"click_count"`
	CreatedAt  string `json:"created_at"`
}

func SaveUrlMappingForUser(shortUrl, longUrl string, userID int64) error {
	_, err := pgPool.Exec(ctx,
		`INSERT INTO urls (short_code, long_url, user_id) VALUES ($1, $2, $3)`,
		shortUrl, longUrl, userID,
	)
	return err
}

func ListUrlsByUser(userID int64) ([]UrlRecord, error) {
	rows, err := pgPool.Query(ctx,
		`SELECT short_code, long_url, click_count, created_at
		 FROM urls WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []UrlRecord
	for rows.Next() {
		var r UrlRecord
		var createdAt time.Time
		if err := rows.Scan(&r.ShortCode, &r.LongUrl, &r.ClickCount, &createdAt); err != nil {
			return nil, err
		}
		r.CreatedAt = createdAt.Format(time.RFC3339)
		records = append(records, r)
	}
	return records, nil
}

// func ListAllUrls() ([]UrlRecord, error) {
// 	rows, err := pgPool.Query(ctx, `SELECT short_code, long_url, click_count, created_at FROM urls ORDER BY created_at DESC`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var records []UrlRecord
// 	for rows.Next() {
// 		var r UrlRecord
// 		var createdAt time.Time
// 		if err := rows.Scan(&r.ShortCode, &r.LongUrl, &r.ClickCount, &createdAt); err != nil {
// 			return nil, err
// 		}
// 		r.CreatedAt = createdAt.Format(time.RFC3339)
// 		records = append(records, r)
// 	}
// 	return records, nil
// } 