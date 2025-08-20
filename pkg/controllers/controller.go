package controllers

import (
	"log"

	"strings"

	errors "github.com/Conding-Student/study_template/pkg/models/errors"
	middleware "github.com/Conding-Student/study_template/pkg/utils/go-utils/database"
	fiberUtils "github.com/Conding-Student/study_template/pkg/utils/go-utils/fiber"
	passwordHashing "github.com/Conding-Student/study_template/pkg/utils/go-utils/passwordHashing"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	user := make(map[string]any)
	db := middleware.DBConn

	if err := c.BodyParser(&user); err != nil {
		return c.JSON(errors.ErrorModel{
			Message:   "Cannot parse JSON",
			IsSuccess: false,
			Error:     err,
		})
	}

	requiredFields := []string{"username", "email", "password", "first_name", "last_name", "age"}
	for _, field := range requiredFields {
		if user[field] == nil {
			return c.JSON(errors.ErrorModel{
				Message:   field + " is required",
				IsSuccess: false,
			})
		}
	}

	var age int
	// if val, ok := user["age"].(float64); ok {
	// 	age = int(val)
	// 	if age <= 0 {
	// 		return c.JSON(errors.ErrorModel{
	// 			Message:   "Age must be a positive number",
	// 			IsSuccess: false,
	// 		})
	// 	}
	// } else {
	// 	return c.JSON(errors.ErrorModel{
	// 		Message:   "Invalid age",
	// 		IsSuccess: false,
	// 	})
	// }

	email := strings.ToLower(strings.TrimSpace(user["email"].(string)))
	password := user["password"].(string) // keep password as is

	// Manual password complexity check
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	if len(password) >= 8 {
		for _, c := range password {
			switch {
			case 'A' <= c && c <= 'Z':
				hasUpper = true
			case 'a' <= c && c <= 'z':
				hasLower = true
			case '0' <= c && c <= '9':
				hasNumber = true
			case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:',.<>?/`~", c):
				hasSpecial = true
			}
		}
	}

	if !(hasUpper && hasLower && hasNumber && hasSpecial) {
		return c.JSON(errors.ErrorModel{
			Message:   "Password must be at least 8 characters long and include uppercase, lowercase, number, and special character",
			IsSuccess: false,
		})
	}

	hashedPassword, err := passwordHashing.HashPassword(password)
	if err != nil {
		return c.JSON(errors.ErrorModel{
			Message:   "Failed to hash password",
			IsSuccess: false,
			Error:     err,
		})
	}

	var feedback string
	err = db.Raw(
		"SELECT create_user($1, $2, $3, $4, $5, $6) AS feedback",
		user["username"],
		email,
		hashedPassword,
		user["first_name"],
		user["last_name"],
		age,
	).Scan(&feedback).Error
	if err != nil {
		log.Println("[DEBUG] Database error during user creation:", err)
		return c.JSON(errors.ErrorModel{
			Message:   "Database error",
			IsSuccess: false,
			Error:     err,
		})
	}

	return c.JSON(errors.ErrorModel{
		Message:   feedback,
		IsSuccess: true,
	})
}

func LoginUser(c *fiber.Ctx) error {
	// Parse request body
	input := make(map[string]string)
	if err := c.BodyParser(&input); err != nil {
		log.Println("[DEBUG] Failed to parse request body:", err)
		return c.JSON(errors.ErrorModel{
			Message:   "Cannot parse JSON",
			IsSuccess: false,
			Error:     err,
		})
	}

	// Get email and password, trim spaces
	email := strings.TrimSpace(input["email"])
	password := input["password"]

	log.Printf("[DEBUG] Login attempt received for email: '%s'\n", email)

	if email == "" {
		log.Println("[DEBUG] Email field is empty")
		return c.JSON(errors.ErrorModel{
			Message:   "Email is required",
			IsSuccess: false,
		})
	}
	if password == "" {
		log.Println("[DEBUG] Password field is empty")
		return c.JSON(errors.ErrorModel{
			Message:   "Password is required",
			IsSuccess: false,
		})
	}

	db := middleware.DBConn

	// Query user dynamically into a map
	user := make(map[string]interface{})
	err := db.Raw("SELECT * FROM get_user_by_email($1)", email).Scan(&user).Error
	if err != nil {
		// err will contain the message from PostgreSQL
		log.Println("[DEBUG] Database query error:", err)

		// Check if the error contains the specific message
		errMsg := err.Error()
		friendlyMsg := "Invalid email or password" // default

		if strings.Contains(errMsg, "No account exists with this email") {
			friendlyMsg = "No account exists with this email"
		}

		return c.JSON(errors.ErrorModel{
			Message:   friendlyMsg,
			IsSuccess: false,
			Error:     err,
		})
	}

	if user["user_id"] == nil || user["password_hash"] == nil {
		log.Printf("[DEBUG] No user found with email: '%s'\n", email)
		return c.JSON(errors.ErrorModel{
			Message:   "No user found with email",
			IsSuccess: false,
		})
	}

	passwordHash, ok := user["password_hash"].(string)
	if !ok {
		log.Printf("[DEBUG] Password hash type assertion failed for email: '%s'\n", email)
		return c.JSON(errors.ErrorModel{
			Message:   "Invalid email or password",
			IsSuccess: false,
		})
	}

	// Debug stored hash
	log.Printf("[DEBUG] Stored hash for '%s': %s\n", email, passwordHash)

	// Compare password
	if !passwordHashing.CheckPasswordHash(password, passwordHash) {
		log.Printf("[DEBUG] Password mismatch for user: '%s'\n", email)
		return c.JSON(errors.ErrorModel{
			Message:   "Invalid email or password",
			IsSuccess: false,
		})
	}

	log.Printf("[DEBUG] Login successful for user: '%s'\n", email)

	// Generate JWT token using fiberUtils
	fiberUtils.Ctx.New(c) // Copy context
	tokenPayload := map[string]interface{}{
		"user_id": user["user_id"],
		"email":   email,
	}

	tokenString, err := fiberUtils.GenerateJWTSignedString(tokenPayload)
	if err != nil {
		log.Println("[DEBUG] Token generation failed:", err)
		return c.JSON(errors.ErrorModel{
			Message:   "Login successful, but failed to generate token",
			IsSuccess: false,
			Error:     err,
		})
	}

	// Return token in response
	return c.JSON(map[string]interface{}{
		"message":   "Login successful",
		"isSuccess": true,
		"token":     tokenString,
	})
}

