package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_USER") + ":" + os.Getenv("DATABASE_PASSWORD") + "@tcp(" + os.Getenv("DATABASE_HOST") + ":" + os.Getenv("DATABASE_PORT") + ")/" + os.Getenv("DATABASE_NAME") + "?parseTime=true"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 12
	}
	var posts []Post
	var total int64
	db.Model(&Post{}).Count(&total)
	db.Order("created_at desc").Limit(limit).Offset((page - 1) * limit).Find(&posts)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts": posts,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"count": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
			"prev_page": func() interface{} {
				if page > 1 {
					return page - 1
				} else {
					return nil
				}
			}(),
			"next_page": func() interface{} {
				if int64(page*limit) < total {
					return page + 1
				} else {
					return nil
				}
			}(),
		},
	})
}

func main() {
	initDB()
	http.HandleFunc("/posts", postsHandler)
	log.Println("Go API server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
