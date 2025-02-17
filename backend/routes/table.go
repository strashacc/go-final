package routes

import (
	"go-final/database"
	"go-final/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func getTable(c *fiber.Ctx) error {
	db := database.DB
	Table := db.GetTable(c.Params("table_id"))
	AuthToken := c.Cookies("AuthToken")
	Cookie := db.GetCookie(map[string]any{"cookie.value": AuthToken})
	
	if Table.Team {
		Team := db.GetTeam(Table.Owner)
		Allowed := 0
		for _, member := range Team.Members {
			if member.ID == Cookie.UserID {
				Allowed = 1
				if member.Privilege == "Admin" || member.Privilege == "Owner" {
					Allowed = 2				
				}
				break;
			}
		}
		if Allowed == 2 {
			return c.Render("teamTable", fiber.Map{"Table": Table, "Editable": true, "Members": Team.Members})
		} else if Allowed == 1 {
			return c.Render("teamTable", fiber.Map{"Table": Table, "Editable": false, "User": Cookie.UserID})
		} else {
			return c.Render("error", nil)
		}
	} else {
		if Table.Owner != Cookie.UserID {
			return c.Render("error", nil)
		}
		return c.Render("personalTable", fiber.Map{"Table": Table})
	}
}

func getCreateTable(c *fiber.Ctx) error {
	// db := database.DB

	// AuthToken := c.Cookies("AuthToken")
	// Cookie := db.GetCookie(map[string]any{"cookie.value": AuthToken})
	// User := db.GetUser(map[string]any{"id": Cookie.UserID})

	return c.Render("createTable", nil)
}

func postCreateTable(c *fiber.Ctx) error {
	db := database.DB
	var Table model.Table
	c.BodyParser(&Table)
	AuthToken := c.Cookies("AuthToken")
	Cookie := db.GetCookie(map[string]any{"cookie.value": AuthToken})
	User := db.GetUser(map[string]any{"id": Cookie.UserID})
	if User.ID == "" {
		return c.Render("error", nil)
	}

	Table.Owner = User.ID
	Table.Team = false
	Table.Items = make([]int, 0)
	for i := range Table.Columns {
		Table.Columns[i].Index = i
	}

	if (!db.AddTable(Table)) {
		return c.Render("error", nil)
	}
	update := map[string]any{
		"$addToSet": map[string]any{
			"tables": Table.ID,
		},
	}
	if !db.UpdateUser(Cookie.UserID, update) {
		return c.Redirect("/error")
	}	

	return c.Redirect("/tables/view/" + Table.ID) //Change
}

func getCreateTeamTable(c *fiber.Ctx) error {
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

	return c.Render("createTable", fiber.Map{"Team": Team})
}

func postCreateTeamTable(c *fiber.Ctx) error {
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

	var Table model.Table
	c.BodyParser(&Table)

	Table.Owner = teamID
	Table.Team = true
	Table.Items = make([]int, 0)
	for i := range Table.Columns {
		Table.Columns[i].Index = i
	}
	for _, member := range Team.Members {
		update := map[string]any{
			"$addToSet": map[string]any{
				"tables": Table.ID,
			},
		}
		if !db.UpdateUser(member.ID, update) {
			return c.Redirect("/error")
		}	
	}

	if (!db.AddTable(Table)) {
		return c.Render("error", nil)
	}
	if (!db.UpdateTeam(Table.Owner, bson.M{"$addToSet": bson.M{"tables": Table.ID} })) {
		return c.Redirect("/");
	}

	return c.Redirect("/tables/view/" + Table.ID) //Change
}

func UpdateTable(c *fiber.Ctx) error {
	db := database.DB
	var Table model.Table
	c.BodyParser(&Table)
	tableID := c.Params("table_id")
	OrigTable := db.GetTable(tableID)
	AuthToken := c.Cookies("AuthToken")
	Cookie := db.GetCookie(map[string]any{"cookie.value": AuthToken})

	if OrigTable.Team {
		Team := db.GetTeam(OrigTable.Owner)
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
	} else {
		if Cookie.UserID != OrigTable.Owner {
			return c.Render("error", nil)
		}
	}
	Table.Owner = OrigTable.Owner
	if len(Table.Columns) != 0{
		for i := range Table.Columns[0].Items {
			Table.Items = append(Table.Items, i)
		}
	}
	Table.Name = OrigTable.Name
	Table.ID = OrigTable.ID
	Table.Team = OrigTable.Team
	for i := range Table.Columns {
		Table.Columns[i].Index = i
		Table.Columns[i].Name = OrigTable.Columns[i].Name
		Table.Columns[i].Type = OrigTable.Columns[i].Type
		Table.Columns[i].Options = OrigTable.Columns[i].Options
	}

	if !db.UpdateTable(tableID, Table) {
		return c.Render("error", nil)
	}

	return c.Redirect("/tables/view/" + c.Params("table_id"))
}

func DeleteTable(c *fiber.Ctx) error {
	db := database.DB
	tableID := c.Params("table_id")
	OrigTable := db.GetTable(tableID)
	AuthToken := c.Cookies("AuthToken")
	Cookie := db.GetCookie(map[string]any{"cookie.value": AuthToken})

	if OrigTable.Team {
		Team := db.GetTeam(OrigTable.Owner)
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
	} else {
		if Cookie.UserID != OrigTable.Owner {
			return c.Render("error", nil)
		}
	}
	if OrigTable.Team {
		Team := db.GetTeam(OrigTable.Owner)

		for _, member := range Team.Members {
			db.UpdateUser(member.ID, bson.M{"$pull": bson.M{"tables": tableID} })
		}

	} else {
		db.UpdateUser(OrigTable.Owner, bson.M{"$pull": bson.M{"tables": tableID} })
	}

	if !db.DeleteTable(tableID) {
		c.Render("error", nil)
	}
	return c.Redirect("/")
}