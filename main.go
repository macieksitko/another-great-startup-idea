package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type JobOffer struct {
	ID          int
	Title       string
	Author     string
	Description string
	ViewsCount  int
	CreatedAt   time.Time
}

type GetJobOffer struct {
	ID          int
	Title       string
	Author     string
	Description string
	DaysAgo     string
}

func main() {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", "./jobs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()
	
	r.Static("/public", "./public")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title, author, description, created_at FROM job_offers")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var jobOffers []GetJobOffer
		for rows.Next() {
			var job GetJobOffer
			var createdAt time.Time

			err := rows.Scan(&job.ID, &job.Title, &job.Author, &job.Description, &createdAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			job.DaysAgo = formatTimeAgo(createdAt)
			jobOffers = append(jobOffers, job)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":     "Job Board",
			"message":   "Available Job Offers",
			"jobOffers": jobOffers,
		})
	})

	r.POST("/search", func(c *gin.Context) {
		query := c.PostForm("query")
		rows, err := db.Query("SELECT id, title, author, description, created_at FROM job_offers WHERE title LIKE ?", "%"+query+"%")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var jobResults []GetJobOffer
		for rows.Next() {
			var job GetJobOffer
			var createdAt time.Time

			err := rows.Scan(&job.ID, &job.Title, &job.Author, &job.Description, &createdAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			job.DaysAgo = formatTimeAgo(createdAt)
			jobResults = append(jobResults, job)
		}

		time.Sleep(2 * time.Second)

		c.HTML(http.StatusOK, "job_list.html", gin.H{
			"title":     "Search Results",
			"query":     query,
			"jobOffers":  jobResults,
		})
	})

	r.GET("/new_offer_modal", func(c* gin.Context) {
		c.HTML(http.StatusOK, "new_offer_modal.html", gin.H{
			"title": "New Offer",
		})
	})
	r.POST("/new_offer", func(c *gin.Context) {
		title := c.PostForm("offer-title")
		description := c.PostForm("offer-description")
		author := c.PostForm("offer-author")
	
		// Insert the new job offer
		result, err := db.Exec("INSERT INTO job_offers (title, description, author) VALUES (?, ?, ?)", title, description, author)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		// Get the ID of the newly inserted job offer
		newID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		// Fetch all job offers, including the new one
		rows, err := db.Query("SELECT id, title, author, description FROM job_offers ORDER BY id DESC")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
	
		var jobOffers []JobOffer
		for rows.Next() {
			var job JobOffer
			err := rows.Scan(&job.ID, &job.Title, &job.Author, &job.Description)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			jobOffers = append(jobOffers, job)
		}
	
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		time.Sleep(2 * time.Second)

	
		c.HTML(http.StatusOK, "job_list.html", gin.H{
			"title":     "Job Offers",
			"jobOffers": jobOffers,
			"newJobID":  newID,
		})
	})

	r.Run(":8080")
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)
	days := int(duration.Hours() / 24)
	
	switch {
	case days == 0:
		return "Today"
	case days == 1:
		return "Yesterday"
	default:
		return fmt.Sprintf("%d days ago", days)
	}
}