package handler
import (
"net/http"
"github.com/gin-gonic/gin"
"auth-user-management-service/internal/usecase"
)
// AuthHandler holds the dependencies for HTTP routing
type AuthHandler struct {
authUsecase usecase.AuthUsecase
}
// NewAuthHandler injects usecase layer into handler
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
return &AuthHandler{
authUsecase: authUsecase,
}
}
// Register (POST /api/v1/auth/register)
func (h *AuthHandler) Register(c *gin.Context) {
// Parse the JSON request body
var req usecase.RegisterRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 400 Bad Request
return
}
// Process business logic
newUser, err := h.authUsecase.Register(&req)
if err != nil {
// Status code can be dynamic (e.g., 409 Conflict if email exists)
c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // 409 Conflict
return
}
// Respond with Created Entity (Password hidden by json tags in Struct)
c.JSON(http.StatusCreated, gin.H{
"message": "User registered successfully",
"user":    newUser,
})
}
