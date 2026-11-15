package function

import (
    "fmt"
    "strconv"

    log "github.com/sirupsen/logrus"

    "crypto/sha256"
)

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
    //    log.Println("Failed get users data: ", err)
    //}

    // Get a list of all returned documents and print them out.
    // See the mongo.Cursor documentation for more examples of using cursors.
    var users []User
        //if err = cursor.All(context.Background(), &users); err != nil {
    //    log.Println("Failed get users data: ", err)
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
    //    log.Println("Failed get user: ", err)
    //    return user, false
    //}
    return user, true
}

// CheckUser returns whether the username and password are correct.
func CheckUser(req map[string]interface{}) bool {
    var res bool

    fmt.Printf("CheckUser: %+v", req)

    sum := sha256.Sum256([]byte(req["Password"].(string)))
    pass := fmt.Sprintf("%x", sum)

    use_cache := false

    if use_cache {
        password := lookupCache(req["Username"].(string))
        res = pass == password
    } else {
        user, _ := lookUpDB(req["Username"].(string))
        res = pass == user.Password
    }

    fmt.Printf(" >> pass: %t\n", res)

    return res
}

func FunctionHandler(context Context) (string, int) {
    return strconv.FormatBool(CheckUser(context.JSON)), 200
}
