package main

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type HotelDB struct {
	HId    string  `bson:"hotelid"`
	HLat   float64 `bson:"lat"`
	HLon   float64 `bson:"lon"`
	HRate  float64 `bson:"rate"`
	HPrice float64 `bson:"price"`
}

type Point struct {
        Pid  string  `bson:"hotelid"`
        Plat float64 `bson:"lat"`
        Plon float64 `bson:"lon"`
}

type RoomType struct {
        bookableRate float64
        totalRate float64
        totalRateInclusive float64
        code string
        currency string
        roomDescription string
}

type RatePlan struct {
        hotelId string
        code string
        inDate string
        outDate string
        roomType RoomType
}

type Image struct {
        url string
        Default bool
}

type Address struct {
        streetNumber string
        streetName string
        city string
        state string
        country string
        postalCode string
        lat float64
        lon float64
}

type Reservation struct {
        HotelId      string `bson:"hotelid"`
        CustomerName string `bson:"customername"`
        InDate       string `bson:"indate"`
        OutDate      string `bson:"outdate"`
        Number       int    `bson:"number"`
}

type Number struct {
        HotelId string `bson:"hotelid"`
        Number  int    `bson:"numberofroom"`
}

type Hotel struct {
        id string
        name string
        phoneNumber string
        description string
        address Address
        images []Image
}

type User struct {
        Username string `bson:"username" json:"username"`
        Password string `bson:"password" json:"password"`
}

