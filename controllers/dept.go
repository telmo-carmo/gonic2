package controllers

import (
	"database/sql"
	"net/http"
	"tac/gonic2/dal"

	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
)

type deptInput struct {
	Deptno int    `json:"deptno" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Loc    string `json:"loc" binding:"required"`
}

// GET /depts
// Find all depts
func FindDepts(c *gin.Context) {
	var depts []dal.Dept
	if err := dal.DB.NewSelect().Model(&depts).Scan(c.Request.Context()); err != nil {
		log.Error().Err(err).Msg("get depts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info().Msgf("get all depts #%d", len(depts))
	c.JSON(http.StatusOK, gin.H{"data": depts})
}

// GET /depts/:deptno
// Find a dept
func FindDept(c *gin.Context) {
	// Get model if exist
	var dept dal.Dept
	deptno := c.Param("deptno")
	if err := dal.DB.NewSelect().Model(&dept).Where("deptno = ?", deptno).Scan(c.Request.Context()); err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msgf("get dept %s", deptno)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record " + deptno + " not found!"})
		return
	}
	log.Info().Msgf("get dept %s", deptno)
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
	dept := dal.Dept{Name: input.Name, Loc: input.Loc}
	if _, err := dal.DB.NewInsert().Model(&dept).Exec(c.Request.Context()); err != nil {
		log.Error().Err(err).Msg("create dept")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info().Msgf("create dept %d", dept.Deptno)
	c.JSON(http.StatusOK, gin.H{"data": dept})
}

func DeleteDept(c *gin.Context) {
	// Get model if exist
	var dept dal.Dept
	deptno := c.Param("deptno")
	if err := dal.DB.NewSelect().Model(&dept).Where("deptno = ?", deptno).Scan(c.Request.Context()); err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msgf("get dept %s", deptno)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record " + deptno + " not found!"})
		return
	}

	if _, err := dal.DB.NewDelete().Model(&dept).WherePK().Exec(c.Request.Context()); err != nil {
		log.Error().Err(err).Msgf("delete dept %s", deptno)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info().Msgf("delete dept %s", deptno)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

func UpdateDept(c *gin.Context) {
	// Get model if exist
	var dept dal.Dept
	deptno := c.Param("deptno")
	if err := dal.DB.NewSelect().Model(&dept).Where("deptno = ?", deptno).Scan(c.Request.Context()); err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msgf("get dept %s", deptno)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record " + deptno + " not found!"})
		return
	}

	// Validate input
	var input deptInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dept.Name = input.Name
	dept.Loc = input.Loc

	if _, err := dal.DB.NewUpdate().Model(&dept).WherePK().Column("name").Column("loc").Exec(c.Request.Context()); err != nil {
		log.Error().Err(err).Msgf("update dept %s", deptno)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info().Msgf("update dept %s", deptno)
	c.JSON(http.StatusOK, gin.H{"data": dept})
}
