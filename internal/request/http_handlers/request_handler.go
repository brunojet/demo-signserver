package http_handlers

import (
	"demo-signserver/internal/request/dtos"
	"demo-signserver/internal/request/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	Service *services.RequestService
}

func NewRequestHandler(service *services.RequestService) *RequestHandler {
	return &RequestHandler{Service: service}
}

func (h *RequestHandler) CreateRequest(c *gin.Context) {
	var dto dtos.CreateSignRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := dto.GetDomainCreateSignRequest()
	ID, err := h.Service.CreateRequest(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Location", "/requests/"+ID)
	c.JSON(http.StatusCreated, gin.H{"id": ID})
}

func (h *RequestHandler) GetRequestByID(c *gin.Context) {
	id := c.Param("id")
	response, err := h.Service.GetRequestByID(id)
	if err != nil || response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *RequestHandler) RegisterRoutes(r *gin.Engine) {
	requests := r.Group("/requests")
	{
		requests.POST("", h.CreateRequest)
		requests.GET(":id", h.GetRequestByID)
	}
}

func RegisterRequestRoutes(r *gin.Engine) {
	requestService := services.NewRequestService()
	handler := NewRequestHandler(requestService)
	handler.RegisterRoutes(r)
}
