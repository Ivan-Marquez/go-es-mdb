package main

import (
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	col := "user"
	pipeline := mongo.Pipeline{
		bson.D{primitive.E{
			Key: "$match", Value: bson.D{primitive.E{
				Key: "operationType", Value: bson.D{primitive.E{
					Key: "$in", Value: []string{"insert", "update"},
				}},
			}},
		}},
	}

	changeStreamHandler(col, pipeline, &User{})
}
