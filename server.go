package main

import (
    "context"
    "log"
    "net/http"
    "os"

    firebase "firebase.google.com/go"
    "google.golang.org/api/option"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    r := gin.Default()
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    opt := option.WithCredentialsFile("starstec-2cf73-firebase-adminsdk-p5slh-58ab5b2048.json")
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

    // Middleware CORS
    configCORS := cors.DefaultConfig()
    configCORS.AllowOrigins = []string{"*"}
    configCORS.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    r.Use(cors.New(configCORS))

    // Endpoint untuk mendapatkan semua data campaign dalam bentuk array JSON
    r.GET("/api/campaigns", func(c *gin.Context) {
        ctx := context.Background()
        campaignRef := client.Collection("campaign")
        docs, err := campaignRef.Documents(ctx).GetAll()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
            return
        }

        var campaigns []gin.H
        for _, doc := range docs {
            campaignData := doc.Data()
            campaign := gin.H{
                "id":          doc.Ref.ID,
                "title":       campaignData["Title"],
                "description": campaignData["Description"],
                "date":        campaignData["Date"],
                "image":       campaignData["image"].(string), // Mengonversi ke string jika diperlukan
            }
            campaigns = append(campaigns, campaign)
        }

        c.JSON(http.StatusOK, campaigns)
    })

    // Endpoint untuk mendapatkan data campaign berdasarkan ID dalam bentuk objek JSON
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
        campaign := gin.H{
            "id":          campaignId,
            "title":       campaignData["Title"],
            "description": campaignData["Description"],
            "date":        campaignData["Date"],
            "image":       campaignData["image"].(string), // Mengonversi ke string jika diperlukan
        }

        c.JSON(http.StatusOK, campaign)
    })

    serverAddr := ":" + port
    if err := r.Run(serverAddr); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}
