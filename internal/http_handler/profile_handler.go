package http_handler

import (
	"demo-signserver/internal/http_handler/application"
	"demo-signserver/internal/http_handler/dtos"
	"demo-signserver/internal/repository/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	Service *application.ProfileService
}

func NewProfileHandler(service *application.ProfileService) *ProfileHandler {
	return &ProfileHandler{Service: service}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var dto dtos.CreateSignerProfileDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile := domain.SignerProfile{
		Signer:      dto.Signer,
		ProfileId:   dto.ProfileId,
		Description: dto.Description,
		Configs:     dto.Configs,
		Upload:      dto.Upload,
		Download:    dto.Download,
	}
	ID, err := h.Service.CreateProfile(&profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Location", "/profiles/"+ID)
	c.JSON(http.StatusCreated, gin.H{"id": ID})
}

func (h *ProfileHandler) GetProfileByID(c *gin.Context) {
	id := c.Param("id")
	profile, err := h.Service.GetProfileByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	id := c.Param("id")
	var dto dtos.UpdateSignerProfileDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Monta struct parcial apenas com os campos preenchidos
	profile := &domain.SignerProfile{
		Description: dto.Description,
		Configs:     dto.Configs,
		Upload:      dto.Upload,
		Download:    dto.Download,
	}
	err := h.Service.UpdateProfile(id, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

func (h *ProfileHandler) RegisterRoutes(r *gin.Engine) {
	profiles := r.Group("/profiles")
	{
		profiles.POST("", h.CreateProfile)
		profiles.GET(":id", h.GetProfileByID)
		profiles.PATCH(":id", h.UpdateProfile)
	}
}

func RegisterProfileRoutes(r *gin.Engine) {
	profileService := application.NewProfileService()
	profileHandler := NewProfileHandler(profileService)
	profileHandler.RegisterRoutes(r)
}
