package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type NasaPhoto struct {
	Copyright   string `json:"copyright"`
	Date        string `json:"date"`
	Explanation string `json:"explanation"`
	Hdurl       string `json:"hdurl"`
	Title       string `json:"title"`
}

func formatDate(date string) (string, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("error parsing date: %w", err)
	}

	formattedDate := parsedDate.Format("01/02/2006")
	return formattedDate, nil
}

func fetchNasaPhoto(apiKey string) (NasaPhoto, error) {
	url := "https://api.nasa.gov/planetary/apod?api_key=" + apiKey
	resp, err := http.Get(url)
	if err != nil {
		return NasaPhoto{}, fmt.Errorf("error fetching data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NasaPhoto{}, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var photo NasaPhoto
	if err := json.NewDecoder(resp.Body).Decode(&photo); err != nil {
		return NasaPhoto{}, fmt.Errorf("error decoding response: %w", err)
	}

	photo.Date, err = formatDate(photo.Date)

	if err != nil {
		fmt.Printf("Problem reformatting date")
	}

  // fmt.Printf(photo.Explanation)

	return photo, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY is not set in the .env file")
	}

	router := gin.Default()
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		photo, err := fetchNasaPhoto(apiKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch NASA photo"})
			return
		}

		c.HTML(http.StatusOK, "index.html", photo)
	})

	router.LoadHTMLGlob("./templates/*")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Server running on http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
