package controller

import (
	"Project_5/main/model"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// HandleVerification godoc
// @Summary Verify user by code
// @Description Checks if the verification code matches and if the phone number exists in the database
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body model.VerificationRequest true "Verify true"
// @Success 200 {object} model.NumVerRes "Verification result"
// @Failure 400 {object} model.NumVerRes "Invalid request body"
// @Failure 500 {object} gin.H "Server error"
// @Router /handleVer [post]
func HandleVerification(c *gin.Context) {

	var req model.VerificationRequest

	// Fixed JSON field names to match the struct
	response := model.NumVerRes{IsVer: false, IsInData: false}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	const verCode = "123456" // This is a security risk; a real app would generate and store this code
	if req.Code != verCode {
		c.JSON(http.StatusOK, response)
		return
	}

	response.IsVer = true

	var id int
	err := DB.QueryRow("SELECT id FROM users WHERE number_phone = ?", req.Phone).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.IsInData = false
			c.JSON(http.StatusOK, response)
		} else {
			// General database error
			log.Printf("Database error in handleVerification: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	response.IsInData = true
	c.JSON(http.StatusOK, response)
} //success

// HandleSignUp godoc
// @Summary Register a new user
// @Description Creates a new user in the system
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body model.SignupJson true "Signup data"
// @Success 200 {object} model.StatusR "Signup successful"
// @Failure 400 {object} model.StatusR "Invalid request body"
// @Failure 500 {object} model.StatusR "Server error"
// @Router /signup [post]
func HandleSignUp(c *gin.Context) {
	var data model.SignupJson
	res := model.StatusR{Status: false}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, res)
		return
	}

	_, err := DB.Exec("INSERT INTO users (name, surname, email, birth_date, passport_id, number_phone) VALUES (?, ?, ?, ?, ?, ?)",
		data.Name, data.Surname, data.Email, data.BirthDate, data.Pid, data.Number)
	if err != nil {
		log.Println("Signup error:", err)
		// We can return a more descriptive error here based on the database error
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Status = true
	c.JSON(http.StatusOK, res)
} //success
