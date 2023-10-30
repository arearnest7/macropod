package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"net"

	log "github.com/sirupsen/logrus"

	"time"

	"crypto/sha256"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	pb "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/user"
	tracing "github.com/vhive-serverless/vSwarm/utils/tracing/go"
)

type User struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

// loadUsers loads hotel users from database
func loadUsers(client *mongo.Client) map[string]string {

	coll := client.Database("user-db").Collection("user")

	filter := bson.D{}
	// filter := bson.M{{"username": username}}
	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		log.Println("Failed get users data: ", err)
	}

	// Get a list of all returned documents and print them out.
	// See the mongo.Cursor documentation for more examples of using cursors.
	var users []User
	if err = cursor.All(context.Background(), &users); err != nil {
		log.Println("Failed get users data: ", err)
	}

	res := make(map[string]string)
	for _, user := range users {
		res[user.Username] = user.Password
	}

	fmt.Printf("Done load users\n")

	return res
}

func lookupCache(username string) string {
	users_cached = loadUsers(MongoClient)
	res, ok := users_cached[username]
	if !ok {
		log.Println("User does not exist: ", username)
	}
	return res
}

//
func lookUpDB(username string) (User, bool) {
	// session := s.MongoClient.Copy()
	// defer session.Close()
	collection := MongoClient.Database("user-db").Collection("user")

	// listAll(collection)

	// unmarshal json profiles
	var user User
	filter := bson.D{primitive.E{Key: "username", Value: username}}
	// filter := bson.M{{"username": username}}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Println("Failed get user: ", err)
		return user, false
	}
	return user, true
}

// CheckUser returns whether the username and password are correct.
func CheckUser(var req) (*pb.Result, error) {
	res := new(pb.Result)

	fmt.Printf("CheckUser: %+v", req)

	sum := sha256.Sum256([]byte(req.Password))
	pass := fmt.Sprintf("%x", sum)

	use_cache := false

	if use_cache {
		password := lookupCache(req.Username)
		res.Correct = pass == password
	} else {
		user, _ := lookUpDB(req.Username)
		res.Correct = pass == user.Password
	}

	fmt.Printf(" >> pass: %t\n", res.Correct)

	return res, nil
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
        fmt.Fprintf(res, CheckUser(req.json)) // echo to caller
}
