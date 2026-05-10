package main

import (
	"log"
	"os"

	_ "team-management-service/docs"
	"team-management-service/config"
	"team-management-service/internal/handler"
	"team-management-service/internal/middleware"
	"team-management-service/internal/repository"
	"team-management-service/internal/usecase"
	"team-management-service/pkg/client"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Team Management Service API
// @version 1.0
// @description Teams, managers, and members. All JSON routes require a valid JWT from the Auth service. Use header: Authorization: Bearer followed by the token.
// @termsOfService http://swagger.io/terms/

// @contact.name Seta Golang Intern Project

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token issued by Auth service (include the word Bearer, a space, then the JWT)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	teamRepo := repository.NewTeamRepository(config.DB)
	authClient := client.NewAuthClient()
	teamUsecase := usecase.NewTeamUsecase(teamRepo, authClient)
	teamHandler := handler.NewTeamHandler(teamUsecase)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
