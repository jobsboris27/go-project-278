package http

import (
	"net/http"
	"strconv"
	"time"

	"app/internal/application/link"
	linkdomain "app/internal/domain/link"

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
	router.GET("/r/:code", h.Redirect)

	api := router.Group("/api/links")
	{
		api.GET("", h.GetAll)
		api.POST("", h.Create)
		api.GET("/:id", h.GetByID)
		api.PUT("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}

	apiVisits := router.Group("/api")
	{
		apiVisits.GET("/link_visits", h.GetVisits)
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

type VisitResponse struct {
	ID        int64  `json:"id"`
	LinkID    int64  `json:"link_id"`
	CreatedAt string `json:"created_at"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
	Status    int    `json:"status"`
}

func toLinkResponse(l *linkdomain.Link, service *link.Service) LinkResponse {
	return LinkResponse{
		ID:          l.ID,
		OriginalURL: l.OriginalURL,
		ShortName:   l.ShortName,
		ShortURL:    service.GetShortURL(l),
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	rangeStr := c.Query("range")
	pagination, err := linkdomain.ParseRange(rangeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid range format"})
		return
	}

	links, total, err := h.service.GetAllLinks(c.Request.Context(), pagination.Offset, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Range", pagination.ContentRange(total))

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

func (h *Handler) Redirect(c *gin.Context) {
	code := c.Param("code")

	linkEntity, err := h.service.GetLinkByShortName(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	referer := c.GetHeader("Referer")

	_ = h.service.RecordVisit(c.Request.Context(), linkEntity.ID, ip, userAgent, referer, http.StatusFound)

	c.Redirect(http.StatusFound, linkEntity.OriginalURL)
}

func (h *Handler) GetVisits(c *gin.Context) {
	rangeStr := c.Query("range")
	pagination, err := linkdomain.ParseRange(rangeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid range format"})
		return
	}

	visits, total, err := h.service.GetVisits(c.Request.Context(), pagination.Offset, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Range", pagination.ContentRange(total))

	response := make([]VisitResponse, len(visits))
	for i, v := range visits {
		response[i] = VisitResponse{
			ID:        v.ID,
			LinkID:    v.LinkID,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
			IP:        v.IP,
			UserAgent: v.UserAgent,
			Referer:   v.Referer,
			Status:    v.Status,
		}
	}

	c.JSON(http.StatusOK, response)
}
