package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"./model"
	"github.com/rs/cors"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	//trying to connect to mongo
	db, err := mgo.Dial("localhost:27017")
	if err != nil {
		//something happens when dialing
		log.Fatal("Cannot dial mongo", err)
	}
	//clear when we're done
	defer db.Close()
	ensureIndex(db)
	//initiate new multiplexer
	mux := goji.NewMux()
	//register all handler for each end point
	mux.HandleFuncC(pat.Get("/swits/:switId"), getSwit(db))
	mux.HandleFuncC(pat.Get("/swits"), getAllSwits(db))
	mux.HandleFuncC(pat.Post("/swits"), createSwit(db))
	mux.HandleFuncC(pat.Post("/users"), addUser(db))
	mux.HandleFuncC(pat.Get("/users/:uid"), getUser(db))
	//to allow cross origin
	handler := cors.Default().Handler(mux)
	//finally, listen and serve in designated host and port
	http.ListenAndServe(":3001", handler)
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("swit_app").C("users")

	index := mgo.Index{
		Key:        []string{"uid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		log.Fatal(err)
	}
}

func addUser(s *mgo.Session) goji.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		//create a copy of a session
		session := s.Copy()
		//clear the copied session once it's done
		defer session.Close()
		//get the body request
		var user model.User
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&user)
		if err != nil {
			fmt.Println("Incorrect body of request")
			return
		}
		c := session.DB("swit_app").C("users")
		//add new user to database
		err = c.Insert(user)
		if err != nil {
			if mgo.IsDup(err) {
				//check if it's duplicate, if so return something to notify
				ResponseSimpleMessage("User is exist", false, w)
				return
			}
			//failed to Insert
			fmt.Println("Failed to insert a new user")
			log.Fatal(err)
		}

		ResponseSimpleMessage("Successfully added a new user", true, w)
	}
}

func getUser(s *mgo.Session) goji.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		//create a copy of a session
		session := s.Copy()
		//clear the copied session once it's done
		defer session.Close()
		//get query param
		uid := pat.Param(ctx, "uid")
		//get collection of users
		c := session.DB("swit_app").C("users")
		//prepare the model
		var user model.User
		//fetch the user in db
		err := c.Find(bson.M{"uid": uid}).One(&user)
		if err != nil {
			//failed to Insert
			fmt.Printf("Unable to search user with uid: %s", uid)
			log.Fatal(err)
		}
		//toJson
		respBody, _ := json.Marshal(user)
		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func createSwit(s *mgo.Session) goji.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		//create a copy of a session
		session := s.Copy()
		//clear the copied session once it's done
		defer session.Close()
		//get the body request
		var swit model.Swit
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&swit)
		if err != nil {
			fmt.Println("Incorrect body of request")
			return
		}

		//should be now
		swit.Time = time.Now()
		swit.SwitId = bson.NewObjectIdWithTime(swit.Time)
		//TODO: should get the user id from req too
		swit.UserId = bson.NewObjectId()
		fmt.Println("incoming swit: ", swit)

		c := session.DB("swit_app").C("swits")
		//insert the new swit and catch error
		err = c.Insert(swit)
		if err != nil {
			//failed to Insert
			fmt.Println("Failed to insert a new book")
			log.Fatal(err)
		}
		resBody := SimpleResponse{"Successfully added a new book", true}
		fmt.Println("res body content: ", resBody)
		//toJson
		json, _ := json.Marshal(resBody)
		//write a response back
		ResponseWithJSON(w, json, http.StatusOK)
	}
}

func getAllSwits(s *mgo.Session) goji.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		//create a copy of a session
		session := s.Copy()
		//clear the copied session once it's done
		defer session.Close()
		//prepare for the response
		var swits []*model.Swit
		//catch them all and if there is an error
		err := session.DB("swit_app").C("swits").Find(nil).All(&swits)
		if err != nil {
			//let's get panic instead
			panic(err)
		}
		// _ as blank identifier,the function returns multiple values, but we don't care
		switsJSON, _ := json.Marshal(swits)
		//write a response
		ResponseWithJSON(w, switsJSON, http.StatusOK)
	}
}

func getSwit(s *mgo.Session) goji.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		//create a copy of a session
		session := s.Copy()
		//clear the copied session once it's done
		defer session.Close()
		//get query param out of the http request
		switID := pat.Param(ctx, "switId")
		fmt.Printf("swit id is %s \n", switID)
		//get swits collection
		c := session.DB("swit_app").C("swits")
		var swit model.Swit
		//convert to object id
		oid := bson.ObjectIdHex(switID)
		err := c.Find(bson.M{"switId": oid}).One(&swit)
		if err != nil {
			panic(err)
		}

		fmt.Printf("the swit object is : %s \n", swit)
		//toJson
		respBody, _ := json.Marshal(swit)
		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func ResponseSimpleMessage(message string, success bool, w http.ResponseWriter) {
	resBody := SimpleResponse{message, success}
	fmt.Println("res body content: ", resBody)
	//toJson
	json, _ := json.Marshal(resBody)
	//write a response back
	ResponseWithJSON(w, json, http.StatusOK)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	fmt.Println("Response Body : ", string(json))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

type SimpleResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
