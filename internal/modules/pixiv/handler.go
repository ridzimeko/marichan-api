package pixiv

import (
	"io"
	"marichan-api/internal/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetArtworkDetail(c *gin.Context) {
	artworkID := c.Param("id")

	if artworkID == "" {
		response.Error(c, http.StatusBadRequest, "artwork id is required", nil)
		return
	}

	result, err := h.service.GetArtworkDetail(c.Request.Context(), artworkID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.JSON(c, http.StatusOK, true, "Pixiv artwork detail fetched", result)
}

func (h *Handler) SearchArtworks(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		response.Error(c, http.StatusBadRequest, "query q is required", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	order := c.DefaultQuery("order", "date_d")
	mode := c.DefaultQuery("mode", "all")
	sMode := c.DefaultQuery("s_mode", "s_tag")

	result, err := h.service.searchArtworks(c.Request.Context(), keyword, order, mode, sMode, page)
	if err != nil {
		response.Error(c, http.StatusBadGateway, err.Error(), nil)
		return
	}

	response.JSON(c, http.StatusOK, true, "Pixiv search success", result)
}

func (h *Handler) GetArtistDetail(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.Error(c, http.StatusBadRequest, "Artist id is required", nil)
		return
	}

	result, err := h.service.GetArtistDetail(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusBadGateway, err.Error(), nil)
		return
	}

	response.JSON(c, http.StatusOK, true, "Pixiv artist detail fetched", result)
}

func (h *Handler) DownloadImage(c *gin.Context) {
	imageURL := c.Query("url")
	if imageURL == "" {
		response.Error(c, http.StatusBadRequest, "image url is required", nil)
		return
	}

	resp, err := h.service.DownloadImage(c.Request.Context(), imageURL)
	if err != nil {
		response.Error(c, http.StatusBadGateway, err.Error(), nil)
		return
	}
	defer resp.Body.Close()

	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		c.Writer.Header().Set("Content-Type", contentType)
	}
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		c.Writer.Header().Set("Content-Length", contentLength)
	}

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func (h *Handler) DownloadArtwork(c *gin.Context) {
	artworkID := c.Query("id")

	if artworkID == "" {
		response.Error(c, http.StatusBadRequest, "valid artwork url or id is required", nil)
		return
	}

	detail, err := h.service.GetArtworkDetail(c.Request.Context(), artworkID)
	if err != nil {
		response.Error(c, http.StatusBadGateway, err.Error(), nil)
		return
	}

	if detail.Pages == nil || len(*detail.Pages) == 0 {
		response.Error(c, http.StatusNotFound, "artwork has no pages", nil)
		return
	}

	// Default to downloading the original quality of the first page
	originalURL := (*detail.Pages)[0].URLs.Original

	resp, err := h.service.DownloadImage(c.Request.Context(), originalURL)
	if err != nil {
		response.Error(c, http.StatusBadGateway, err.Error(), nil)
		return
	}
	defer resp.Body.Close()

	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		c.Writer.Header().Set("Content-Type", contentType)
	}
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		c.Writer.Header().Set("Content-Length", contentLength)
	}

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
