package main

import (
	"fmt"
	"strconv"
	"context"
	"crypto/sha256"

	log "github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	bson2 "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

type Hotel struct {
        id string
        name string
        phoneNumber string
        description string
        address Address
        images []Image
}

func initializeDatabase()  {
	url := "http://hotel-app-database.default.10.125.188.36.sslip.io"

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

	c := session.DB("profile-db").C("hotels")
	// First we clear the collection to have always a new one
	if err = c.DropCollection(); err != nil {
		log.Print("DropCollection: ", err)
	}

	count, err := c.Find(&bson.M{"id": "1"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Hotel{
			Id:          "1",
			Name:        "Clift Hotel",
			PhoneNumber: "(415) 775-4700",
			Description: "A 6-minute walk from Union Square and 4 minutes from a Muni Metro station, this luxury hotel designed by Philippe Starck features an artsy furniture collection in the lobby, including work by Salvador Dali.",
			Address: &pb.Address{
				StreetNumber: "495",
				StreetName:   "Geary St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94102",
				Lat:          37.7867,
				Lon:          -122.4112}})
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
			Id:          "2",
			Name:        "W San Francisco",
			PhoneNumber: "(415) 777-5300",
			Description: "Less than a block from the Yerba Buena Center for the Arts, this trendy hotel is a 12-minute walk from Union Square.",
			Address: &pb.Address{
				StreetNumber: "181",
				StreetName:   "3rd St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94103",
				Lat:          37.7854,
				Lon:          -122.4005}})
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
			Id:          "3",
			Name:        "Hotel Zetta",
			PhoneNumber: "(415) 543-8555",
			Description: "A 3-minute walk from the Powell Street cable-car turnaround and BART rail station, this hip hotel 9 minutes from Union Square combines high-tech lodging with artsy touches.",
			Address: &pb.Address{
				StreetNumber: "55",
				StreetName:   "5th St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94103",
				Lat:          37.7834,
				Lon:          -122.4071}})

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
			Id:          "4",
			Name:        "Hotel Vitale",
			PhoneNumber: "(415) 278-3700",
			Description: "This waterfront hotel with Bay Bridge views is 3 blocks from the Financial District and a 4-minute walk from the Ferry Building.",
			Address: &pb.Address{
				StreetNumber: "8",
				StreetName:   "Mission St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94105",
				Lat:          37.7936,
				Lon:          -122.3930}})

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
			Id:          "5",
			Name:        "Phoenix Hotel",
			PhoneNumber: "(415) 776-1380",
			Description: "Located in the Tenderloin neighborhood, a 10-minute walk from a BART rail station, this retro motor lodge has hosted many rock musicians and other celebrities since the 1950s. It’s a 4-minute walk from the historic Great American Music Hall nightclub.",
			Address: &pb.Address{
				StreetNumber: "601",
				StreetName:   "Eddy St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94109",
				Lat:          37.7831,
				Lon:          -122.4181}})

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
			Id:          "6",
			Name:        "St. Regis San Francisco",
			PhoneNumber: "(415) 284-4000",
			Description: "St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
			Address: &pb.Address{
				StreetNumber: "125",
				StreetName:   "3rd St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94109",
				Lat:          37.7863,
				Lon:          -122.4015}})

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
		lat := 37.7835 + float32(i)/500.0*3
		lon := -122.41 + float32(i)/500.0*4
		if count == 0 {
			err = c.Insert(&Hotel{
				Id:          hotel_id,
				Name:        "St. Regis San Francisco",
				PhoneNumber: phone_num,
				Description: "St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
				Address: &pb.Address{
					StreetNumber: "125",
					StreetName:   "3rd St",
					City:         "San Francisco",
					State:        "CA",
					Country:      "United States",
					PostalCode:   "94109",
					Lat:          lat,
					Lon:          lon}})

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

	c := session.DB("rate-db").C("inventory")
	// First we clear the collection to have always a new one
	if err = c.DropCollection(); err != nil {
		log.Print("DropCollection: ", err)
	}
	count, err := c.Count()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" Collection contains: %d elements\n", count)
	count, err = c.Find(&bson.M{"hotelid": "2"}).Count()
	if err != nil {
		log.Fatal(err)
	}

	item := RatePlan{
		HotelId: "1",
		Code:    "RACK",
		InDate:  "2015-04-09",
		OutDate: "2015-04-10",
		RoomType: &pb.RoomType{
			BookableRate:       109.00,
			Code:               "KNG",
			RoomDescription:    "King sized bed",
			TotalRate:          109.00,
			TotalRateInclusive: 123.17,
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

	item.HotelId = "2"
	item.Code = "RACK"
	item.InDate = "2015-04-09"
	item.OutDate = "2015-04-10"
	item.RoomType.BookableRate = 139.00
	item.RoomType.Code = "QN"
	item.RoomType.RoomDescription = "Queen sized bed"
	item.RoomType.TotalRate = 139.00
	item.RoomType.TotalRateInclusive = 153.09

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

	item.HotelId = "3"
	item.Code = "RACK"
	item.InDate = "2015-04-09"
	item.OutDate = "2015-04-10"
	item.RoomType.BookableRate = 109.00
	item.RoomType.Code = "KNG"
	item.RoomType.RoomDescription = "King sized bed"
	item.RoomType.TotalRate = 109.00
	item.RoomType.TotalRateInclusive = 123.17

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

				item.HotelId = hotel_id
				item.Code = "RACK"
				item.InDate = "2015-04-09"
				item.OutDate = end_date
				item.RoomType.BookableRate = rate
				item.RoomType.Code = "KNG"
				item.RoomType.RoomDescription = "King sized bed"
				item.RoomType.TotalRate = rate
				item.RoomType.TotalRateInclusive = rate_inc

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

	c := session.DB("recommendation-db").C("recommendation")
	// First we clear the collection to have always a new one
	if err = c.DropCollection(); err != nil {
		log.Print("DropCollection: ", err)
	}

	count, err := c.Find(&bson.M{"hotelid": "1"}).Count()
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
		_, err = c.Find(&bson.M{"hotelid": hotel_id}).Count()
		if err != nil {
			log.Fatal(err)
		}
		lat := 37.7835 + float64(i)/500.0*3
		lon := -122.41 + float64(i)/500.0*4

		count, err = c.Find(&bson.M{"hotelid": hotel_id}).Count()
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
			err = c.Insert(&HotelDB{hotel_id, lat, lon, rate, rate_inc})
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

	c := session.DB("reservation-db").C("reservation")
	// First we clear the collection to have always a new one
	if err = c.DropCollection(); err != nil {
		log.Print("DropCollection: ", err)
	}

	count, err := c.Find(&bson.M{"hotelid": "4"}).Count()
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
	count, err = c.Find(&bson.M{"hotelid": "1"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Number{"1", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "2"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Number{"2", 10})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "3"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Number{"3", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "4"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Number{"4", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "5"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Number{"5", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelid": "6"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		err = c.Insert(&Number{"6", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 7; i <= 80; i++ {
		hotel_id := strconv.Itoa(i)
		count, err = c.Find(&bson.M{"hotelid": hotel_id}).Count()
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
			err = c.Insert(&Number{hotel_id, room_num})
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
	users := make([]User, 500)
	users[0].Username = "hello"
	users[0].Password = "hello"

	// readJson(&users, *initdata)
	// fmt.Println(users)

	// Create users
	for i := 1; i < len(users); i++ {
		suffix := strconv.Itoa(i)
		users[i].Username = "user_" + suffix
		users[i].Password = "pass_" + suffix
	}

	// Insert in Database
	coll := c.Database("user-db").Collection("user")
	// Clear any existing data

	if err := coll.Drop(ctx); err != nil {
		log.Print("DropCollection: ", err)
	}

	// u := User{"Hello", "test"}
	// res, err := coll.InsertOne(ctx, u)

	// sum := sha256.Sum256([]byte(u.Password))
	// pass := fmt.Sprintf("%x", sum)

	// fmt.Println("Insert User in DB: ", u.Username, " ", u.Password, " ", pass)

	// create the database records
	elements := make([]interface{}, len(users))
	for i := range users {
		sum := sha256.Sum256([]byte(users[i].Password))
		pass := fmt.Sprintf("%x", sum)
		elements[i] = bson.M{"username": users[i].Username, "password": pass}
	}

	// Insert them into the data base
	opts := options.InsertMany().SetOrdered(false)
	_, err := coll.InsertMany(ctx, elements, opts)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initializeDatabase()
}
