package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/antontuzov/minimalytics/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	TrackPageView(path, referrer, userAgent, ip string) error
	GetDailyStats() ([]models.DailyStat, error)
	GetUniqueVisits() ([]models.UniqueVisitStat, error)
	GetTopPages() ([]models.PageStat, error)
	GetReferrers() ([]models.ReferrerStat, error)
	GetDevices() ([]models.DeviceStat, error)
	GetBrowsers() ([]models.BrowserStat, error)
	Close() error
}

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(path string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	s := &SQLiteStorage{db: db}
	go s.startRetentionPolicy()

	return s, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS page_views (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			path TEXT NOT NULL,
			referrer TEXT,
			user_agent TEXT,
			ip_address TEXT
		)`)
	return err
}

func (s *SQLiteStorage) TrackPageView(path, referrer, userAgent, ip string) error {
	_, err := s.db.Exec(
		"INSERT INTO page_views (timestamp, path, referrer, user_agent, ip_address) VALUES (?, ?, ?, ?, ?)",
		time.Now().Format(time.RFC3339),
		path,
		referrer,
		userAgent,
		anonymizeIP(ip),
	)
	return err
}

func (s *SQLiteStorage) GetDailyStats() ([]models.DailyStat, error) {
	return queryStats(s.db, `
		SELECT DATE(timestamp) as day, COUNT(*) as count 
		FROM page_views 
		GROUP BY day 
		ORDER BY day DESC
		LIMIT 30`, models.DailyStat{})
}

func (s *SQLiteStorage) GetUniqueVisits() ([]models.UniqueVisitStat, error) {
	return queryStats(s.db, `
		SELECT DATE(timestamp) as day, COUNT(DISTINCT ip_address) as count 
		FROM page_views 
		GROUP BY day 
		ORDER BY day DESC
		LIMIT 30`, models.UniqueVisitStat{})
}

func (s *SQLiteStorage) GetTopPages() ([]models.PageStat, error) {
	return queryStats(s.db, `
		SELECT path, COUNT(*) as count 
		FROM page_views 
		GROUP BY path 
		ORDER BY count DESC
		LIMIT 10`, models.PageStat{})
}

func (s *SQLiteStorage) GetReferrers() ([]models.ReferrerStat, error) {
	return queryStats(s.db, `
		SELECT referrer, COUNT(*) as count 
		FROM page_views 
		WHERE referrer != ''
		GROUP BY referrer 
		ORDER BY count DESC
		LIMIT 10`, models.ReferrerStat{})
}

func (s *SQLiteStorage) GetDevices() ([]models.DeviceStat, error) {
	return queryStats(s.db, `
		SELECT 
			CASE
				WHEN user_agent LIKE '%Mobile%' THEN 'Mobile'
				WHEN user_agent LIKE '%Tablet%' THEN 'Tablet'
				ELSE 'Desktop'
			END as device,
			COUNT(*) as count
		FROM page_views
		GROUP BY device
		ORDER BY count DESC`, models.DeviceStat{})
}

func (s *SQLiteStorage) GetBrowsers() ([]models.BrowserStat, error) {
	return queryStats(s.db, `
		SELECT 
			SUBSTR(user_agent, 0, INSTR(user_agent, '/')) as browser,
			COUNT(*) as count
		FROM page_views
		GROUP BY browser
		ORDER BY count DESC
		LIMIT 10`, models.BrowserStat{})
}

func queryStats[T any](db *sql.DB, query string, model T) ([]T, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []T
	for rows.Next() {
		var s T
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func anonymizeIP(ip string) string {
	if strings.Contains(ip, ":") {
		parts := strings.Split(ip, ":")
		if len(parts) >= 3 {
			return strings.Join(parts[:3], ":") + "::xxxx"
		}
	} else {
		parts := strings.Split(ip, ".")
		if len(parts) == 4 {
			return strings.Join(parts[:3], ".") + ".0"
		}
	}
	return "unknown"
}

func (s *SQLiteStorage) startRetentionPolicy() {
	retentionDays := os.Getenv("RETENTION_DAYS")
	if retentionDays == "" {
		retentionDays = "90"
	}

	for {
		_, err := s.db.Exec(`
			DELETE FROM page_views 
			WHERE timestamp < datetime('now', ?)
		`, "-"+retentionDays+" days")

		if err != nil {
			log.Printf("Retention policy error: %v", err)
		}
		time.Sleep(24 * time.Hour)
	}
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
