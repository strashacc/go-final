package routes

import (
	"go-final/database"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
	"github.com/gofiber/fiber/v2/middleware/redirect"
)

func Setup() {
	engine := django.New("../frontend/views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/": "/auth/login",
			"/profile": "/auth/login",
		},
		StatusCode: 301,
	}))
	app.Static("/static", "../frontend/static")

	app.Get("/", func(c *fiber.Ctx) error {
		db := database.DB

		if AuthToken := c.Cookies("AuthToken"); AuthToken != "" {
			Cookie := db.GetCookie(map[string]any{"cookie.value": AuthToken})
			log.Println(AuthToken)
			return c.Redirect("/profile/" + Cookie.UserID)
		} else {
			return c.Redirect("/auth/login")
		}
	})
	app.Get("/error", func(c *fiber.Ctx) error {
		return c.Render("error", nil)
	})

	auth := app.Group("/auth")

	auth.Get("/login", GetLogin)
	auth.Post("/login", PostLogin)
	auth.Get("/signup", GetSignUp)
	auth.Post("/signup", PostSignUp)

	profile := app.Group("/profile")

	profile.Get("/:username", GetProfile)
	profile.Post("/logout", Logout)
	profile.Post("/delete", DeleteProfile)

	teams := app.Group("/teams")
	teams.Post("/create", PostCreateTeam)
	teams.Get("/create", GetCreateTeam)
	teams.Get("/team/:team_id", GetTeam)
	teams.Get("/search", SearchTeams)
	teams.Post("/delete/:team_id", DeleteTeam)

	tables := app.Group("/tables")
	tables.Get("/create", getCreateTable)
	tables.Post("/create", postCreateTable)
	tables.Get("/:team_id/create", getCreateTeamTable)
	tables.Post("/:team_id/create", postCreateTeamTable)
	tables.Get("/view/:table_id", getTable)
	tables.Post("/update/:table_id", UpdateTable)
	tables.Post("/delete/:table_id", DeleteTable)

	Test(app)

	log.Println(app.Listen(":" + os.Getenv("PORT")))
}

func Test(app *fiber.App) {
	db := database.DB
	db.Ping()
}
