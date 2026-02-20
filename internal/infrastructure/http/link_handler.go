package http

import (
	"net/http"
	"strconv"

	"app/internal/application/link"
	domainLink "app/internal/domain/link"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *link.Service
}

func NewHandler(service *link.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/links")
	{
		api.GET("", h.GetAll)
		api.POST("", h.Create)
		api.GET("/:id", h.GetByID)
		api.PUT("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
}

type LinkResponse struct {
	ID          int64  `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortURL    string `json:"short_url"`
}

func toLinkResponse(l *domainLink.Link, service *link.Service) LinkResponse {
	return LinkResponse{
		ID:          l.ID,
		OriginalURL: l.OriginalURL,
		ShortName:   l.ShortName,
		ShortURL:    service.GetShortURL(l),
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	links, err := h.service.GetAllLinks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]LinkResponse, len(links))
	for i, l := range links {
		response[i] = toLinkResponse(l, h.service)
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	linkEntity, err := h.service.CreateLink(c.Request.Context(), req.OriginalURL, req.ShortName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toLinkResponse(linkEntity, h.service))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	linkEntity, err := h.service.GetLink(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}

	c.JSON(http.StatusOK, toLinkResponse(linkEntity, h.service))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	linkEntity, err := h.service.UpdateLink(c.Request.Context(), id, req.OriginalURL, req.ShortName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toLinkResponse(linkEntity, h.service))
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteLink(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
