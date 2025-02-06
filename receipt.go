package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

func validateReceipt(receipt Receipt) bool {
	if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || len(receipt.Items) == 0 || receipt.Total == "" {
		return false
	}
	// Validate retailer name format
	regRetailer := regexp.MustCompile(`^[\w\s\-&]+$`)
	if !regRetailer.MatchString(receipt.Retailer) {
		return false
	}
	// Validate total amount format (must be in the format of xx.xx)
	regTotal := regexp.MustCompile(`^\d+\.\d{2}$`)
	if !regTotal.MatchString(receipt.Total) {
		return false
	}
	// Validate purchase date format
	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		return false
	}
	// Validate purchase time format
	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		return false
	}
	return true
}

func calculatePoints(receipt Receipt) int {
	points := 0

	// 1 point for every alphanumeric character in the retailer name
	retailerPoints := 0
	for _, c := range receipt.Retailer {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			retailerPoints++
		}
	}
	points += retailerPoints
	fmt.Println("Retailer Name Points:", retailerPoints)

	// Check if total amount is a round dollar amount
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == float64(int(total)) {
		points += 50
		fmt.Println("Round Dollar Bonus: +50")
	}

	// Check if total amount is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
		fmt.Println("Multiple of 0.25 Bonus: +25")
	}

	// 5 points for every 2 items on the receipt
	itemPoints := (len(receipt.Items) / 2) * 5
	points += itemPoints
	fmt.Println("Item Count Bonus:", itemPoints)

	// Check if item description length is a multiple of 3 and calculate bonus points
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			itemDescPoints := int(math.Ceil(price * 0.2))
			points += itemDescPoints
			fmt.Println("Item Description Bonus:", itemDescPoints, "for", item.ShortDescription)
		}
	}

	// Check if the purchase date's day is odd
	dateParts := strings.Split(receipt.PurchaseDate, "-")
	if day, err := strconv.Atoi(dateParts[2]); err == nil && day%2 != 0 {
		points += 6
		fmt.Println("Odd Day Bonus: +6")
	}

	// Check if the purchase time is between 2:00 PM - 3:59 PM
	if t, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		hour, minute := t.Hour(), t.Minute()
		if hour == 14 || (hour == 15 && minute <= 59) {
			points += 10
			fmt.Println("Afternoon Time Bonus: +10")
		}
	}

	fmt.Println("Total Points:", points)
	return points
}
