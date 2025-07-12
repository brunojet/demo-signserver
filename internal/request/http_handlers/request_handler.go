package http_handlers

import (
	"demo-signserver/internal/observability"
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

	domain := dto.GetDomainCreateSignRequest()
	response, err := h.Service.CreateRequest(domain)
	if err != nil {
		httpStatus := http.StatusInternalServerError
		if contains(response.SignerError.Message, "Erro ao buscar item com ID") {
			httpStatus = http.StatusConflict
		}
		c.JSON(httpStatus, response)
		return
	}
	c.Header("Location", "/requests/"+*response.ID)
	c.JSON(http.StatusCreated, response)
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

func (h *RequestHandler) GetSignerStatusByID(c *gin.Context) {
	id := c.Param("id")
	response, err := h.Service.GetSignerStatusByID(id)
	if err != nil || response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *RequestHandler) RegisterRoutes(r *gin.Engine) {
	requests := r.Group("/requests")
	{
		requests.POST("", observability.ObservableMiddleware(func(c *gin.Context, obs *observability.Observability) {
			h.CreateRequest(c)
		}))
		requests.GET(":id", observability.ObservableMiddleware(func(c *gin.Context, obs *observability.Observability) {
			h.GetSignerStatusByID(c)
		}))
	}
	requests = r.Group("/manage_requests")
	{
		requests.GET(":id", observability.ObservableMiddleware(func(c *gin.Context, obs *observability.Observability) {
			h.GetRequestByID(c)
		}))
	}
}

func RegisterRequestRoutes(r *gin.Engine) {
	requestService := services.NewRequestService()
	handler := NewRequestHandler(requestService)
	handler.RegisterRoutes(r)
}
