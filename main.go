package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"./model"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

func main() {
	//initiate new multiplexer
	mux := goji.NewMux()
	mux.HandleFuncC(pat.Get("/swit/:switId"), getSwit)
	mux.HandleFuncC(pat.Get("/swits"), getAllSwits)
	http.ListenAndServe("localhost:3001", mux)
}

func getAllSwits(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	//prepare for the response
	var swits []*model.Swit
	//catch them all and if there is an error
	err := getSession().DB("swit_app").C("swits").Find(nil).All(&swits)
	if err != nil {
		//let's get panic instead
		panic(err)
	}
	// _ as blank identifier,the function returns multiple values, but we don't care
	switsJSON, _ := json.Marshal(swits)
	//write a response
	ResponseWithJSON(w, switsJSON, http.StatusOK)
}

func createUser(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var user model.User
	err := decoder.Decode(user)
	if err != nil {
		panic(err)
	}

	defer req.Body.Close()
	log.Println(user.Name)
}

func getSession() *mgo.Session {
	if session == nil {
		log.Println("empty session")
		var err error
		session, err = mgo.Dial("localhost:27017")
		//check error connection, is mongo running?
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	}

	copiedSession := session.Copy()
	defer copiedSession.Clone()

	return copiedSession
}

func getSwit(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	switID := pat.Param(ctx, "switId")
	fmt.Printf("swit id is %s \n", switID)
	//get the db session
	swits := getSession().DB("swit_app").C("swits")
	var swit model.Swit
	err := swits.Find(bson.M{"switId": switID}).One(&swit)
	if err != nil {
		panic(err)
	}

	fmt.Printf("the swit object is : %s \n", swit)
	respBody, _ := json.Marshal(swit)
	ResponseWithJSON(w, respBody, http.StatusOK)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	fmt.Printf("Response Body : %q", json)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}
