package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rate-limiter/Middleware"
	models "rate-limiter/Models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	var wg sync.WaitGroup
	errchan := make(chan error, 2)
	var RedisClient *redis.Client
	var mongoClient *mongo.Client
	wg.Add(2)
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//redis connection
		defer cancel()
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		pingRedis := RedisClient.Ping(ctx).Err()
		if pingRedis != nil {
			errchan <- fmt.Errorf("error in connecting the redis %w", pingRedis)
			return
		}
		fmt.Println("Redis connected")
	}()
	//MongoDB connection
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		mongo_url := "mongodb://localhost:27017"
		clientOptions := options.Client().ApplyURI(mongo_url)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			errchan <- fmt.Errorf("failed to create mongoClient %w", err)
		}
		mongoClient = client
		pingMongo := mongoClient.Ping(ctx, nil)
		if pingMongo != nil {
			errchan <- fmt.Errorf("error in connecting datatbase %w ", pingMongo)
			return
		}
		fmt.Println("MongoDb Connected")
	}()
	wg.Wait()
	close(errchan)
	for err := range errchan {
		if err != nil {
			log.Fatal(err)
		}
	}
	//Create Server
	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server Healthy"})
	})
	app.POST("/signup", Middleware.CheckUserInput, func(c *gin.Context) {
		val, exists := c.Get("validatedUser")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "error in recieve the input from the middleware"})
			return

		}
		newUser := val.(models.User)
		// if err := c.ShouldBindJSON(&newUser); err != nil {
		// 	fmt.Println("Error ->", err)
		// 	c.JSON(http.StatusBadRequest, gin.H{

		// 		"message": "error in reading the body of the request %w",
		// 	})
		// 	return
		// }
		//Hash the password
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "Failed to hash the password"})
			return
		}
		newUser.Password = string(hashPassword)
		//Create mongo collection
		collection := mongoClient.Database("rate-limiter").Collection("Users")
		_, err = collection.InsertOne(context.TODO(), newUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "error in putting new user in the database",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"messge": "SignUp Successfully",
		})
	})

	app.Run()
}
