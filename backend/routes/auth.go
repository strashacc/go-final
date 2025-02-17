package routes

import (
	"go-final/database"
	"go-final/model"
	"go-final/scripts"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetLogin(c *fiber.Ctx) error {
	if Cookie := c.Cookies("AuthToken"); Cookie != "" {
		db := database.DB
		getCookie := db.GetCookie(map[string]any{"cookie.value": Cookie})
		User := db.GetUser(map[string]any{"id": getCookie.UserID})
		if User.ID == "" {
			c.Cookie(&fiber.Cookie{Name: "AuthToken", Path: "/", Value: "", Expires: time.Time{}})
			return c.Render("login", fiber.Map{})
		}
		return c.Redirect("/profile/" + getCookie.UserID)
	}
	return c.Render("login", fiber.Map{})
}
func PostLogin(c *fiber.Ctx) error {
	if Cookie := c.Cookies("AuthToken"); Cookie != "" {
		return c.Redirect("/")
	}

	db := database.DB

	var credentials struct {
		ID       string
		Password string
	}
	c.BodyParser(&credentials)

	query := make(map[string]any)
	query["id"] = credentials.ID
	dbResponse := db.GetUser(query)
	if dbResponse.Password != credentials.Password {
		return c.Render("error", fiber.Map{"Error": struct{ Title string }{Title: "Some error with logging in"}})
	}

	AuthToken := scripts.GenerateSecureToken(16)
	Cookie := fiber.Cookie{Name: "AuthToken", Value: AuthToken}
	toDB := struct {
		UserID string
		Cookie fiber.Cookie
	}{UserID: dbResponse.ID, Cookie: Cookie}
	if !db.AddCookie(toDB) {
		return c.SendString("Error")
	}

	c.Cookie(&Cookie)
	return c.Redirect("/")
}

func GetSignUp(c *fiber.Ctx) error {
	if Cookie := c.Cookies("AuthToken"); Cookie != "" {
		return c.Redirect("/")
	}
	return c.Render("signup", fiber.Map{})
}
func PostSignUp(c *fiber.Ctx) error {
	if Cookie := c.Cookies("AuthToken"); Cookie != "" {
		return c.Redirect("/")
	}

	db := database.DB
	var user model.User

	c.BodyParser(&user)

	user.Teams = make([]string, 0)
	user.Tables = make([]string, 0)

	if !db.AddUser(user) {
		return c.Render("error", fiber.Map{"Error": struct{ Title string }{Title: "Could not create user"}})
	}

	AuthToken := scripts.GenerateSecureToken(16)
	Cookie := fiber.Cookie{Name: "AuthToken", Value: AuthToken}
	toDB := struct {
		UserID string
		Cookie fiber.Cookie
	}{UserID: user.ID, Cookie: Cookie}
	if !db.AddCookie(toDB) {
		return c.SendString("Error")
	}

	c.Cookie(&Cookie)
	return c.Redirect("/")
}
