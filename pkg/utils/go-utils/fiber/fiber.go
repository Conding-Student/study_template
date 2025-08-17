// Package fiber provides utility functions for gofiber v2, jwt-go
// With additional validation functions, sending JSON response and parsing request bodies, getting JWT claims
package fiber

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

// Context GoFiber Context
type Context struct {
	c *fiber.Ctx
}

// JWTConfig configuration for JWT
type JWTConfig struct {
	Duration     time.Duration
	CookieMaxAge int
	SetCookies   bool
	SecretKey    []byte
}

// Ctx Context to be initiated by the New function
var Ctx Context
var jwtConfig JWTConfig

// New Copies GoFiber context as new current context
func (ctx *Context) New(c *fiber.Ctx) {
	Ctx = Context{
		c: c,
	}
}

type AccessGranted struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Token   string `json:"token"`
}

// Message struct for GoFiber context response
type Message struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Message struct for GoFiber context response
type MessageInterface struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

// ParseBody Parses the request body from the copied current context
func ParseBody(in interface{}) error {
	err := Ctx.c.BodyParser(in)

	if err != nil {
		LogError(err)
		return Ctx.c.Status(503).SendString(err.Error())
	}

	return err
}

// GetParamValue Gets the parameter value from the copied current context
func GetParamValue(param string, message string) string {
	paramValue := Ctx.c.Params(param)

	if paramValue == "" {
		err := SendJSONMessage(message, false, 400)
		LogError(err)
	}

	return paramValue
}

// SendJSONMessage Sends JSON Message with HTTP Status code to current context
func SendJSONMessageInterface(data interface{}, isSuccess bool, httpStatusCode int) error {
	status := "failed"

	if isSuccess {
		status = "success"
	}

	return Ctx.c.Status(httpStatusCode).JSON(MessageInterface{
		Status:     status,
		StatusCode: httpStatusCode,
		Data:       data,
	})
}

// SendJSONMessage Sends JSON Message with HTTP Status code to current context
func SendJSONMessage(message string, isSuccess bool, httpStatusCode int) error {
	status := "failed"

	if isSuccess {
		status = "success"
	}

	return Ctx.c.Status(httpStatusCode).JSON(Message{
		Message: message,
		Status:  status,
	})
}

func AccessGrantedJSONMessage(message string, token string, isSuccess bool, httpStatusCode int) error {
	status := "failed"

	if isSuccess {
		status = "success"
	}

	return Ctx.c.Status(httpStatusCode).JSON(AccessGranted{
		Message: message,
		Status:  status,
		Token:   token,
	})
}

func AccessGrantedResponse(message string, token string) error {
	err := AccessGrantedJSONMessage(message, token, true, 200)
	LogError(err)
	return err
}

// SendSuccessResponse Wrapper function for SendJSONMessage of 200 Success
func SendSuccessResponse(message string) error {
	err := SendJSONMessage(message, true, 200)
	LogError(err)
	return err
}

// SendBadRequestResponse Wrapper function for SendJSONMessage of 400 Bad request
func SendBadRequestResponse(message string) error {
	err := SendJSONMessage(message, false, 400)
	LogError(err)
	return err
}

// ValidateField Validation of strings and return if valid based on specification and error message if invalid
func ValidateField(field, title string, isMandatory bool, max, min int, format string) (ok bool, message string) {
	ok = true

	if !isMandatory {
		return
	}

	if len(field) == 0 {
		message += fmt.Sprintf("'%s' cannot be empty.", title)
		_ = SendBadRequestResponse(message)
		ok = false
	} else {
		switch format {
		case "S":
			if len(field) > 2 {
				message += fmt.Sprintf("The length of '%s' cannot be greater than 2.", title)
				_ = SendBadRequestResponse(message)
				ok = false
			}
		case "N":
			if _, err := strconv.Atoi(field); err != nil {
				message += fmt.Sprintf("'%s' should only contain numbers.", title)
				_ = SendBadRequestResponse(message)
				ok = false
			}

			fallthrough
		case "ANS":
			cflOK, cflMessage := CheckFieldLength(field, title, max, min)

			if !cflOK {
				ok = false
				message += "\n" + cflMessage
			}
		}
	}

	return ok, message
}

// CheckFieldLength Validation of strings' length and return if valid based on maximum and minimum length specified and error message if invalid
func CheckFieldLength(field, title string, max, min int) (ok bool, message string) {
	ok = true

	if len(field) > max {
		message += fmt.Sprintf("The length of '%s' cannot be greater than %d.", title, max)
		_ = SendBadRequestResponse(message)
		ok = false
	}

	if len(field) < min {
		message += fmt.Sprintf("The length of '%s' cannot be less than %d.", title, min)
		_ = SendBadRequestResponse(message)
		ok = false
	}

	return
}

// GetJWTClaims Get User JWT claims of the current context
func GetJWTClaims() jwt.MapClaims {
	userToken := Ctx.c.Locals("user").(*jwt.Token)
	return userToken.Claims.(jwt.MapClaims)
}

// GetJWTClaim Wrapper function for getting a JWT claim by key
func GetJWTClaim(key string) map[string]interface{} {
	return GetJWTClaims()[key].(map[string]interface{})
}

// GetJSONFieldValues Returns a map of JSON keys and values of a struct
func GetJSONFieldValues(q interface{}) map[string]string {
	val := reflect.ValueOf(q).Elem()
	fields := make(map[string]string)

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		fields[tag.Get("json")] = val.Field(i).String()
	}

	return fields
}

// ValidateJSONField Wrapper function for JSON field validation of a struct
func ValidateJSONField(q interface{}, tag string, isMandatory bool, max, min int, format string) (bool, string) {
	return ValidateField(GetJSONFieldValues(q)[tag], tag, isMandatory, max, min, format)
}

// LogError Logs errors
func LogError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

// AuthenticationMiddleware ...
func AuthenticationMiddleware(j JWTConfig) func(*fiber.Ctx) error {
	jwtConfig = j
	ret := jwtware.New(jwtware.Config{
		SigningKey: jwtConfig.SecretKey,
	})

	return ret
}

// GenerateJWTSignedString ...
func GenerateJWTSignedString(claimMaps fiber.Map) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(jwtConfig.Duration).Unix()

	for key, value := range claimMaps {
		claims[key] = value
	}

	t, err := token.SignedString(jwtConfig.SecretKey)

	if jwtConfig.SetCookies && err == nil {
		Ctx.c.Cookie(&fiber.Cookie{
			Name:   "token",
			Value:  t,
			MaxAge: jwtConfig.CookieMaxAge,
		})
	}

	return t, err
}

// GetJWTClaimOfType ...
func GetJWTClaimOfType(key string, valueType interface{}) error {
	userInfoJSON, err := json.Marshal(GetJWTClaim(key))

	if err == nil {
		err = json.Unmarshal(userInfoJSON, &valueType)
	}

	return err
}
