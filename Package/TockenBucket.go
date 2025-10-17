package Package

// import (
// 	"context"
// 	"fmt"
// 	"math"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/redis/go-redis/v9"
// )

// var IntailToken int
// var BucketSize int

// func TokenBuucketAlgorithm(RedisClient *redis.Client, Bucketsize int64, refill float64) gin.HandlerFunc {
// 	ctx := context.Background()
// 	return func(c *gin.Context) {
// 		//Get the Cilent Ip
// 		ipAddress := c.ClientIP()
// 		hashkey := fmt.Sprintf("Rate limit %s", ipAddress)
// 		// fmt.Println("HashKey", hashkey)

// 		var tokens float64
// 		var lastTimeStap int64
// 		//Now make the transcationn
// 		pipe := RedisClient.TxPipeline()
// 		bucketstate := pipe.HGetAll(ctx, hashkey)
// 		pipe.Exec(ctx) //for existing the command
// 		bucketData, err := bucketstate.Result()
// 		// fmt.Println("Bucket Data ->", bucketData)
// 		if err != nil || len(bucketData) == 0 {
// 			//This mean it is new user
// 			initialToken := Bucketsize - 1
// 			pipe = RedisClient.TxPipeline()
// 			pipe.HSet(ctx, hashkey, "token", initialToken)
// 			pipe.HSet(ctx, hashkey, "timeStap", time.Now().Unix())
// 			//Persist the token //Clear the token after 1 hour
// 			pipe.Expire(ctx, hashkey, 1*time.Hour)
// 			_, err := pipe.Exec(ctx)
// 			if err != nil {
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error in setting the value to the new user"})
// 				return

// 			}
// 			c.Next()
// 			return
// 		}
// 		//Check for existing user
// 		//
// 		// fmt.Println("Existing user part")
// 		RemainingTokenStr := bucketData["token"]
// 		lastTimeStapStr := bucketData["timeStap"]
// 		RemainingToken, _ := strconv.ParseFloat(RemainingTokenStr, 64)
// 		lastTimeStap, _ = strconv.ParseInt(lastTimeStapStr, 10, 64)

// 		//Calcualte the time to add the token
// 		Now := time.Now().Unix()
// 		elapsedSecond := Now - lastTimeStap
// 		//Refill rate
// 		tokenToAdd := float64(elapsedSecond) * refill
// 		// fmt.Println("Token to Add", tokenToAdd)
// 		tokens = math.Min(float64(Bucketsize), RemainingToken+tokenToAdd)
// 		// fmt.Println(tokens)

// 		if tokens < 1 {
// 			//Exceed the request abort the chain
// 			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests are made ... Try after some time"})
// 			return
// 		}
// 		RemainingToken = tokens - 1
// 		pipe = pipe.TxPipeline()
// 		pipe.HSet(ctx, hashkey, "token", RemainingToken)
// 		pipe.HSet(ctx, hashkey, "timeStap", Now)
// 		_, err = pipe.Exec(ctx)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 				"message": "Error in setting the value",
// 			})
// 			return

// 		}
// 		fmt.Printf("Request from %s and remaining token %0.2f", ipAddress, RemainingToken)
// 		c.Next()

// 	}
// }
