package routes

import (
	"go-final/database"
	"go-final/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func PostCreateTeam(c *fiber.Ctx) error {
	db := database.DB
	AuthToken := c.Cookies("AuthToken")
	var newTeam model.Team
	c.BodyParser(&newTeam)

	Cookie := db.GetCookie(map[string]any{"cookie.value": AuthToken})
	User := db.GetUser(map[string]any{"id": Cookie.UserID})

	if newTeam.Members == nil {
		newTeam.Members = make([]model.UserPrivilege, 0)
	}
	newTeam.Members = append(newTeam.Members, model.UserPrivilege{ID: User.ID, Privilege: "Owner"})
	newTeam.Tables = make([]string, 0)

	if !db.AddTeam(newTeam) {
		return c.Redirect("/error")
	}
	for _, member := range newTeam.Members {
		update := map[string]any{
			"$addToSet": map[string]any{
				"teams": newTeam.ID,
			},
		}
		if !db.UpdateUser(member.ID, update) {
			return c.Redirect("/error")
		}	
	}

	return c.Redirect("/teams/team/" + newTeam.ID)
}

func GetCreateTeam(c *fiber.Ctx) error {
	AuthToken := c.Cookies("AuthToken")

	if AuthToken == "" {
		return c.Redirect("/error")
		// return c.Render("createTeam", fiber.Map{})
	}

	return c.Render("createTeam", fiber.Map{})
}

func GetTeam(c *fiber.Ctx) error {
	db := database.DB
	AuthToken := c.Cookies("AuthToken")
	Cookie := db.GetCookie(map[string]any{"cookie.value": AuthToken})
	User := db.GetUser(map[string]any{"id": Cookie.UserID})
	Team := db.GetTeam(c.Params("team_id"))

	for _, member := range Team.Members {
		if member.ID == User.ID {
			return myTeam(c, Team, member.Privilege)
		}
	}
	return team(c, Team)
}

func SearchTeams(c *fiber.Ctx) error {
	query := c.Query("search")
	if query == "" {
		return c.Render("searchTeams", fiber.Map{})
	}
	db := database.DB
	Teams := db.GetTeams(query)

	return c.Render("searchTeams", fiber.Map{"Teams": Teams})
}

func DeleteTeam(c *fiber.Ctx) error {
	db := database.DB
	AuthToken := c.Cookies("AuthToken")
	Cookie := db.GetCookie(map[string]any{"cookie.value": AuthToken})
	teamID := c.Params("team_id")
	Team := db.GetTeam(teamID)
	if Team.ID == "" {
		return c.Render("error", nil)
	}
	
	Allowed := false
	for _, member := range Team.Members {
		if member.ID == Cookie.UserID {
			if member.Privilege == "Admin" || member.Privilege == "Owner" {
				Allowed = true
				break;				
			}
		}
	}
	if !Allowed {
		return c.Render("error", nil)
	}

	for _, table := range Team.Tables {	
		for _, member := range Team.Members {
			db.UpdateUser(member.ID, bson.M{"$pull": bson.M{"tables": table} })
		}
		db.DeleteTable(table)
	}
	for _, member := range Team.Members {
		db.UpdateUser(member.ID, bson.M{"$pull": bson.M{"teams": c.Params("team_id")} })
	}

	if (!db.DeleteTeam(c.Params("team_id"))) {
		return c.Render("error", nil)
	}
	return c.Redirect("/")
}

func myTeam(c *fiber.Ctx, Team model.Team, Privilege string) error {
	return c.Render("myteam", fiber.Map{"Team": Team, "Privilege": Privilege})
}

func team(c *fiber.Ctx, Team model.Team) error {
	return c.Render("team", fiber.Map{"Team": Team})
}
