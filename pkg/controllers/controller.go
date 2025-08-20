package controllers

import (
	"log"
	"regexp"
	"strings"

	errors "github.com/Conding-Student/study_template/pkg/models/errors"
	response "github.com/Conding-Student/study_template/pkg/models/response"
	middleware "github.com/Conding-Student/study_template/pkg/utils/go-utils/database"
	fiberUtils "github.com/Conding-Student/study_template/pkg/utils/go-utils/fiber"
	passwordHashing "github.com/Conding-Student/study_template/pkg/utils/go-utils/passwordHashing"
	"github.com/gofiber/fiber/v2"
)

// ========================= Create User =========================
func CreateUser(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Parse request body
	var user map[string]any
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid request body",
			Data: errors.ErrorModel{
				Message:   "Cannot parse JSON",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Extract fields
	username := user["username"].(string)
	firstName := user["first_name"].(string)
	lastName := user["last_name"].(string)
	ageVal := user["age"].(float64)
	email := strings.ToLower(strings.TrimSpace(user["email"].(string)))
	password := user["password"].(string)

	// Validate username, names, and age
	if validationErr := ValidateUserInput(username, firstName, lastName, ageVal); validationErr != nil {
		return c.JSON(validationErr)
	}

	// Validate email
	if !isEmailValid(email) {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Email format is not valid",
			Data: errors.ErrorModel{
				Message:   "Invalid email format",
				IsSuccess: false,
			},
		})
	}

	// Validate and hash password (registration requires password)
	hashedPassword, passwordErr := ValidateAndHashPassword(password, false)
	if passwordErr != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Password must be at least 8 characters long and include uppercase, lowercase, number, and special character",
			Data:    passwordErr,
		})
	}

	// Call DB function
	var feedback string
	if err := db.Raw(
		"SELECT create_user($1,$2,$3,$4,$5,$6) AS feedback",
		username,
		email,
		hashedPassword,
		firstName,
		lastName,
		ageVal,
	).Scan(&feedback).Error; err != nil {
		log.Println("[DEBUG] Database error:", err)
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Page has been broken",
			Data: errors.ErrorModel{
				Message:   "Database error during user creation",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Success response
	if feedback == "Email already exists" || feedback == "Username already exists" {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: feedback,
			Data: errors.ErrorModel{
				Message:   feedback,
				IsSuccess: false,
			},
		})
	}

	// Success case
	return c.JSON(response.ResponseModel{
		RetCode: "100",
		Message: "User created successfully",
		Data: errors.ErrorModel{
			Message:   feedback,
			IsSuccess: true,
		},
	})
}

// ========================= Login User =========================
func LoginUser(c *fiber.Ctx) error {
	input := make(map[string]string)
	if err := c.BodyParser(&input); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid request body",
			Data: errors.ErrorModel{
				Message:   "Cannot parse JSON",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	email := strings.TrimSpace(input["email"])
	password := input["password"]

	if email == "" || password == "" {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Email and password are required",
			Data: errors.ErrorModel{
				Message:   "Missing credentials",
				IsSuccess: false,
			},
		})
	}

	db := middleware.DBConn
	user := make(map[string]interface{})

	err := db.Raw("SELECT * FROM get_user_by_email($1)", email).Scan(&user).Error
	if err != nil || user["user_id"] == nil {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid email or password",
			Data: errors.ErrorModel{
				Message:   "Invalid input body",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	passwordHash, ok := user["password_hash"].(string)
	if !ok || !passwordHashing.CheckPasswordHash(password, passwordHash) {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid email or password",
			Data: errors.ErrorModel{
				Message:   "Invalid input body",
				IsSuccess: false,
			},
		})
	}

	// Generate JWT
	fiberUtils.Ctx.New(c)
	tokenPayload := map[string]interface{}{
		"user_id": user["user_id"],
		"email":   email,
	}
	token, err := fiberUtils.GenerateJWTSignedString(tokenPayload)
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Login succeeded, but the problem is on us",
			Data: errors.ErrorModel{
				Message:   "Token generation failed",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "100",
		Message: "Login successfully!",
		Data: map[string]any{
			"message":   "Login successful",
			"token":     token,
			"isSuccess": true,
		},
	})
}

// ========================= Get Personal Details =========================
// GetPersonalDetails handles fetching user details
func GetPersonalDetails(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Copy context for fiberUtils
	fiberUtils.Ctx.New(c)

	// Extract user_id from JWT token
	claims := fiberUtils.GetJWTClaims()
	userID, ok := claims["user_id"].(float64) // JWT numbers come as float64
	if !ok {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid token: user_id missing",
			Data: errors.ErrorModel{
				Message:   "Invalid token",
				IsSuccess: false,
			},
		})
	}

	// ✅ Call your Postgres function (RAW SQL)
	var user map[string]any
	err := db.Raw("SELECT * FROM get_personal_details_by_id($1)", int(userID)).Scan(&user).Error
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "404",
			Message: "User not found",
			Data: errors.ErrorModel{
				Message:   "User not found",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// ✅ Return user data + success info
	return c.JSON(response.ResponseModel{
		RetCode: "100",
		Message: "User fetched successfully",
		Data: map[string]any{
			"user": user,
			"info": errors.ErrorModel{
				Message:   "User fetched successfully",
				IsSuccess: true,
			},
		},
	})
}

