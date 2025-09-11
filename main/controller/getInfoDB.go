package controller

import (
	"Project_5/main/model"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetLoansUsers godoc
// @Summary Get all loans for a user
// @Description Returns all loans where the user is either lender or borrower
// @Tags Loans
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {array} model.UserLoans
// @Failure 400 {object} model.StatusR
// @Failure 500 {object} model.StatusR
// @Security BearerAuth
// @Router /loans/getAllLoans [get]
func GetLoansUsers(c *gin.Context) {

	var arr []model.UserLoans
	idStr := c.Query("id")
	num, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	// Выбираем все кредиты и займы, где пользователь либо lender, либо borrower
	rows, err := DB.Query(`
		SELECT id, lender_id, borrower_id, amount, interest_rate, created_at, due_date, status 
		FROM loans 
		WHERE lender_id = ? OR borrower_id = ?`,
		num, num)
	if err != nil {
		log.Printf("Database query error in getLoansUsers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var res model.UserLoans
		if err := rows.Scan(
			&res.Id,
			&res.LenderId,
			&res.BorrowedId,
			&res.Amount,
			&res.InterestRate,
			&res.CreatedAt,
			&res.DueDate,
			&res.Status,
		); err != nil {
			log.Printf("Scan error in getLoansUsers: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process database results"})
			return
		}
		arr = append(arr, res)
	}

	c.JSON(http.StatusOK, arr)
} // success

// GetUser godoc
// @Summary Get user by passport ID
// @Description Returns detailed information about a user by passport ID
// @Tags Loans
// @Accept json
// @Produce json
// @Param pid query string true "Passport ID"
// @Success 200 {object} model.UserInfo
// @Failure 404 {object} model.StatusR
// @Failure 500 {object} model.StatusR
// @Security BearerAuth
// @Router /user [get]
func GetUser(c *gin.Context) {
	id := c.Query("pid")

	var req model.UserInfo
	err := DB.QueryRow("SELECT id, name, surname, email, birth_date, number_phone, passport_id, created_at FROM users WHERE passport_id = ?", id).
		Scan(&req.Id, &req.Name, &req.Surname, &req.Email,
			&req.BirthDate, &req.NumberPhone, &req.PassportID, &req.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			log.Printf("Database query error in getUser: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, req)
} // success

// GetHistory godoc
// @Summary Get payment history for a user
// @Description Returns all payments where the user is lender or borrower
// @Tags Loans
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {array} model.Payment
// @Failure 400 {object} model.StatusR
// @Failure 500 {object} model.StatusR
// @Security BearerAuth
// @Router /reports/getHistory [get]
func GetHistory(c *gin.Context) {
	idStr := c.Query("id")
	num, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	rows, err := DB.Query(`
      SELECT p.id, p.loan_id, p.amount, p.type, p.paid_at
      FROM payments p
      JOIN loans l ON p.loan_id = l.id AND (l.lender_id = ? OR l.borrower_id = ?)  
    `, num, num)
	if err != nil {
		log.Printf("Database query error in getHistory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(&p.Id, &p.LoanId, &p.Amount, &p.Type, &p.PaidAt); err != nil {
			log.Printf("Scan error in getHistory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process database results"})
			return
		}
		payments = append(payments, p)
	}

	c.JSON(http.StatusOK, payments)
} // success

// GetMonthSummary godoc
// @Summary Get monthly loan summary for a user
// @Description Returns the total owed to and by the user for a specific month
// @Tags Loans
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} model.MonthSummary
// @Failure 400 {object} model.StatusR
// @Failure 500 {object} model.StatusR
// @Security BearerAuth
// @Router /reports [get]
func GetMonthSummary(c *gin.Context) {
	// Retrieve the user's ID from the request's query parameters.
	// This is less secure than using a token but fulfills the request.
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query parameter is required"})
		return
	}

	// Get year and month from the request's query parameters.
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month query parameters are required"})
		return
	}

	// Convert year and month strings to integers.
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format"})
		return
	}

	// Create a time range for the specified month.
	// This is more database-agnostic than using YEAR() and MONTH() functions.
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	// SQL query to calculate the sums.
	// It uses conditional aggregation to sum amounts based on whether the user
	// is the lender or the borrower. COALESCE ensures we get 0 instead of NULL if no loans are found.
	query := `
        SELECT
            COALESCE(SUM(CASE WHEN lender_id = ? THEN amount ELSE 0 END), 0) AS total_owed_to_user,
            COALESCE(SUM(CASE WHEN borrower_id = ? THEN amount ELSE 0 END), 0) AS total_owed_by_user
        FROM loans
        WHERE (lender_id = ? OR borrower_id = ?)
        AND due_date >= ? AND due_date < ?
    `

	var summary model.MonthSummary
	// Execute the query and scan the results directly into the summary struct.
	err = DB.QueryRow(query, userId, userId, userId, userId, startDate, endDate).Scan(&summary.TotalOwedToUser, &summary.TotalOwedByUser)
	if err != nil {
		// sql.ErrNoRows is not a critical error here; it just means no debts for that month.
		// We can safely return a zeroed summary.
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, model.MonthSummary{TotalOwedToUser: 0, TotalOwedByUser: 0})
		} else {
			log.Printf("Database error in getMonthSummary: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Return the summary data as JSON.
	c.JSON(http.StatusOK, summary)
} //success