func initializeDatabase()  {
	url := "mongodb://192.168.10.18:27017"

	// GEO
	fmt.Printf("geo db ip addr = %s\n", url)
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	c := session.DB("geo-db").C("geo")
	// First we clear the collection to have always a new one
	if err = c.DropCollection(); err != nil {
		log.Print("DropCollection: ", err)
	}

	count, err := c.Find(&bson.M{"hotelid": "1"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Point{"1", 37.7867, -122.4112})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "2"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Point{"2", 37.7854, -122.4005})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "3"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Point{"3", 37.7854, -122.4071})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "4"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Point{"4", 37.7936, -122.3930})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "5"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Point{"5", 37.7831, -122.4181})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "6"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Point{"6", 37.7863, -122.4015})
		if err != nil {
			log.Fatal(err)
		}
	}

	// add up to 80 hotels
	for i := 7; i <= 80; i++ {
		hotel_id := strconv.Itoa(i)
		count, err = c.Find(&bson.M{"hotelid": hotel_id}).Count()
		if err != nil {
			log.Fatal(err)
		}
		lat := 37.7835 + float64(i)/500.0*3
		lon := -122.41 + float64(i)/500.0*4
		if count == 0 {
			err = c.Insert(&Point{hotel_id, lat, lon})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = c.EnsureIndexKey("hotelid")
	if err != nil {
		log.Fatal(err)
	}

	//PROFILE
	fmt.Printf("profile db ip addr = %s\n", url)
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	c = session.DB("profile-db").C("hotels")
	// First we clear the collection to have always a new one
	if err = c.DropCollection(); err != nil {
		log.Print("DropCollection: ", err)
	}

	count, err = c.Find(&bson.M{"id": "1"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Hotel{
			id:          "1",
			name:        "Clift Hotel",
			phoneNumber: "(415) 775-4700",
			description: "A 6-minute walk from Union Square and 4 minutes from a Muni Metro station, this luxury hotel designed by Philippe Starck features an artsy furniture collection in the lobby, including work by Salvador Dali.",
			address: Address{
				streetNumber: "495",
				streetName:   "Geary St",
				city:         "San Francisco",
				state:        "CA",
				country:      "United States",
				postalCode:   "94102",
				lat:          37.7867,
				lon:          -122.4112}})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"id": "2"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Hotel{
			id:          "2",
			name:        "W San Francisco",
			phoneNumber: "(415) 777-5300",
			description: "Less than a block from the Yerba Buena Center for the Arts, this trendy hotel is a 12-minute walk from Union Square.",
			address: Address{
				streetNumber: "181",
				streetName:   "3rd St",
				city:         "San Francisco",
				state:        "CA",
				country:      "United States",
				postalCode:   "94103",
				lat:          37.7854,
				lon:          -122.4005}})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"id": "3"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Hotel{
			id:          "3",
			name:        "Hotel Zetta",
			phoneNumber: "(415) 543-8555",
			description: "A 3-minute walk from the Powell Street cable-car turnaround and BART rail station, this hip hotel 9 minutes from Union Square combines high-tech lodging with artsy touches.",
			address: Address{
				streetNumber: "55",
				streetName:   "5th St",
				city:         "San Francisco",
				state:        "CA",
				country:      "United States",
				postalCode:   "94103",
				lat:          37.7834,
				lon:          -122.4071}})

		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"id": "4"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Hotel{
			id:          "4",
			name:        "Hotel Vitale",
			phoneNumber: "(415) 278-3700",
			description: "This waterfront hotel with Bay Bridge views is 3 blocks from the Financial District and a 4-minute walk from the Ferry Building.",
			address: Address{
				streetNumber: "8",
				streetName:   "Mission St",
				city:         "San Francisco",
				state:        "CA",
				country:      "United States",
				postalCode:   "94105",
				lat:          37.7936,
				lon:          -122.3930}})

		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"id": "5"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Hotel{
			id:          "5",
			name:        "Phoenix Hotel",
			phoneNumber: "(415) 776-1380",
			description: "Located in the Tenderloin neighborhood, a 10-minute walk from a BART rail station, this retro motor lodge has hosted many rock musicians and other celebrities since the 1950s. It’s a 4-minute walk from the historic Great American Music Hall nightclub.",
			address: Address{
				streetNumber: "601",
				streetName:   "Eddy St",
				city:         "San Francisco",
				state:        "CA",
				country:      "United States",
				postalCode:   "94109",
				lat:          37.7831,
				lon:          -122.4181}})

		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"id": "6"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Hotel{
			id:          "6",
			name:        "St. Regis San Francisco",
			phoneNumber: "(415) 284-4000",
			description: "St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
			address: Address{
				streetNumber: "125",
				streetName:   "3rd St",
				city:         "San Francisco",
				state:        "CA",
				country:      "United States",
				postalCode:   "94109",
				lat:          37.7863,
				lon:          -122.4015}})

		if err != nil {
			log.Fatal(err)
		}
	}

	// add up to 80 hotels
	for i := 7; i <= 80; i++ {
		hotel_id := strconv.Itoa(i)
		count, err = c.Find(&bson.M{"id": hotel_id}).Count()
		if err != nil {
			log.Fatal(err)
		}
		phone_num := "(415) 284-40" + hotel_id
		lat := 37.7835 + float64(i)/500.0*3
		lon := -122.41 + float64(i)/500.0*4
		if count == 0 {
			err = c.Insert(&Hotel{
				id:          hotel_id,
				name:        "St. Regis San Francisco",
				phoneNumber: phone_num,
				description: "St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
				address: Address{
					streetNumber: "125",
					streetName:   "3rd St",
					city:         "San Francisco",
					state:        "CA",
					country:      "United States",
					postalCode:   "94109",
					lat:          lat,
					lon:          lon}})

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = c.EnsureIndexKey("id")
	if err != nil {
		log.Fatal(err)
	}

	//RATE
	fmt.Printf("rate db ip addr = %s\n", url)
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	c = session.DB("rate-db").C("inventory")
	// First we clear the collection to have always a new one
	if err = c.DropCollection(); err != nil {
		log.Print("DropCollection: ", err)
	}
	count, err = c.Count()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" Collection contains: %d elements\n", count)
	count, err = c.Find(&bson.M{"hotelid": "2"}).Count()
	if err != nil {
		log.Fatal(err)
	}

	item := RatePlan{
		hotelId: "1",
		code:    "RACK",
		inDate:  "2015-04-09",
		outDate: "2015-04-10",
		roomType: RoomType{
			bookableRate:       109.00,
			code:               "KNG",
			roomDescription:    "King sized bed",
			totalRate:          109.00,
			totalRateInclusive: 123.17,
		},
	}

	if count == 0 {
		err = c.Insert(&item)
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "2"}).Count()
	if err != nil {
		log.Fatal(err)
	}

	item.hotelId = "2"
	item.code = "RACK"
	item.inDate = "2015-04-09"
	item.outDate = "2015-04-10"
	item.roomType.bookableRate = 139.00
	item.roomType.code = "QN"
	item.roomType.roomDescription = "Queen sized bed"
	item.roomType.totalRate = 139.00
	item.roomType.totalRateInclusive = 153.09

	if count == 0 {
		err = c.Insert(&item)
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "3"}).Count()
	if err != nil {
		log.Fatal(err)
	}

	item.hotelId = "3"
	item.code = "RACK"
	item.inDate = "2015-04-09"
	item.outDate = "2015-04-10"
	item.roomType.bookableRate = 109.00
	item.roomType.code = "KNG"
	item.roomType.roomDescription = "King sized bed"
	item.roomType.totalRate = 109.00
	item.roomType.totalRateInclusive = 123.17

	if count == 0 {
		err = c.Insert(&item)
		if err != nil {
			log.Fatal(err)
		}
	}

	// add up to 80 hotels
	for i := 7; i <= 80; i++ {
		if i%3 == 0 {
			hotel_id := strconv.Itoa(i)
			count, err = c.Find(&bson.M{"hotelid": hotel_id}).Count()
			if err != nil {
				log.Fatal(err)
			}
			end_date := "2015-04-"
			rate := 109.00
			rate_inc := 123.17
			if i%2 == 0 {
				end_date = end_date + "17"
			} else {
				end_date = end_date + "24"
			}

			if i%5 == 1 {
				rate = 120.00
				rate_inc = 140.00
			} else if i%5 == 2 {
				rate = 124.00
				rate_inc = 144.00
			} else if i%5 == 3 {
				rate = 132.00
				rate_inc = 158.00
			} else if i%5 == 4 {
				rate = 232.00
				rate_inc = 258.00
			}

			if count == 0 {

				item.hotelId = hotel_id
				item.code = "RACK"
				item.inDate = "2015-04-09"
				item.outDate = end_date
				item.roomType.bookableRate = rate
				item.roomType.code = "KNG"
				item.roomType.roomDescription = "King sized bed"
				item.roomType.totalRate = rate
				item.roomType.totalRateInclusive = rate_inc

				err = c.Insert(&item)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	err = c.EnsureIndexKey("hotelid")
	if err != nil {
		log.Fatal(err)
	}

	//RECOMMENDATION
	fmt.Printf("recommendation db ip addr = %s\n", url)
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	c = session.DB("recommendation-db").C("recommendation")
	// First we clear the collection to have always a new one
	if err = c.DropCollection(); err != nil {
		log.Print("DropCollection: ", err)
	}

	count, err = c.Find(&bson.M{"hotelid": "1"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&HotelDB{"1", 37.7867, -122.4112, 109.00, 150.00})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "2"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&HotelDB{"2", 37.7854, -122.4005, 139.00, 120.00})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "3"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&HotelDB{"3", 37.7834, -122.4071, 109.00, 190.00})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "4"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&HotelDB{"4", 37.7936, -122.3930, 129.00, 160.00})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "5"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&HotelDB{"5", 37.7831, -122.4181, 119.00, 140.00})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "6"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&HotelDB{"6", 37.7863, -122.4015, 149.00, 200.00})
		if err != nil {
			log.Fatal(err)
		}
	}

	// add up to 80 hotels
	for i := 7; i <= 80; i++ {
		hotel_id := strconv.Itoa(i)
		_, err = c.Find(bson.M{"hotelid": hotel_id}).Count()
		if err != nil {
			log.Fatal(err)
		}
		lat := 37.7835 + float64(i)/500.0*3
		lon := -122.41 + float64(i)/500.0*4

		count, err = c.Find(bson.M{"hotelid": hotel_id}).Count()
		if err != nil {
			log.Fatal(err)
		}

		rate := 135.00
		rate_inc := 179.00
		if i%3 == 0 {
			if i%5 == 0 {
				rate = 109.00
				rate_inc = 123.17
			} else if i%5 == 1 {
				rate = 120.00
				rate_inc = 140.00
			} else if i%5 == 2 {
				rate = 124.00
				rate_inc = 144.00
			} else if i%5 == 3 {
				rate = 132.00
				rate_inc = 158.00
			} else if i%5 == 4 {
				rate = 232.00
				rate_inc = 258.00
			}
		}

		if count == 0 {
			err = c.Insert(HotelDB{hotel_id, lat, lon, rate, rate_inc})
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	err = c.EnsureIndexKey("hotelid")
	if err != nil {
		log.Fatal(err)
	}

	//RESERVATION
	fmt.Printf("reservation db ip addr = %s\n", url)
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	c = session.DB("reservation-db").C("reservation")
	// First we clear the collection to have always a new one
	if err = c.DropCollection(); err != nil {
		log.Print("DropCollection: ", err)
	}

	count, err = c.Find(bson.M{"hotelid": "4"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Reservation{"4", "Alice", "2015-04-09", "2015-04-10", 1})
		if err != nil {
			log.Fatal(err)
		}
	}

	c = session.DB("reservation-db").C("number")
	count, err = c.Find(bson.M{"hotelid": "1"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(Number{"1", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(bson.M{"hotelid": "2"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(Number{"2", 10})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(bson.M{"hotelid": "3"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(Number{"3", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(bson.M{"hotelid": "4"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(Number{"4", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(bson.M{"hotelid": "5"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(Number{"5", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(bson.M{"hotelid": "6"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(Number{"6", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 7; i <= 80; i++ {
		hotel_id := strconv.Itoa(i)
		count, err = c.Find(bson.M{"hotelid": hotel_id}).Count()
		if err != nil {
			log.Fatal(err)
		}
		room_num := 200
		if i%3 == 1 {
			room_num = 300
		} else if i%3 == 2 {
			room_num = 250
		}
		if count == 0 {
			err = c.Insert(Number{hotel_id, room_num})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = c.EnsureIndexKey("hotelid")
	if err != nil {
		log.Fatal(err)
	}

	//USER
	// Read the initial users
	//users := make([]User, 500)
	//users[0].Username = "hello"
	//users[0].Password = "hello"

	// readJson(&users, *initdata)
	// fmt.Println(users)

	// Create users
	//for i := 1; i < len(users); i++ {
	//	suffix := strconv.Itoa(i)
	//	users[i].Username = "user_" + suffix
	//	users[i].Password = "pass_" + suffix
	//}

	// Insert in Database
	//coll := session.Database("user-db").Collection("user")
	// Clear any existing data

	// u := User{"Hello", "test"}
	// res, err := coll.InsertOne(ctx, u)

	// sum := sha256.Sum256([]byte(u.Password))
	// pass := fmt.Sprintf("%x", sum)

	// fmt.Println("Insert User in DB: ", u.Username, " ", u.Password, " ", pass)

	// create the database records
	//elements := make([]interface{}, len(users))
	//for i := range users {
	//	sum := sha256.Sum256([]byte(users[i].Password))
	//	pass := fmt.Sprintf("%x", sum)
	//	elements[i] = bson.M{"username": users[i].Username, "password": pass}
	//}

	// Insert them into the data base
	//opts := options.InsertMany().SetOrdered(false)
	//_, err = coll.InsertMany(elements, opts)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func main() {
	initializeDatabase()
}
