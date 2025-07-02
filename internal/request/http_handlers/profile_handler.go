package http_handlers

import (
	"demo-signserver/internal/observability"
	"demo-signserver/internal/request/dtos"
	"demo-signserver/internal/request/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	Service *services.ProfileService
}

func NewProfileHandler(service *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{Service: service}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var dto dtos.CreateSignerProfileDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile := dto.GetDomainCreateSignerProfile()
	ID, err := h.Service.CreateProfile(profile)
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
	profile := dto.GetDomainUpdateSignerProfile()

	err := h.Service.UpdateProfile(id, profile)
	if err != nil {
		if isDynamoDBValidationException(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

func (h *ProfileHandler) RegisterRoutes(r *gin.Engine) {
	profiles := r.Group("/profiles")
	{
		profiles.POST("", observability.ObservableMiddleware(func(c *gin.Context, obs *observability.Observability) {
			h.CreateProfile(c)
		}))
		profiles.GET(":id", observability.ObservableMiddleware(func(c *gin.Context, obs *observability.Observability) {
			h.GetProfileByID(c)
		}))
		profiles.PATCH(":id", observability.ObservableMiddleware(func(c *gin.Context, obs *observability.Observability) {
			h.UpdateProfile(c)
		}))
	}
}

func RegisterProfileRoutes(r *gin.Engine) {
	profileService := services.NewProfileService()
	profileHandler := NewProfileHandler(profileService)
	profileHandler.RegisterRoutes(r)
}

func isDynamoDBValidationException(err error) bool {
	if err != nil &&
		(len(err.Error()) > 0 &&
			((contains(err.Error(), "ValidationException") && contains(err.Error(), "number of conditions on the keys is invalid")) ||
				contains(err.Error(), "ConditionalCheckFailed"))) {
		fmt.Printf("DynamoDB validation exception: %v\n", err)
		return true
	}
	return false
}
