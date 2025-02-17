package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-final/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *Database

type Database struct {
	connection mongo.Client
}

func New(URI string) {

	uri := os.Getenv("MONGODB_URI")
	docs := "www.mongodb.com/docs/drivers/go/current/"
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " + docs +
			"usage-examples/#environment-variable")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	result := Database{connection: *client}

	DB = &result
}

func (db *Database) Disconnect() {
	if err := db.connection.Disconnect(context.TODO()); err != nil {
		panic(err)
	} else {
		fmt.Println("Disconnected Successfully")
	}
}

func (db *Database) Ping() {
	// Send a ping to confirm a successful connection
	ping := db.connection.Database(os.Getenv("DATABASE")).RunCommand(context.TODO(), bson.M{"ping": 1}).Err()
	if err := ping; err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func (db *Database) GetUsers(query map[string]any) []model.User {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("users")

	cur, err := coll.Find(context.TODO(), query)
	if err != nil {
		panic(err)
	} else {
		var results []model.User
		if err := cur.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		return results
	}
}

func (db *Database) GetUser(query map[string]any) model.User {
	var result model.User

	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("users")
	if err := coll.FindOne(context.TODO(), query).Decode(&result); err != nil {
		log.Println(err)
		return model.User{}
	} else {
		return result
	}
}

func (db *Database) AddUser(user model.User) bool {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("users")
	if _, err := coll.InsertOne(context.TODO(), user); err != nil {
		log.Println("Failed to insert a user into database")
		return false
	}

	return true
}

func (db *Database) UpdateUser(username string, update any) bool {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("users")
	if _, err := coll.UpdateOne(context.TODO(), bson.M{"id": username}, update); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (db *Database) DeleteUser(username string) bool {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("users")
	if _, err := coll.DeleteOne(context.TODO(), bson.M{"id": username}); err != nil {
		log.Panicln(err)
		return false
	}
	return true
}

func (db *Database) GetTeam(id string) model.Team {
	var result model.Team

	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("teams")
	if err := coll.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result); err != nil {
		log.Println(err)
	}
	return result
}

func (db *Database) GetTeams(search string) []model.Team {
	var results []model.Team

	match := bson.M{"$text": bson.M{"$search": search} }
	sort := bson.M{"score": bson.M{"$meta": "textScore"} }
	var req []bson.M
	req = append(req, bson.M{"$match": match})
	req = append(req, bson.M{"$sort": sort})

	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("teams")
	if cursor, err := coll.Aggregate(context.TODO(), req ); err != nil {
		log.Println(err)
	} else {

		if err := cursor.All(context.TODO(), &results); err != nil {
			log.Println(err)
		}
	}
	// if cursor, err := coll.Find(context.TODO(), bson.M{"$text": bson.M{"$search": search} }); err != nil {
	// 	log.Println(err)
	// } else {

	// 	if err := cursor.All(context.TODO(), &results); err != nil {
	// 		log.Println(err)
	// 	}
	// }
	return results
}

func (db *Database) AddTeam(team model.Team) bool {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("teams")
	if _, err := coll.InsertOne(context.TODO(), team); err != nil {
		return false
	}

	return true
}

func (db *Database) UpdateTeam(id string, update any) bool {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("teams")
	if _, err := coll.UpdateOne(context.TODO(), bson.M{"id": id}, update); err != nil {
		return false
	}

	return true
}

func (db *Database) DeleteTeam(id string) bool {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("teams")
	if _, err := coll.DeleteOne(context.TODO(), bson.M{"id": id}); err != nil {
		return false
	}

	return true
}

func (db *Database) AddCookie(Cookie interface{}) bool {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("cookies")

	if _, err := coll.InsertOne(context.TODO(), Cookie); err != nil {
		return false
	}

	return true
}

func (db *Database) GetCookie(query map[string]any) struct {
	UserID string
	Cookie fiber.Cookie
} {
	var Cookie struct {
		UserID string
		Cookie fiber.Cookie
	}

	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("cookies")
	if err := coll.FindOne(context.TODO(), query).Decode(&Cookie); err != nil {
		log.Println(err)
		return Cookie
	}

	return Cookie
}

func (db *Database) DeleteCookies(query map[string]any) bool {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("cookies")
	if _, err := coll.DeleteMany(context.TODO(), query); err != nil {
		return false
	}
	return true
}

func (db *Database) AddTable(Table model.Table) bool {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("tables")
	if _, err := coll.InsertOne(context.TODO(), Table); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (db *Database) GetTable(id string) model.Table {	
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("tables")
	var table model.Table
	if err := coll.FindOne(context.TODO(), bson.M{"id": id}).Decode(&table); err != nil {
		log.Println(err)
	}
	return table
}

func (db *Database) UpdateTable(id string, Table model.Table) bool {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("tables")
	if _, err := coll.ReplaceOne(context.TODO(), bson.M{"id": id}, Table); err != nil {
		log.Println(err)
		return false
	}
	return true
}
func (db *Database) DeleteTable(id string) bool {
	coll := db.connection.Database(os.Getenv("DATABASE")).Collection("tables")
	if _, err := coll.DeleteOne(context.TODO(), bson.M{"id": id}); err != nil {
		log.Println(err)
		return false
	}
	return true
}