package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Income struct {
	MonthlyIncome   int `json:"monthly_income"`
	WorkedMonth     int `json:"worked_month"`
	Bonus           int `json:"bonus"`
	FreelanceIncome int `json:"freelance_income"`
}

type Deduction struct {
	SpouseDeduction               bool `json:"spouse_deduction"`
	PregnancyExpenses             int  `json:"pregnancy_expenses"`
	SecondaryCities               int  `json:"secondary_cities"`
	ShopdeeMeekhun                int  `json:"shopdee_meekhun"`
	HomeLoanInterest              int  `json:"home_loan_interest"`
	PurchaseOtopProducts          int  `json:"purchase_otop_products"`
	PurchaseCommunityEnterprises  int  `json:"purchase_community_enterprises"`
	PurchaseFromSocialEnterprises int  `json:"purchase_from_social_enterprises"`
	PurchaseWithVatETax           int  `json:"purchase_with_vat_etax"`
	PurchaseWithEReceipt          int  `json:"purchase_with_e_receipt"`
}

type TaxRequest struct {
	Income    Income    `json:"income"`
	Deduction Deduction `json:"deduction"`
}

type TaxResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TotalIncome    int `json:"total_income"`
		TotalDeduction int `json:"total_deduction"`
		TotalTax       int `json:"total_tax"`
		Refund         int `json:"refund"`
	} `json:"data"`
}

type TaxResult struct {
	gorm.Model
	TotalIncome    int
	TotalDeduction int
	TotalTax       int
	Refund         int
}

// var db *gorm.DB

func main() {
	// var err error

	// dsn := "root:photo0877071818@TH@tcp(127.0.0.1:3306)/tax_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal("Failed to connect to database:", err)
	// }

	// Create TaxResult table if not exists
	// db.AutoMigrate(&TaxResult{})

	// Setup Echo server
	e := echo.New()
	e.POST("/tax-calculation", calculateTaxHandler)

	log.Println("Server started on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func calculateTaxHandler(c echo.Context) error {
	var req TaxRequest

	// Parse JSON request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	totalIncome := (req.Income.MonthlyIncome * req.Income.WorkedMonth) +
		req.Income.Bonus + req.Income.FreelanceIncome

	totalDeduction := calculateTotalDeduction(req.Deduction)
	netIncome := totalIncome - totalDeduction

	totalTax := 0
	if netIncome > 0 {
		totalTax = calculateTax(netIncome)
	}

	taxPaid := 0
	refund := taxPaid - totalTax
	if refund < 0 {
		refund = 0
	}
	// Save result to database
	// db.Create(&TaxResult{
	// 	TotalIncome:    totalIncome,
	// 	TotalDeduction: totalDeduction,
	// 	TotalTax:       totalTax,
	// 	Refund:         refund,
	// })

	fmt.Println("Total Income:", totalIncome)

	var resp TaxResponse
	resp.Code = 1000
	resp.Message = "success"
	resp.Data.TotalIncome = totalIncome
	resp.Data.TotalDeduction = totalDeduction
	resp.Data.TotalTax = totalTax
	resp.Data.Refund = refund

	return c.JSON(http.StatusOK, resp)
}

// Sum up all applicable deductions
func calculateTotalDeduction(d Deduction) int {
	total := 60000 // Personal deduction

	if d.SpouseDeduction {
		total += 60000
	}

	// Cap and add each deduction
	total += min(d.PregnancyExpenses, 60000)
	total += min(d.SecondaryCities, 15000)
	total += min(d.ShopdeeMeekhun, 50000)
	total += min(d.HomeLoanInterest, 100000)
	total += min(d.PurchaseOtopProducts, 20000)
	total += min(d.PurchaseCommunityEnterprises, 20000)
	total += min(d.PurchaseFromSocialEnterprises, 20000)
	total += min(d.PurchaseWithVatETax, 30000)
	total += min(d.PurchaseWithEReceipt, 30000)

	return total
}

// Apply tax rates to net income
func calculateTax(income int) int {
	brackets := []struct {
		limit int
		rate  float64
	}{
		{150000, 0.00}, // No tax for income up to 150,000
		{300000, 0.05}, // 5% for income between 150,001 and 300,000
		{500000, 0.10},
		{750000, 0.15},
		{1000000, 0.20},
		{2000000, 0.25},
		{5000000, 0.30},
		{1<<63 - 1, 0.35},
	}

	tax := 0.0
	prevLimit := 0

	for _, b := range brackets {
		if income <= b.limit {
			tax += float64(income-prevLimit) * b.rate
			break
		}
		tax += float64(b.limit-prevLimit) * b.rate
		prevLimit = b.limit
	}

	return int(tax)
}

// Utility function to return minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
