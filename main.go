package main

import (
	"crypto/md5"
	"encoding/hex"
	//"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"strconv"
)

var usersCollection *mgo.Collection

func genId() int {
	id, err := usersCollection.Count()
	if err != nil {
		fmt.Println(err.Error())
	}
	var res User
	var idI int
	for i := 0; i <= id; i++ {
		err := usersCollection.Find(bson.M{"id": i}).One(&res)
		if err != nil {
			idI = i
			break
		}
	}
	return idI
}

type User struct {
	Id       int    `json:" id"`
	Login    string `json:" login"`
	Password string `json:" password"`
	Age      int    `json:" age"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	userDocuments := []User{}
	err := usersCollection.Find(bson.M{}).All(&userDocuments)
	if err != nil {
		fmt.Println(err.Error())
	}
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	t.ExecuteTemplate(w, "index", userDocuments)
}
func regHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("reg.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "reg", nil)
}
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	rmap := r.URL.Query()
	login := rmap.Get("login")
	password := rmap.Get("pass")
	hasher := md5.New()
	hasher.Write([]byte(password))
	hashPass := hex.EncodeToString(hasher.Sum(nil))
	tempAge, _ := strconv.ParseInt(rmap.Get("age"), 10, 0)
	age := int(tempAge)
	userDocument := User{genId(), login, hashPass, age}
	err := usersCollection.Insert(userDocument)
	if err != nil {
		fmt.Println(err.Error())
	}
	http.Redirect(w, r, "/", 302)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	rmap := r.URL.Query()
	tempId, _ := strconv.ParseInt(rmap.Get("id"), 10, 0)
	id := int(tempId)
	err := usersCollection.Remove(bson.M{"id": id})
	if err != nil {
		fmt.Println(err.Error())
	}
	http.Redirect(w, r, "/", 302)
}

func main() {
	fmt.Println("Listening on port :4000")
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	usersCollection = session.DB("learn").C("newBaseDD")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/reg", regHandler)
	http.HandleFunc("/user.create", createUserHandler)
	http.HandleFunc("/user.delete", deleteUserHandler)
	http.ListenAndServe(":4000", nil)
}
