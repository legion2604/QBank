package model

import (
	"database/sql"
	"time"
)

type MonthSummary struct {
	TotalOwedToUser int `json:"total_owed_to_user"`
	TotalOwedByUser int `json:"total_owed_by_user"`
}

type StringData struct {
	Text string `json:"text"`
}
type DayDebts struct {
	Date   string   `json:"date"`
	ToMe   []string `json:"to_me"`
	FromMe []string `json:"from_me"`
}

type StatusR struct {
	Status bool `json:"status"`
}

type SignupJson struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	BirthDate string `json:"birth_date"`
	Pid       string `json:"pid"`
	Number    string `json:"number_phone"`
}

type NumVer struct {
	Number string `json:"number"`
	Code   string `json:"code"`
}

type NumVerRes struct {
	IsVer    bool `json:"is_ver"`
	IsInData bool `json:"is_in_data"`
}

type SigninForm struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// report structs
type monthSumReq struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	ID    string `json:"id"`
}

type monthSum struct {
	Principal int `json:"principal"`
	Interest  int `json:"interest"`
}

// data from users
type UserInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	BirthDate   string `json:"birth_date"`
	NumberPhone string `json:"number_phone"`
	PassportID  string `json:"passport_id"`
	CreatedAt   string `json:"created_at"`
}

// Dept
type userDept struct {
	Id           string    `json:"id"`
	LenderId     string    `json:"lender_id"`
	BorrowerId   string    `json:"borrower_id"`
	Amount       string    `json:"amount"`
	InterestRate string    `json:"interest_rate"`
	CreatedAt    time.Time `json:"created_at"`
	DueDate      string    `json:"due_date"`
	Status       string    `json:"status"`
}

type UserLoans struct {
	Id           int            `json:"id"`
	LenderId     sql.NullInt64  `json:"lender_id"`
	BorrowedId   sql.NullInt64  `json:"borrower_id"`
	Amount       int            `json:"amount"`
	InterestRate int            `json:"interest_rate"`
	CreatedAt    sql.NullString `json:"created_at"`
	DueDate      sql.NullString `json:"due_date"`
	Status       sql.NullString `json:"status"`
}

type AddDeptStruct struct {
	LenderId   int       `json:"lender_id"`
	BorrowerId int       `json:"borrower_id"`
	Amount     int       `json:"amount"`
	DueDate    time.Time `json:"due_date"` // Gin будет парсить дату автоматически
	Interest   int       `json:"interest"`
	Status     string    `json:"status"`
}

type RemDeptStruct struct {
	Id int `json:"id"`
}

type Payment struct {
	Id     int            `json:"id"`
	LoanId int            `json:"loan_id"`
	Amount int            `json:"amount"`
	Type   string         `json:"type"`
	PaidAt sql.NullString `json:"paid_at"`
}
