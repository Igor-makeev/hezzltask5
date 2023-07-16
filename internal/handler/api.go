package handler

import (
	"errors"
	"hezzltask5/internal/models"
	projerrors "hezzltask5/internal/models/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createHandler(c *gin.Context) {
	q := c.Request.URL.Query()
	campaignId, err := strconv.Atoi(q.Get("campaignId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req models.CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Create(c, campaignId, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) updateHandler(c *gin.Context) {
	q := c.Request.URL.Query()
	campaignId, err := strconv.Atoi(q.Get("campaignId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(q.Get("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Update(c, campaignId, id, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, projerrors.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)

}

func (h *Handler) removeHandler(c *gin.Context) {
	q := c.Request.URL.Query()
	campaignId, err := strconv.Atoi(q.Get("campaignId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(q.Get("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//todo
	response, err := h.service.Delete(c, id, campaignId)
	if err != nil {
		if errors.Is(err, projerrors.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)

}

func (h *Handler) listHandler(c *gin.Context) {
	response, err := h.service.GetList(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
