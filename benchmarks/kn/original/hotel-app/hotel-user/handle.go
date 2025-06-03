package function

import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"

	log "github.com/sirupsen/logrus"

	"crypto/sha256"
)

type RequestBody struct {
        Request string `json:"Request"`
        RequestType string `json:"RequestType"`
        Lat float64 `json:"Lat"`
        Lon float64 `json:"Lon"`
        HotelId string `json:"HotelId"`
        HotelIds []string `json:"HotelIds"`
        RoomNumber int `json:"RoomNumber"`
        CustomerName string `json:"CustomerName"`
        Username string `json:"Username"`
        Password string `json:"Password"`
        Require string `json:"Require"`
        InDate string `json:"InDate"`
        OutDate string `json:"OutDate"`
}

type User struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

// loadUsers loads hotel users from database
func loadUsers() map[string]string {

	//coll := "{}"

	//filter := bson.D{}
	// filter := bson.M{{"username": username}}
        //cursor, err := coll.Find(context.Background(), filter)
	//if err != nil {
	//	log.Println("Failed get users data: ", err)
	//}

	// Get a list of all returned documents and print them out.
	// See the mongo.Cursor documentation for more examples of using cursors.
	var users []User
        //if err = cursor.All(context.Background(), &users); err != nil {
	//	log.Println("Failed get users data: ", err)
	//}

	res := make(map[string]string)
	for _, user := range users {
		res[user.Username] = user.Password
	}

	fmt.Printf("Done load users\n")

	return res
}

func lookupCache(username string) string {
        //ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	//defer cancel()
	//MongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("HOTEL_APP_DATABASE")))
	//if err != nil { return "" }
	users_cached := loadUsers()
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
        //ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
        //defer cancel()
        //MongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("HOTEL_APP_DATABASE")))
	//collection := "{}"

	// listAll(collection)

	// unmarshal json profiles
	var user User
	//filter := bson.D{primitive.E{Key: "username", Value: username}}
	// filter := bson.M{{"username": username}}
        //err = collection.FindOne(context.Background(), filter).Decode(&user)
	//if err != nil {
	//	log.Println("Failed get user: ", err)
	//	return user, false
	//}
	return user, true
}

// CheckUser returns whether the username and password are correct.
func CheckUser(req RequestBody) bool {
	var res bool

	fmt.Printf("CheckUser: %+v", req)

	sum := sha256.Sum256([]byte(req.Password))
	pass := fmt.Sprintf("%x", sum)

	use_cache := false

	if use_cache {
		password := lookupCache(req.Username)
		res = pass == password
	} else {
		user, _ := lookUpDB(req.Username)
		res = pass == user.Password
	}

	fmt.Printf(" >> pass: %t\n", res)

        return res
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
        body, _ := ioutil.ReadAll(req.Body)
        body_u := RequestBody{}
        json.Unmarshal(body, &body_u)
        defer req.Body.Close()
	ret := strconv.FormatBool(CheckUser(body_u))
        fmt.Fprintf(res, ret) // echo to caller
}
