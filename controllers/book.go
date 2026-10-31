package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"tac/gonic2/dal"

	"github.com/gin-gonic/gin"
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
		log.Printf("[ERROR] get books: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[INFO] get all books #%d", len(books))
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
			log.Printf("[ERROR] get book %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record " + id + " not found!"})
		return
	}
	log.Printf("[INFO] get book %s", id)
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
		log.Printf("[ERROR] create book: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[INFO] created book '%v'", book)
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
			log.Printf("[ERROR] get book %s for update: %v", id, err)
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
	log.Printf("[INFO] updating book '%v'", book)
	if _, err := dal.DB.NewUpdate().Model(&book).WherePK().Column("title", "author").Exec(c.Request.Context()); err != nil {
		log.Printf("[ERROR] update book %s: %v", id, err)
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
			log.Printf("[ERROR] get book %s for delete: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	log.Printf("[INFO] deleting book '%v'", book)
	if _, err := dal.DB.NewDelete().Model(&book).WherePK().Exec(c.Request.Context()); err != nil {
		log.Printf("[ERROR] delete book %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
