package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"./model"
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
	//initiate new multiplexer
	mux := goji.NewMux()
	//register all handler for each end point
	mux.HandleFuncC(pat.Get("/swit/:switId"), getSwit(db))
	mux.HandleFuncC(pat.Get("/swits"), getAllSwits(db))
	mux.HandleFuncC(pat.Post("/swits"), createSwit(db))
	//finally, listen and serve in designated host and port
	http.ListenAndServe("localhost:3001", mux)
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
