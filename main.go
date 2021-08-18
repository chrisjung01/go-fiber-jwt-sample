package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
	"log"
	"strconv"
	"time"
)

func authorizeRequired() func(ctx *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:    []byte(secret),
		SigningMethod: "HS512",
	})
}

var secret string

func main() {
	buildNewSecret()

	// create a new fiber app
	app := fiber.New()

	app.Get("/user", authorizeRequired(), returnUser)

	app.Get("/login", login)

	app.Get("/hello", authorizeRequired(), hello)

	err := app.Listen(":3000")
	if err != nil {
		log.Println(err.Error())
	}
}

/*
rest token to test
curl --header "Content-Type: application/json" --header "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNjI5NTc1Mzg1LCJuYW1lIjoiSm9obiBEb2UifQ.Vf0zI0YQADwvFUYrFxaRQLEgRdL0qXW_aRRafWPH5ZZR4fr4EECRHhsMSV4Gv27GbjYEwfSuAIhnlrK2AitAPw" --request GET --data '{"email": "email","password": "password"}' http://localhost:3000/hello
 */
func hello(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "okay",
	})
}

func buildNewSecret() {
	var hmacSampleSecret = make([]byte, 256)

	// create random key
	_, err := rand.Read(hmacSampleSecret)
	if err != nil {
		log.Println(err.Error())
	}

	secret = base64.StdEncoding.EncodeToString(hmacSampleSecret)
	println("secret: " + secret)
}

/*
rest token to test
curl --header "Content-Type: application/json" --header "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNjI5NTc1Mzg1LCJuYW1lIjoiSm9obiBEb2UifQ.Vf0zI0YQADwvFUYrFxaRQLEgRdL0qXW_aRRafWPH5ZZR4fr4EECRHhsMSV4Gv27GbjYEwfSuAIhnlrK2AitAPw" --request GET --data '{"email": "email","password": "password"}' http://localhost:3000/user
 */
func returnUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	isAdmin := claims["admin"].(bool)
	return c.SendString("Welcome " + name + " admin: " + strconv.FormatBool(isAdmin))
}

/*
curl --header "Content-Type: application/json" --request GET --data '{"email": "email","password": "password"}' http://localhost:3000/login
 */
func login(ctx *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body request
	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot pars json",
		})
	}

	if body.Email != "email" || body.Password != "password" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Bad Credentials",
		})
	}

	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "John Doe"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(secret))

	println("tokenstring: ", tokenString)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": tokenString,
	})
}
