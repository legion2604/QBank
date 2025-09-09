package controller

import (
	"Project_5/main/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func RemoveDept(c *gin.Context) {
	var deptstruct model.RemDeptStruct
	if err := c.ShouldBindJSON(&deptstruct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result, err := DB.Exec("DELETE FROM loans WHERE id = ?", deptstruct.Id)
	if err != nil {
		log.Printf("Database delete error in remDept: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete loan"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("RowsAffected error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete loan"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "error": "Loan not found"})
		return
	}

	c.JSON(http.StatusOK, model.StatusR{Status: true})
} //success

func AddDept(c *gin.Context) {
	var deptstruct model.AddDeptStruct
	if err := c.ShouldBindJSON(&deptstruct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	query := `
	INSERT INTO loans (lender_id, borrower_id, amount, interest_rate, created_at, due_date, status)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := DB.Exec(
		query,
		deptstruct.LenderId,
		deptstruct.BorrowerId,
		deptstruct.Amount,
		deptstruct.Interest,
		time.Now().UTC(),
		deptstruct.DueDate,
		deptstruct.Status,
	)
	if err != nil {
		log.Printf("addDept error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add loan"})
		return
	}

	c.JSON(http.StatusOK, model.StatusR{Status: true})
} //success
