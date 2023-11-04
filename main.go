package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"os"
	"time"
)

type Bike struct {
	ID           string `bson:"_id"`
	Name         string `bson:"name"`
	Type         string `bson:"type"`
	PricePerHour uint16 `bson:"pricePerHour"`
	Owner        string `bson:"owner"`
}

type User struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	Surname     string `bson:"surname"`
	PhoneNumber string `bson:"phoneNumber"`
	Email       string `bson:"email"`
}

type Rent struct {
	ID   string    `bson:"id"`
	Day  time.Time `bson:"day"`
	User `bson:"customer"`
	Bike `bson:"bike"`
}

var (
	mongoPrefix = "mongodb://"
	host        = getStrEnv("HOST", "localhost")
	port        = getStrEnv("HOST_PORT", "27017")
	dbName      = getStrEnv("DB", "bikeService")
	urlSplitter = "/"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoPrefix+host+":"+port+urlSplitter+dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	database := client.Database(dbName)

	bikes := make([]Bike, 0)
	users := make([]User, 0)
	rents := make([]Rent, 0)

	for i := 0; i < 10000; i++ {
		bike := Bike{
			ID:           primitive.NewObjectID().Hex(),
			Name:         []string{"Specialized", "Giant", "Trek", "Scott", "BMC", "Santa Cruz", "Norco", "Cube"}[i%8],
			Type:         []string{"MTB", "Downhill", "Freeride", "Gravel"}[i%4],
			PricePerHour: uint16((i % 20) + 21),
			Owner:        []string{"John", "Mike", "Tom", "Jack", "Thomas", "Jonny", "Martin"}[i%7],
		}
		bikes = append(bikes, bike)
	}

	rand.Seed(time.Now().Unix())
	for i := 0; i < 10000; i++ {
		user := User{
			ID:          primitive.NewObjectID().Hex(),
			Name:        []string{"John", "Mike", "Tom", "Jack", "Thomas", "Jonny", "Martin"}[rand.Intn(7)],
			Surname:     []string{"Smith", "Johnson", "Brown", "Wilson", "Lee", "Davis", "Evans"}[rand.Intn(7)],
			PhoneNumber: []string{"555-1234", "555-5678", "555-9876", "555-4321", "555-8765"}[rand.Intn(5)],
			Email:       fmt.Sprintf("email%d@%s", i, []string{"example.com", "gmail.com", "yahoo.com", "hotmail.com"}[rand.Intn(4)]),
		}
		users = append(users, user)
	}

	for i := 0; i < 10000; i++ {
		rent := Rent{
			ID:   primitive.NewObjectID().Hex(),
			Day:  time.Now(),
			Bike: bikes[i%20],
			User: users[i%20],
		}
		rents = append(rents, rent)
	}

	bikesCollection := database.Collection("bikes")
	usersCollection := database.Collection("users")
	rentsCollection := database.Collection("rents")

	for _, bike := range bikes {
		_, err = bikesCollection.InsertOne(context.Background(), bike)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, user := range users {
		_, err = usersCollection.InsertOne(context.Background(), user)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, rent := range rents {
		_, err = rentsCollection.InsertOne(context.Background(), rent)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Данные успешно вставлены в коллекции")
}

func getStrEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); len(value) == 0 {
		return defaultValue
	} else {
		return value
	}
}
