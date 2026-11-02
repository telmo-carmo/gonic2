package controllers

import (
	"database/sql"
	"net/http"
	"tac/gonic2/dal"

	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
)

func GetAllBonus(c *gin.Context) {
	myDb := dal.DB.DB

	log.Info().Msg("GetAllBonus")
	rows, err := myDb.Query("SELECT ename, job, sal, comm  FROM bonus")
	if err != nil {
		log.Warn().Err(err).Msg("Query error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	myBonus := []dal.Bonus{}

	for rows.Next() {
		var tempBn dal.Bonus
		rows.Scan(&tempBn.Ename, &tempBn.Job, &tempBn.Sal, &tempBn.Comm)
		log.Info().Msgf("Ename:%s, Job:%s, Sal:%d", tempBn.Ename, tempBn.Job, tempBn.Sal)
		myBonus = append(myBonus, tempBn)
	}
	rows.Close()

	if err = rows.Err(); err != nil {
		log.Warn().Err(err).Msg("Query rows error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, myBonus)
}

func GetOneBonus(c *gin.Context) {
	myDb := dal.DB.DB
	id := c.Params.ByName("id") // c.Param("id")
	uid := c.GetString("uid")
	log.Info().Msgf("GetOneBonus id: %s (%s)", id, uid)

	var tempBn dal.Bonus
	err := myDb.QueryRow("SELECT ename, job, sal,comm FROM bonus where ename = ?", id).Scan(&tempBn.Ename, &tempBn.Job, &tempBn.Sal, &tempBn.Comm)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
		log.Warn().Err(err).Msg("Query row error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		log.Info().Msgf("Ename:%s, Job:%s, Sal:%d", tempBn.Ename, tempBn.Job, tempBn.Sal)
	}
	c.JSON(http.StatusOK, tempBn)
}

func CreateBonus(c *gin.Context) {
	// Validate input
	var bn dal.Bonus
	if err := c.ShouldBindJSON(&bn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	myDb := dal.DB.DB
	_, err := myDb.Exec("INSERT INTO bonus (ename, job, sal,comm) VALUES (?,?,?,?)", bn.Ename, bn.Job, bn.Sal, bn.Comm)
	if err != nil {
		log.Error().Err(err).Msg("insert bonus")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bn)
}

func UpdateBonus(c *gin.Context) {
	id := c.Params.ByName("id")
	// Validate input
	var bn dal.Bonus
	if err := c.ShouldBindJSON(&bn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if id != bn.Ename {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input payload and id"})
		return
	}
	myDb := dal.DB.DB
	_, err := myDb.Exec("UPDATE bonus SET job=?, sal=?,comm=?  WHERE ename = ?", bn.Job, bn.Sal, bn.Comm, bn.Ename)
	if err != nil {
		log.Error().Err(err).Msg("update bonus")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bn)
}

func DeleteBonus(c *gin.Context) {
	id := c.Params.ByName("id")
	log.Info().Msgf("delete %s", id)
	myDb := dal.DB.DB
	var bn dal.Bonus
	err := myDb.QueryRow("SELECT ename, job, sal,comm FROM bonus where ename = ?", id).Scan(&bn.Ename, &bn.Job, &bn.Sal, &bn.Comm)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
		log.Warn().Err(err).Msg("Query row error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		log.Info().Msgf("Ename:%s, Job:%s, Sal:%d", bn.Ename, bn.Job, bn.Sal)
	}

	_, err = myDb.Exec("DELETE FROM  bonus  WHERE ename = ?", bn.Ename)
	if err != nil {
		log.Error().Err(err).Msg("delete bonus")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data:": id})
}
