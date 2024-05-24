package controllers

import (
	"gofiber-marketplace/src/helpers"
	"gofiber-marketplace/src/middlewares"
	"gofiber-marketplace/src/models"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *fiber.Ctx) error {
	var register models.Register
	if err := c.BodyParser(&register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	user := middlewares.XSSMiddleware(&register).(*models.Register)
	if authErrors := helpers.PasswordValidation(user.Password, helpers.StructValidation(user)); len(authErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "unprocessable entity",
			"statusCode": 422,
			"message":    "Validation failed",
			"errors":     authErrors,
		})
	}

	if existUser := models.SelectUserbyEmail(user.Email); existUser.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Email already exists",
		})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to hash password",
		})
	}

	newUser := models.User{
		Email:    user.Email,
		Password: string(hashPassword),
		Role:     user.Role,
	}

	userId, err := models.CreateUser(&newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to create user",
		})
	}

	if newUser.Role == "seller" {
		newSeller := models.Seller{
			UserID: userId,
			Name:   user.Name,
			Phone:  user.Phone,
		}

		if err := models.CreateSeller(&newSeller); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "server error",
				"statusCode": 500,
				"message":    "Failed to create seller",
			})
		}
	} else if newUser.Role == "customer" {
		newCustomer := models.Customer{
			UserID: userId,
			Name:   user.Name,
			Phone:  user.Phone,
		}

		if err := models.CreateCustomer(&newCustomer); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "server error",
				"statusCode": 500,
				"message":    "Failed to create customer",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "User created successfully",
	})
}

func LoginUser(c *fiber.Ctx) error {
	var login models.User
	if err := c.BodyParser(&login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	user := middlewares.XSSMiddleware(&login).(*models.User)
	if authErrors := helpers.PasswordValidation(user.Password, helpers.StructValidation(user)); len(authErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "unprocessable entity",
			"statusCode": 422,
			"message":    "Validation failed",
			"errors":     authErrors,
		})
	}

	existUser := models.SelectUserbyEmail(user.Email)
	if existUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Email not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(login.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "unauthorized",
			"statusCode": 401,
			"message":    "Invalid password",
		})
	}

	payload := map[string]interface{}{
		"id":    existUser.ID,
		"email": existUser.Email,
		"role":  existUser.Role,
	}

	token, err := helpers.GenerateToken(os.Getenv("SECRETKEY"), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to generate token",
		})
	}

	refreshToken, err := helpers.GenerateRefreshToken(os.Getenv("SECRETKEY"), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Could not generate refresh token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":        "success",
		"statusCode":    200,
		"message":       "Login successfully",
		"email":         existUser.Email,
		"role":          existUser.Role,
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func RequestResetPassword(c *fiber.Ctx) error {
	var requestEmail models.User
	if err := c.BodyParser(&requestEmail); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	user := middlewares.XSSMiddleware(&requestEmail).(*models.User)
	if authErrors := helpers.FieldRequiredValidation(user.Email, "required,email"); authErrors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "unprocessable entity",
			"statusCode": 422,
			"message":    "Validation failed",
			"errors":     authErrors,
		})
	}

	if existUser := models.SelectUserbyEmail(user.Email); existUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Email not found",
		})
	}

	payload := map[string]interface{}{
		"email": user.Email,
	}

	token, err := helpers.GenerateToken(os.Getenv("SECRETKEY"), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":     "accepted",
		"statusCode": 202,
		"message":    "Password reset email sent",
		"token":      token,
	})
}

func ResetPassword(c *fiber.Ctx) error {
	email, err := middlewares.JWTResetPassword(c)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"status":     fiberErr.Message,
				"statusCode": fiberErr.Code,
				"message":    fiberErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	if existUser := models.SelectUserbyEmail(email); existUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Email not found",
		})
	}

	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}
	updatedUser.Email = email

	user := middlewares.XSSMiddleware(&updatedUser).(*models.User)
	if authErrors := helpers.PasswordValidation(user.Password, helpers.StructValidation(user)); len(authErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "unprocessable entity",
			"statusCode": 422,
			"message":    "Validation failed",
			"errors":     authErrors,
		})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to hash password",
		})
	}
	user.Password = string(hashPassword)

	if err := models.UpdateUserbyEmail(user.Email, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to reset password",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 201,
		"message":    "Reset password successfully",
	})
}

func CreateRefreshToken(c *fiber.Ctx) error {
	var refreshData struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&refreshData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to parse request body",
		})
	}

	token, err := helpers.GenerateToken(os.Getenv("SECRETKEY"), map[string]interface{}{"refreshToken": refreshData.RefreshToken})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Could not generate access token",
		})
	}

	refreshToken, err := helpers.GenerateRefreshToken(os.Getenv("SECRETKEY"), map[string]interface{}{"refreshToken": refreshData.RefreshToken})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Could not generate refresh token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":        "success",
		"statusCode":    201,
		"message":       "Refresh successfully",
		"token":         token,
		"refresh_token": refreshToken,
	})
}
