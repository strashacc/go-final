package model


type User struct {
	ID string `bson: "_id", json: ID`
	Name string `json: Name`
	Email string
	Password string
	Teams []string
	Tables []string
}