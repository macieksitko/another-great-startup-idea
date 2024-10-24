package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type JobOffer struct {
	ID          int
	Title       string
	Company     string
	Description string
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
		rows, err := db.Query("SELECT id, title, company, description FROM job_offers")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var jobOffers []JobOffer
		for rows.Next() {
			var job JobOffer
			err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Description)
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

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":     "Job Board",
			"message":   "Available Job Offers",
			"jobOffers": jobOffers,
		})
	})

	r.POST("/search", func(c *gin.Context) {
		query := c.PostForm("query")
		rows, err := db.Query("SELECT id, title, company, description FROM job_offers WHERE title LIKE ?", "%"+query+"%")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var jobResults []JobOffer
		for rows.Next() {
			var job JobOffer
			err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Description)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			jobResults = append(jobResults, job)
		}

		time.Sleep(3 * time.Second)

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
		company := c.PostForm("offer-company")
	
		// Insert the new job offer
		result, err := db.Exec("INSERT INTO job_offers (title, description, company) VALUES (?, ?, ?)", title, description, company)
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
		rows, err := db.Query("SELECT id, title, company, description FROM job_offers ORDER BY id DESC")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
	
		var jobOffers []JobOffer
		for rows.Next() {
			var job JobOffer
			err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Description)
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

		time.Sleep(3 * time.Second)

	
		c.HTML(http.StatusOK, "job_list.html", gin.H{
			"title":     "Job Offers",
			"jobOffers": jobOffers,
			"newJobID":  newID,
		})
	})

	r.Run(":8080")
}