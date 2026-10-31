package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"tac/gonic2/dal"

	"github.com/gin-gonic/gin"
)

type deptInput struct {
	Name string `json:"name" binding:"required"`
	Loc  string `json:"loc" binding:"required"`
}

// GET /depts
// Find all depts
func FindDepts(c *gin.Context) {
	var depts []dal.Dept
	if err := dal.DB.NewSelect().Model(&depts).Scan(c.Request.Context()); err != nil {
		log.Printf("[ERROR] get depts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[INFO] get all depts #%d", len(depts))
	c.JSON(http.StatusOK, gin.H{"data": depts})
}

// GET /depts/:id
// Find a dept
func FindDept(c *gin.Context) {
	// Get model if exist
	var dept dal.Dept
	id := c.Param("id")
	if err := dal.DB.NewSelect().Model(&dept).Where("id = ?", id).Scan(c.Request.Context()); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("[ERROR] get dept %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record " + id + " not found!"})
		return
	}
	log.Printf("[INFO] get dept %s", id)
	c.JSON(http.StatusOK, gin.H{"data": dept})
}

// POST /depts
// Create new dept
func CreateDept(c *gin.Context) {
	// Validate input
	var input deptInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create dept
	dept := dal.Dept{Dname: input.Name, Loc: input.Loc}
	if _, err := dal.DB.NewInsert().Model(&dept).Exec(c.Request.Context()); err != nil {
		log.Printf("[ERROR] create dept: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[INFO] create dept %s", dept.Dname)
	c.JSON(http.StatusOK, gin.H{"data": dept})
}

func DeleteDept(c *gin.Context) {
	// Get model if exist
	var dept dal.Dept
	id := c.Param("id")
	if err := dal.DB.NewSelect().Model(&dept).Where("id = ?", id).Scan(c.Request.Context()); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("[ERROR] get dept %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record " + id + " not found!"})
		return
	}

	if _, err := dal.DB.NewDelete().Model(&dept).Where("id = ?", id).Exec(c.Request.Context()); err != nil {
		log.Printf("[ERROR] delete dept %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[INFO] delete dept %s", id)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

func UpdateDept(c *gin.Context) {
	// Get model if exist
	var dept dal.Dept
	id := c.Param("id")
	if err := dal.DB.NewSelect().Model(&dept).Where("id = ?", id).Scan(c.Request.Context()); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("[ERROR] get dept %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record " + id + " not found!"})
		return
	}

	// Validate input
	var input deptInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dept.Dname = input.Name

	if _, err := dal.DB.NewUpdate().Model(&dept).Where("id = ?", id).Exec(c.Request.Context()); err != nil {
		log.Printf("[ERROR] update dept %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[INFO] update dept %s", id)
	c.JSON(http.StatusOK, gin.H{"data": dept})
}
