package model

type Team struct {
	ID          string `bson: "id"`
	Name        string `bson: "name"`
	Visibility  bool `bson: "visibility"`
	Description string `bson: "description"`
	Members     []UserPrivilege `bson: "members"` //Members[username] = privelege
	Tables      []string
}

type UserPrivilege struct {
	ID string `bson: "id"`
	Privilege string `bson: "privilege"`
}