// ========================= Update Personal Details =========================
// UpdatePersonalDetails updates user info using stored procedure
func UpdatePersonalDetails(c *fiber.Ctx) error {
	db := middleware.DBConn
	fiberUtils.Ctx.New(c)

	// Extract user_id from JWT token
	claims := fiberUtils.GetJWTClaims()
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Unauthorized",
			Data: errors.ErrorModel{
				Message:   "Missing or invalid user_id in token",
				IsSuccess: false,
			},
		})
	}

	// Parse request body
	var req map[string]any
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid request body",
			Data: errors.ErrorModel{
				Message:   err.Error(),
				IsSuccess: false,
			},
		})
	}

	// Extract fields from request
	username, _ := req["username"].(string)
	email, _ := req["email"].(string)
	password, _ := req["password"].(string) // raw password from request
	firstName, _ := req["first_name"].(string)
	lastName, _ := req["last_name"].(string)

	var age *int
	if val, ok := req["age"].(float64); ok {
		tmp := int(val)
		age = &tmp
	}

	//val;idate email
	if !isEmailValid(email) {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid email format",
			Data: errors.ErrorModel{
				Message:   "Email format is not valid",
				IsSuccess: false,
			},
		})
	}
	// Validate and hash password
	hashedPassword, passErr := ValidateAndHashPassword(password, true) // empty allowed
	if passErr != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Password must be at least 8 characters long and include uppercase, lowercase, number, and special character",
			Data:    passErr,
		})
	}

	// Call stored procedure
	var feedback string
	dbErr := db.Raw(
		"SELECT update_user($1, $2, $3, $4, $5, $6, $7)",
		int(userID), username, email, hashedPassword, firstName, lastName, age,
	).Scan(&feedback).Error

	if dbErr != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data: errors.ErrorModel{
				Message:   dbErr.Error(),
				IsSuccess: false,
			},
		})
	}

	// ✅ Determine success/failure from feedback
	success := strings.Contains(strings.ToLower(feedback), "success")

	return c.JSON(response.ResponseModel{
		RetCode: func() string {
			if success {
				return "100"
			}
			return "400"
		}(),
		Message: feedback, // direct message from function (e.g. "Email already exists")
		Data: map[string]any{
			"feedback": feedback,
			"user": map[string]any{
				"username":   username,
				"email":      email,
				"first_name": firstName,
				"last_name":  lastName,
				"age":        age,
			},
			"info": errors.ErrorModel{
				Message:   feedback,
				IsSuccess: success,
			},
		},
	})
}

// ========================= Delete User =========================
func DeleteUser(c *fiber.Ctx) error {
	db := middleware.DBConn
	fiberUtils.Ctx.New(c)

	claims := fiberUtils.GetJWTClaims()
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid token",
			Data: errors.ErrorModel{
				Message:   "Invalid token: user_id not found",
				IsSuccess: false,
			},
		})
	}

	err := db.Exec("SELECT delete_user($1)", int(userID)).Error
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data: errors.ErrorModel{
				Message:   "Failed to delete user",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "100",
		Message: "User deleted successfully",
		Data: errors.ErrorModel{
			Message:   "User deleted successfully",
			IsSuccess: true,
		},
	})
}

// ========================= Helpers =========================
// ===================== Helper: Validate Required Fields =====================
func ValidateUserInput(username, firstName, lastName string, age float64) *response.ResponseModel {
	if strings.ToLower(strings.TrimSpace(username)) == "" {
		return &response.ResponseModel{
			RetCode: "400",
			Message: "Username is required",
			Data: errors.ErrorModel{
				Message:   "Username cannot be empty",
				IsSuccess: false,
			},
		}
	}

	if strings.TrimSpace(firstName) == "" {
		return &response.ResponseModel{
			RetCode: "400",
			Message: "First name is required",
			Data: errors.ErrorModel{
				Message:   "First name cannot be empty",
				IsSuccess: false,
			},
		}
	}

	if strings.TrimSpace(lastName) == "" {
		return &response.ResponseModel{
			RetCode: "400",
			Message: "Last name is required",
			Data: errors.ErrorModel{
				Message:   "Last name cannot be empty",
				IsSuccess: false,
			},
		}
	}

	if age <= 0 {
		return &response.ResponseModel{
			RetCode: "400",
			Message: "Age is required",
			Data: errors.ErrorModel{
				Message:   "Age must be a positive number",
				IsSuccess: false,
			},
		}
	}

	return nil
}

func ValidateAndHashPassword(password string, isUpdate bool) (string, *errors.ErrorModel) {
	// If blank password
	if password == "" {
		if isUpdate {
			return "", nil // allow empty password on update
		} else {
			return "", &errors.ErrorModel{
				Message:   "Password cannot be empty",
				IsSuccess: false,
			}
		}
	}

	// Password complexity check
	if !isPasswordStrong(password) {
		return "", &errors.ErrorModel{
			Message:   "Password validation failed",
			IsSuccess: false,
		}
	}

	// Hash password
	hashed, err := passwordHashing.HashPassword(password)
	if err != nil {
		return "", &errors.ErrorModel{
			Message:   "Failed to hash password: " + err.Error(),
			IsSuccess: false,
			Error:     err,
		}
	}

	return hashed, nil
}

func isPasswordStrong(password string) bool {
	// Reject empty password immediately (spaces are allowed)

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	// Password must be at least 8 characters
	if len(password) < 8 {
		return false
	}

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

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// isEmailValid checks if email matches allowed domains
func isEmailValid(email string) bool {
	regex := `^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.(com|net|org|edu|gov)$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
