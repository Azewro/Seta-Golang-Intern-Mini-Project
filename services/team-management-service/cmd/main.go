package main

import (
	"log"
	"os"

	"team-management-service/config"
	"team-management-service/internal/handler"
	"team-management-service/internal/middleware"
	"team-management-service/internal/repository"
	"team-management-service/internal/usecase"
	"team-management-service/pkg/client"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	teamRepo := repository.NewTeamRepository(config.DB)
	authClient := client.NewAuthClient()
	teamUsecase := usecase.NewTeamUsecase(teamRepo, authClient)
	teamHandler := handler.NewTeamHandler(teamUsecase)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "UP",
			"service": "team-management-service",
		})
	})

	api := r.Group("/api/v1")
	{
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired(authClient))
		{
			teams := protected.Group("/teams")
			{
				teams.GET("/my", teamHandler.ListMyTeams)
				teams.GET("/:teamId", teamHandler.GetTeam)

				manager := teams.Group("")
				manager.Use(middleware.ManagerOnly())
				{
					manager.POST("", teamHandler.CreateTeam)
					manager.POST("/:teamId/members", teamHandler.AddMember)
					manager.DELETE("/:teamId/members/:userId", teamHandler.RemoveMember)
					manager.POST("/:teamId/managers", teamHandler.AddManager)
					manager.DELETE("/:teamId/managers/:userId", teamHandler.RemoveManager)
				}
			}
		}
	}

	port := os.Getenv("TEAM_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Team service is starting on port http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
