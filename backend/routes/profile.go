package routes

import (
	"go-final/database"
	"go-final/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func GetProfile(c *fiber.Ctx) error {
	db := database.DB

	AuthToken := c.Cookies("AuthToken")

	Profile := db.GetUser(map[string]any{"id": c.Params("username")})

	if AuthToken != "" {
		Cookie := db.GetCookie(map[string]any{"cookie.value": AuthToken})
		User := db.GetUser(map[string]any{"id": Cookie.UserID})

		if User.ID == c.Params("username") {
			return myProfile(c, Profile)
		} else {
			return profile(c, Profile)
		}
	}
	return profile(c, Profile)
}
func DeleteProfile(c *fiber.Ctx) error {
	db := database.DB

	AuthToken := c.Cookies("AuthToken")
	if AuthToken == "" {
		return c.Redirect("/")
	}
	Cookie := db.GetCookie(map[string]any{"cookie.value": AuthToken})

	User := db.GetUser(map[string]any{"id": Cookie.UserID})

	if !db.DeleteUser(User.ID) {
		log.Error("Could not delete user from the database")
		return c.Redirect("/")
	}
	if !db.DeleteCookies(map[string]any{"userid": Cookie.UserID}) {
		log.Error("Could not delete user cookies from the database")
	}

	c.Cookie(&fiber.Cookie{Name: "AuthToken", Path: "/", Value: "", Expires: time.Time{}})
	return c.Redirect("/")
}
func Logout(c *fiber.Ctx) error {
	db := database.DB

	AuthToken := c.Cookies("AuthToken")
	db.DeleteCookies(map[string]any{"cookie.value": AuthToken})

	c.Cookie(&fiber.Cookie{Name: "AuthToken", Path: "/", Value: "", Expires: time.Time{}})
	return c.Redirect("/")
}

func myProfile(c *fiber.Ctx, Profile model.User) error {
	return c.Render("myprofile", fiber.Map{"User": Profile})
}

func profile(c *fiber.Ctx, Profile model.User) error {
	return c.Render("profile", fiber.Map{"User": Profile})
}