// REQUIRED JWT TOKEN
func GetPersonalDetails(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Copy context for fiberUtils
	fiberUtils.Ctx.New(c)

	// Extract user_id from JWT token
	claims := fiberUtils.GetJWTClaims()
	userID, ok := claims["user_id"].(float64) // JWT numbers come as float64
	if !ok {
		return c.JSON(errors.ErrorModel{
			Message:   "Invalid token: user_id not found",
			IsSuccess: false,
		})
	}

	// Query user personal details
	user := make(map[string]interface{})
	// get_personal_details_by_id is a function that returns user details by user_id
	err := db.Raw("SELECT * FROM get_personal_details_by_id($1)", userID).Scan(&user).Error

	if err != nil {
		return c.JSON(errors.ErrorModel{
			Message:   "Failed to fetch user details",
			IsSuccess: false,
			Error:     err,
		})
	}

	if user["user_id"] == nil {
		return c.JSON(errors.ErrorModel{
			Message:   "User not found",
			IsSuccess: false,
		})
	}

	// Return personal details
	return c.JSON(map[string]interface{}{
		"message":   "User details fetched successfully",
		"isSuccess": true,
		"data":      user,
	})
}
func UpdatePersonalDetails(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Parse request body into a map
	user := make(map[string]any)
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(errors.ErrorModel{
			Message:   "Cannot parse JSON",
			IsSuccess: false,
			Error:     err,
		})
	}

	// Extract user_id from JWT
	fiberUtils.Ctx.New(c)
	claims := fiberUtils.GetJWTClaims()
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(errors.ErrorModel{
			Message:   "Invalid token: user_id not found",
			IsSuccess: false,
		})
	}
	userID := int(userIDFloat)

	// Extract and sanitize fields
	username := ""
	if v, ok := user["username"].(string); ok {
		username = strings.TrimSpace(v)
	}

	email := ""
	if v, ok := user["email"].(string); ok {
		email = strings.ToLower(strings.TrimSpace(v))
	}

	firstName := ""
	if v, ok := user["first_name"].(string); ok {
		firstName = strings.TrimSpace(v)
	}

	lastName := ""
	if v, ok := user["last_name"].(string); ok {
		lastName = strings.TrimSpace(v)
	}

	var agePtr *int
	if v, ok := user["age"].(float64); ok {
		a := int(v)
		agePtr = &a
	}

	passwordPtr := (*string)(nil) // nil if password not provided
	if v, ok := user["password"].(string); ok && v != "" {
		password := v

		// Password complexity check
		var hasUpper, hasLower, hasNumber, hasSpecial bool
		if len(password) >= 8 {
			for _, c := range password {
				switch {
				case 'A' <= c && c <= 'Z':
					hasUpper = true
				case 'a' <= c && c <= 'z':
					hasLower = true
				case '0' <= c && c <= '9':
					hasNumber = true
				case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:',.<>?/`~", c):
					hasSpecial = true
				}
			}
		}

		if !(hasUpper && hasLower && hasNumber && hasSpecial) {
			return c.JSON(errors.ErrorModel{
				Message:   "Password must be at least 8 characters long and include uppercase, lowercase, number, and special character",
				IsSuccess: false,
			})
		}

		// Hash password
		hashed, err := passwordHashing.HashPassword(password)
		if err != nil {
			return c.JSON(errors.ErrorModel{
				Message:   "Failed to hash password",
				IsSuccess: false,
				Error:     err,
			})
		}
		passwordPtr = &hashed
	}

	// Call DB function: exactly 7 parameters
	var feedback string
	err := db.Raw(
		"SELECT update_user($1,$2,$3,$4,$5,$6,$7)",
		userID,
		username,
		email,
		passwordPtr, // pointer, nil will be NULL in DB
		firstName,
		lastName,
		agePtr,
	).Scan(&feedback).Error
	if err != nil {
		return c.JSON(errors.ErrorModel{
			Message:   "Database error",
			IsSuccess: false,
			Error:     err,
		})
	}

	return c.JSON(errors.ErrorModel{
		Message:   feedback,
		IsSuccess: true,
	})
}

// REQUIRED JWT TOKEN
func DeleteUser(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Copy context for fiberUtils
	fiberUtils.Ctx.New(c)

	// Extract user_id from JWT token
	claims := fiberUtils.GetJWTClaims()
	userID, ok := claims["user_id"].(float64) // JWT numbers come as float64
	if !ok {
		return c.JSON(errors.ErrorModel{
			Message:   "Invalid token: user_id not found",
			IsSuccess: false,
		})
	}

	// Execute the delete function in PostgreSQL
	err := db.Exec("SELECT delete_user($1)", int(userID)).Error
	if err != nil {
		return c.JSON(errors.ErrorModel{
			Message:   "Failed to delete user",
			IsSuccess: false,
			Error:     err,
		})
	}

	return c.JSON(map[string]interface{}{
		"message":   "User deleted successfully",
		"isSuccess": true,
	})
}
