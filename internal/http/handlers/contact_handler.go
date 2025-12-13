package handlers

import (
	"log"
	"net/http"
	"strings"

	"contactFormAPI/internal/repository"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	repo *repository.ContactRepository
}

func NewContactHandler(repo *repository.ContactRepository) *ContactHandler {
	return &ContactHandler{repo: repo}
}

type CreateContactRequest struct {
	Contact string `json:"contact" binding:"required" example:"example@example.com"`
	Name    string `json:"name" binding:"required" example:"山田太郎"`
	Message string `json:"message" binding:"required" example:"お問い合わせ内容"`
}

type CreateContactResponse struct {
	Status string `json:"status" example:"created"`
}

// CreateContact godoc
// @Summary 問い合わせを登録する
// @Description 問い合わせフォームから送信された情報を登録する
// @Tags Public API
// @Accept json
// @Produce json
// @Param request body CreateContactRequest true "問い合わせ情報"
// @Success 201 {object} CreateContactResponse
// @Failure 400 {object} map[string]string
// @Router /contact [post]
func (h *ContactHandler) CreateContact(c *gin.Context) {
	var req CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// バリデーション
	req.Contact = strings.TrimSpace(req.Contact)
	req.Name = strings.TrimSpace(req.Name)
	req.Message = strings.TrimSpace(req.Message)

	if req.Contact == "" || req.Name == "" || req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contact, name, and message are required and cannot be empty"})
		return
	}

	if len(req.Message) > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message must be 2000 characters or less"})
		return
	}

	// IPアドレスとUser-Agentを取得
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	var ipPtr, userAgentPtr *string
	if ip != "" {
		ipPtr = &ip
	}
	if userAgent != "" {
		userAgentPtr = &userAgent
	}

	_, err := h.repo.Create(req.Contact, req.Name, req.Message, ipPtr, userAgentPtr)
	if err != nil {
		// エラーの詳細をログに出力（本番環境では機密情報に注意）
		log.Printf("ERROR: Failed to create contact: %v", err)
		c.Error(err) // Ginのエラーログに記録
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create contact"})
		return
	}

	c.JSON(http.StatusCreated, CreateContactResponse{Status: "created"})
}
