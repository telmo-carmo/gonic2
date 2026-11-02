package controllers

import (
	"database/sql"
	"net/http"
	"tac/gonic2/dal"

	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
)

type bookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

// GET /books
// Find all books
func FindBooks(c *gin.Context) {
	var books []dal.Book
	if err := dal.DB.NewSelect().Model(&books).Scan(c.Request.Context()); err != nil {
		log.Error().Err(err).Msg("get books")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info().Msgf("get all books #%d", len(books))
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// GET /books/:id
// Find a book
func FindBook(c *gin.Context) {
	// Get model if exist
	var book dal.Book
	id := c.Param("id")
	if err := dal.DB.NewSelect().Model(&book).Where("id = ?", id).Scan(c.Request.Context()); err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msgf("get book %s", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record " + id + " not found!"})
		return
	}
	log.Info().Msgf("get book %s", id)
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// POST /books
// Create new book
func CreateBook(c *gin.Context) {
	// Validate input
	var input bookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	book := dal.Book{Title: input.Title, Author: input.Author}
	if _, err := dal.DB.NewInsert().Model(&book).Exec(c.Request.Context()); err != nil {
		log.Error().Err(err).Msg("create book")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info().Msgf("created book '%v'", book)
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// PUT /books/:id
// Update a book
func UpdateBook(c *gin.Context) {
	// Get model if exist
	var book dal.Book
	id := c.Param("id")
	if err := dal.DB.NewSelect().Model(&book).Where("id = ?", id).Scan(c.Request.Context()); err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msgf("get book %s for update", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record " + id + " not found!"})
		return
	}

	// Validate input
	var input bookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book.Title = input.Title
	book.Author = input.Author
	log.Info().Msgf("updating book '%v'", book)
	if _, err := dal.DB.NewUpdate().Model(&book).WherePK().Column("title", "author").Exec(c.Request.Context()); err != nil {
		log.Error().Err(err).Msgf("update book %s", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// DELETE /books/:id
// Delete a book
func DeleteBook(c *gin.Context) {
	// Get model if exist
	var book dal.Book
	id := c.Param("id")
	if err := dal.DB.NewSelect().Model(&book).Where("id = ?", id).Scan(c.Request.Context()); err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msgf("get book %s for delete", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	log.Info().Msgf("deleting book '%v'", book)
	if _, err := dal.DB.NewDelete().Model(&book).WherePK().Exec(c.Request.Context()); err != nil {
		log.Error().Err(err).Msgf("delete book %s", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
