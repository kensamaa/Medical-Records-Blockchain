package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/kensamaa/blockchain-medical-records/rest-api/services"

	"github.com/gin-gonic/gin"
)

// CreateRecord handles POST /api/records requests.
func CreateRecord(c *gin.Context) {
	var record map[string]interface{}
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal record"})
		return
	}

	if err := services.CreateRecord(string(recordJSON)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record created successfully"})
}

// GetRecord handles GET /api/records/:id requests.
func GetRecord(c *gin.Context) {
	id := c.Param("id")
	record, err := services.GetRecord(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, record)
}

// UpdateRecord handles PUT /api/records/:id requests.
func UpdateRecord(c *gin.Context) {
	var record map[string]interface{}
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal record"})
		return
	}

	if err := services.UpdateRecord(string(recordJSON)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
}
