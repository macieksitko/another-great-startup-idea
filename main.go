package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	client "github.com/macieksitko/another-great-startup-idea/embeddings/client"
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

type CreateJobOffer struct {
	Title    string `form:"title" validate:"required"`
	Description   string `form:"description" validate:"required"`
	Author     string `form:"author" validate:"required"`
}

func main() {
	client := client.NewClient(
        "http://localhost:8000",
        "",
    )
	// Open the SQLite database
	sqlite_vec.Auto()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

		// Run initialization script
	err = initDatabaseFromFile(db, "init.sql")
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	
	r.Static("/public", "./public")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title, author, description, created_at FROM job_offers ORDER BY id DESC")
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
		rows, err := db.Query("SELECT id, title, author, description, created_at FROM job_offers WHERE title LIKE ? ORDER BY id DESC", "%"+query+"%")
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

		time.Sleep(1 * time.Second)

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
		var newOffer CreateJobOffer

		validate := validator.New()

		// Bind the form data to the struct
		if err := c.ShouldBind(&newOffer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the struct
		if err := validate.Struct(newOffer); err != nil {
			errorMessages := make(map[string]string)
			for _, err := range err.(validator.ValidationErrors) {
				fieldName := err.Field()
				errorMessages[fieldName] = fmt.Sprintf("%s is %s", fieldName, err.Tag())
			}
			c.HTML(http.StatusBadRequest, "new_offer_modal.html", gin.H{
				"title": "New Offer with error",
				"errors": errorMessages,
			})
			return
		}

	
		// Insert the new job offer
		result, err := db.Exec("INSERT INTO job_offers (title, description, author) VALUES (?, ?, ?)", newOffer.Title, newOffer.Description, newOffer.Author)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		// Get the ID of the newly inserted job offer
		_, err = result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		// Fetch all job offers, including the new one
		rows, err := db.Query("SELECT id, title, author, description, created_at FROM job_offers ORDER BY id DESC")
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
			fmt.Println((job.Description))

			jobOffers = append(jobOffers, job)

			embeddings := getEmbeddings(job.Description)
		}
	
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		time.Sleep(2 * time.Second)

	
		c.HTML(http.StatusOK, "job_list.html", gin.H{
			"title":     "Job Offers",
			"jobOffers": jobOffers,
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

func initDatabaseFromFile(db *sql.DB, filename string) error {
	// Read the contents of the SQL file
	initScript, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading init file: %v", err)
	}

	// Execute the script
	_, err = db.Exec(string(initScript))
	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}

	return nil
}