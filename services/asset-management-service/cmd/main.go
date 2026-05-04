package main

import (
	"log"
	"os"

	"asset-management-service/config"
	"asset-management-service/internal/handler"
	"asset-management-service/internal/middleware"
	"asset-management-service/internal/repository"
	"asset-management-service/internal/usecase"
	"asset-management-service/pkg/client"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	repo := repository.NewAssetRepository(config.DB)
	authClient := client.NewAuthClient()
	teamClient := client.NewTeamClient()
	assetUsecase := usecase.NewAssetUsecase(repo, authClient, teamClient)
	assetHandler := handler.NewAssetHandler(assetUsecase)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "UP",
			"service": "asset-management-service",
		})
	})

	api := r.Group("/api/v1")
	{
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired(authClient))
		{
			protected.POST("/folders", assetHandler.CreateFolder)
			protected.GET("/folders", assetHandler.ListFolders)
			protected.GET("/folders/:folderId", assetHandler.GetFolder)
			protected.PATCH("/folders/:folderId", assetHandler.UpdateFolder)
			protected.DELETE("/folders/:folderId", assetHandler.DeleteFolder)

			protected.POST("/folders/:folderId/notes", assetHandler.CreateNote)
			protected.GET("/folders/:folderId/notes", assetHandler.ListNotesByFolder)
			protected.GET("/notes/:noteId", assetHandler.GetNote)
			protected.PATCH("/notes/:noteId", assetHandler.UpdateNote)
			protected.DELETE("/notes/:noteId", assetHandler.DeleteNote)

			protected.POST("/shares", assetHandler.ShareAsset)
			protected.DELETE("/shares/:shareId", assetHandler.RevokeShare)
			protected.GET("/shares/received", assetHandler.ListReceivedShares)
			protected.GET("/shares/granted", assetHandler.ListGrantedShares)
		}
	}

	port := os.Getenv("ASSET_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Asset service is starting on port http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
