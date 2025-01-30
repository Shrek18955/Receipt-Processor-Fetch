SUMMARY:

main.go - Handles API requests and routing.
receipt.go - Handles validation and point calculation (which is open in the canvas).
go.mod - Manages dependencies.
go.sum - Ensures dependency integrity.

This Go-based receipt processor evaluates receipts and assigns reward points based on predefined rules. The implementation consists of two primary functions:

validateReceipt(receipt Receipt) bool
- Ensures that the receipt follows the correct format.
- Checks for required fields such as Retailer, PurchaseDate, PurchaseTime, Items, and Total.
- Validates formats using regular expressions (regexp), ensuring:
- Retailer names are properly formatted.
- The total price follows a valid decimal structure (xx.xx).
- The date and time match the expected format.

calculatePoints(receipt Receipt) int

- Retailer Name Bonus: 1 point per alphanumeric character in the retailer’s name.
- Round Dollar Bonus: +50 points if the total amount is a whole number ($50.00).
- Multiple of 0.25 Bonus: +25 points if the total is divisible by 0.25.
- Item Count Bonus: +5 points for every two items.
- Item Description Bonus: If the trimmed length of an item’s description is a multiple of 3, the item price is multiplied by 0.2, rounded up, and added as points.
- Odd Purchase Day Bonus: +6 points if the purchase date is an odd-numbered day.
- Afternoon Purchase Bonus: +10 points if the purchase was made between 2:00 PM - 4:00 PM.



- A receipt is validated using validateReceipt().
- If valid, calculatePoints() determines the reward points.
