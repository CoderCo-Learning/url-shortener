package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)
	waitForDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handleHealth)
	mux.HandleFunc("/summary", handleSummary)
	mux.HandleFunc("/url/", handleURLStats)
	mux.HandleFunc("/recent", handleRecent)
	mux.HandleFunc("/top", handleTop)

	port := getEnv("PORT", "8081")
	log.Printf("Dashboard API listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	if err := db.Ping(); err != nil {
		status = "unhealthy"
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status, "service": "dashboard"})
}

func handleSummary(w http.ResponseWriter, r *http.Request) {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	var totalURLs, totalClicks, clicksToday int
	db.QueryRow("SELECT COUNT(*) FROM urls").Scan(&totalURLs)
	db.QueryRow("SELECT COALESCE(SUM(clicks), 0) FROM urls").Scan(&totalClicks)
	db.QueryRow(
		"SELECT COALESCE(SUM(clicks), 0) FROM click_stats_hourly WHERE hour >= $1",
		today,
	).Scan(&clicksToday)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_urls":   totalURLs,
		"total_clicks": totalClicks,
		"clicks_today": clicksToday,
	})
}

func handleURLStats(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/url/")
	if code == "" {
		httpError(w, "provide a short code", http.StatusBadRequest)
		return
	}

	var url string
	var clicks int
	var createdAt string
	err := db.QueryRow(
		"SELECT url, clicks, created_at FROM urls WHERE id = $1", code,
	).Scan(&url, &clicks, &createdAt)
	if err != nil {
		httpError(w, "not found", http.StatusNotFound)
		return
	}

	// Hourly stats for last 24 hours
	rows, err := db.Query(
		`SELECT hour, clicks FROM click_stats_hourly
		 WHERE short_code = $1 AND hour >= NOW() - INTERVAL '24 hours'
		 ORDER BY hour DESC`,
		code,
	)

	type HourlyStat struct {
		Hour   string `json:"hour"`
		Clicks int    `json:"clicks"`
	}

	hourly := []HourlyStat{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var h HourlyStat
			rows.Scan(&h.Hour, &h.Clicks)
			hourly = append(hourly, h)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"short_code":   code,
		"url":          url,
		"total_clicks": clicks,
		"created_at":   createdAt,
		"hourly":       hourly,
	})
}

func handleRecent(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(
		`SELECT short_code, ip_address, user_agent, clicked_at
		 FROM click_events ORDER BY clicked_at DESC LIMIT 50`,
	)
	if err != nil {
		httpError(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Click struct {
		ShortCode string `json:"short_code"`
		IP        string `json:"ip"`
		UserAgent string `json:"user_agent"`
		ClickedAt string `json:"clicked_at"`
	}

	clicks := []Click{}
	for rows.Next() {
		var c Click
		rows.Scan(&c.ShortCode, &c.IP, &c.UserAgent, &c.ClickedAt)
		clicks = append(clicks, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clicks)
}

func handleTop(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(
		"SELECT id, url, clicks FROM urls ORDER BY clicks DESC LIMIT 10",
	)
	if err != nil {
		httpError(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type TopURL struct {
		ShortCode string `json:"short_code"`
		URL       string `json:"url"`
		Clicks    int    `json:"clicks"`
	}

	top := []TopURL{}
	for rows.Next() {
		var t TopURL
		rows.Scan(&t.ShortCode, &t.URL, &t.Clicks)
		top = append(top, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(top)
}

func httpError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func waitForDB() {
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			return
		}
		log.Printf("Waiting for database... (%d/30)", i+1)
		time.Sleep(time.Second)
	}
	log.Fatal("Database not ready after 30s")
}
