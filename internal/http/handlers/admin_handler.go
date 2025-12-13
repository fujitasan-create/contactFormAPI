package handlers

import (
	"net/http"

	"contactFormAPI/internal/auth"
	"contactFormAPI/internal/config"
	"contactFormAPI/internal/repository"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	repo   *repository.ContactRepository
	config *config.Config
}

func NewAdminHandler(repo *repository.ContactRepository, cfg *config.Config) *AdminHandler {
	return &AdminHandler{repo: repo, config: cfg}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Login godoc
// @Summary Adminログイン
// @Description AdminユーザーでログインしてJWTトークンを取得する
// @Tags Admin API
// @Accept json
// @Produce json
// @Param request body LoginRequest true "ログイン情報"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string
// @Router /admin/login [post]
func (h *AdminHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// configから取得
	adminUsername := h.config.AdminUsername
	adminPasswordHash := h.config.AdminPasswordHash

	// Usernameの検証
	if req.Username != adminUsername {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Passwordの検証（bcrypt）
	if !auth.VerifyPassword(adminPasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// JWTトークンを発行
	token, err := auth.GenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{AccessToken: token})
}

// GetMessages godoc
// @Summary 問い合わせ一覧取得
// @Description Adminユーザーが問い合わせ一覧を取得する（JWT認証必須）
// @Tags Admin API
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} repository.Contact
// @Failure 401 {object} map[string]string
// @Router /admin/messages [get]
func (h *AdminHandler) GetMessages(c *gin.Context) {
	contacts, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
		return
	}

	// IPとUser-Agentはレスポンスから除外（セキュリティ上の理由）
	type ContactResponse struct {
		ID        int64  `json:"id"`
		Contact   string `json:"contact"`
		Name      string `json:"name"`
		Message   string `json:"message"`
		CreatedAt string `json:"created_at"`
	}

	var response []ContactResponse
	for _, contact := range contacts {
		response = append(response, ContactResponse{
			ID:        contact.ID,
			Contact:   contact.Contact,
			Name:      contact.Name,
			Message:   contact.Message,
			CreatedAt: contact.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, response)
}

