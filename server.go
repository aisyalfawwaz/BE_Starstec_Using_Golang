package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	opt := option.WithCredentialsFile("starstec-2cf73-firebase-adminsdk-p5slh-58ab5b2048.json") // Sesuaikan dengan path ke serviceAccountKey Anda
	config := &firebase.Config{DatabaseURL: "https://starstec-2cf73-default-rtdb.firebaseio.com"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v", err)
	}
	defer client.Close()

	r.GET("/api/campaign/:campaignId", func(c *gin.Context) {
		campaignId := c.Param("campaignId")

		campaignDoc, err := client.Collection("campaign").Doc(campaignId).Get(context.Background())
		if err != nil {
			if campaignDoc == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		campaignData := campaignDoc.Data()
		c.JSON(http.StatusOK, campaignData)
	})

	serverAddr := fmt.Sprintf(":%s", port)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
