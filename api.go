package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	obj map[string]interface{}
	arr []interface{}
)
type user struct {
	Id       bson.ObjectId `form: "id" 	    json:"id"       bson:"id"`
	Username string        `form:"login"    json:"login"    bson:"login"    binding:"required,min=3,max=20"`
	Password string        `form:"pass"     json:"pass"     bson:"password" binding:"required,min=5"`
	Age      int           `form:"age"      json:"age"      bson:"age" 	    binding:"required"`
}

type group struct {
	Id    bson.ObjectId `form: "id" 	   json:"id"         bson:"id"`
	Title string        `form:"title"      json:"title"      bson:"title"    binding:"required,min=3,max=20"`
	Users []string      `form:"users"      json:"users"      bson:"users"`
}
type userGroup struct {
	Id     string `form: "id"	json:"id"	bson:"id"	binding:"required"`
	UserId string `form:"userid"	json:"userid"	bson:"userid" binding:"required"`
}
type groupGroup struct {
	First  string `form: "first"	json:"first"	bson:"first"	binding:"required"`
	Second string `form:"second"	json:"second"	bson:"second" binding:"required"`
}

func hashPass(pass string) string {
	hasher := md5.New()
	hasher.Write([]byte(pass))
	hp := hex.EncodeToString(hasher.Sum(nil))
	return hp
}

func main() {

	gin.SetMode(gin.ReleaseMode)
	session, err := mgo.Dial("localhost")
	sG := session.DB("mydb2").C("groups")
	sU := session.DB("mydb2").C("users")
	if err != nil {
		panic(err)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	// Sample: http://localhost:4000/api/user/create/?login=Tomaki&pass=qwerty123&age=20
	router.GET("/api/user/create/", func(c *gin.Context) {

		var u user
		if err := c.Bind(&u); err != nil {
			c.JSON(400, obj{"msg": "Wrong data format", "err": err.Error()})
			return
		}
		u.Id = bson.NewObjectId()
		u.Password = hashPass(u.Password)
		if err := sU.Insert(u); err != nil {
			c.JSON(400, obj{"msg": "DB error"})
			return
		}
		c.JSON(200, u.Id)
	})
	// Sample: http://localhost:4000/api/find/Vasya
	router.GET("/api/user/find/:login", func(c *gin.Context) {
		var u user

		if err := sU.Find(bson.M{"login": c.Param("login")}).One(&u); err != nil {
			c.JSON(400, obj{"msg": "User not found", "err": err.Error()})
			return
		}

		c.JSON(200, u.Id)
	})
	// Sample: http://localhost:4000/api/user/all
	router.GET("/api/user/all", func(c *gin.Context) {

		var u []user

		if err := sU.Find(bson.M{}).All(&u); err != nil {
			c.JSON(400, obj{"msg": "User not found", "err": err.Error()})
			return
		}
		c.JSON(200, u)
	})
	// Sample: http://localhost:4000/api/user/delete/56689358f071fe0c3c145943
	router.GET("/api/user/delete/:id", func(c *gin.Context) {
		if err := sU.Remove(bson.M{"id": bson.ObjectIdHex(c.Param("id"))}); err != nil {
			c.JSON(400, obj{"msg": "User not found", "err": err.Error()})
			return
		}
		c.JSON(200, obj{"msg": "User removed"})
	})
	// Sample: http://localhost:4000/api/group/create/?title=Programmers
	router.GET("/api/group/create/", func(c *gin.Context) {
		var g group
		if err := c.Bind(&g); err != nil {
			c.JSON(400, obj{"msg": "Wrong data format", "err": err.Error()})
			return
		}
		g.Id = bson.NewObjectId()
		if err := sG.Insert(g); err != nil {
			c.JSON(400, obj{"msg": "DB error"})
			return
		}
		c.JSON(200, g.Id)
	})
	// Sample: http://localhost:4000/api/group/delete/56689358f071fe0c3c145943
	router.GET("/api/group/delete/:id", func(c *gin.Context) {
		if err := sG.Remove(bson.M{"id": bson.ObjectIdHex(c.Param("id"))}); err != nil {
			c.JSON(400, obj{"msg": "Group not found", "err": err.Error()})
			return
		}
		c.JSON(200, obj{"msg": "Group removed"})
	})
	// Sample: http://localhost:4000/api/group/all
	router.GET("/api/group/all", func(c *gin.Context) {
		var g []group
		if err := sG.Find(bson.M{}).All(&g); err != nil {
			c.JSON(400, obj{"msg": "Group not found", "err": err.Error()})
			return
		}
		c.JSON(200, g)
	})
	// Sample: http://localhost:4000/api/group/addUser/56689358f071fe0c3c145943
	router.GET("/api/group/addUser/", func(c *gin.Context) {
		var ug userGroup
		if err := c.Bind(&ug); err != nil {
			c.JSON(400, obj{"msg": "Wrong data format", "err": err.Error()})
			return
		}

		var g group
		fmt.Println(ug.Id, ug.UserId)
		if err := sG.Find(bson.M{"id": bson.ObjectIdHex(ug.Id)}).One(&g); err != nil {
			c.JSON(400, obj{"msg": "Group not found", "err": err.Error()})
			return
		}

		g.Users = append(g.Users, ug.UserId)

		if err := sG.Update(bson.M{"id": g.Id}, bson.M{"id": g.Id, "title": g.Title, "users": g.Users}); err != nil {
			c.JSON(400, obj{"msg": "Group dont updated", "err": err.Error()})
			return
		}

		c.JSON(200, obj{"msg": "User added"})
	})
	// Sample: http://localhost:4000/api/group/getUsers/56689358f071fe0c3c145943
	router.GET("/api/group/getUsers/:id", func(c *gin.Context) {

		var g group

		if err := sG.Find(bson.M{"id": bson.ObjectIdHex(c.Param("id"))}).One(&g); err != nil {
			c.JSON(400, obj{"msg": "Group not found", "err": err.Error()})
			return
		}
		var u user
		var ug []user
		i := 0
		for _, us := range g.Users {
			fmt.Println(i)
			if err := sU.Find(bson.M{"id": bson.ObjectIdHex(us)}).One(&u); err != nil {
				c.JSON(400, obj{"msg": "User not found", "err": err.Error()})
				return

			}
			ug = append(ug, u)
		}
		c.JSON(200, ug)
	})
	// Sample: http://localhost:4000/api/group/deleteUser/?Id=56689358f071fe0c3c145943&userid=5668936cf071fe0c3c145947
	router.GET("/api/group/deleteUser/", func(c *gin.Context) {

		var ug userGroup
		if err := c.Bind(&ug); err != nil {
			c.JSON(400, obj{"msg": "Wrong data format", "err": err.Error()})
			return
		}

		var g group
		if err := sG.Find(bson.M{"id": bson.ObjectIdHex(ug.Id)}).One(&g); err != nil {
			c.JSON(400, obj{"msg": "Group not found", "err": err.Error()})
			return
		}

		var tempZ int
		for i := range g.Users {
			if g.Users[i] == ug.UserId {
				tempZ = i
			}

		}
		g.Users = append(g.Users[:tempZ], g.Users[tempZ+1:]...)
		fmt.Println(g.Users)
		if err := sG.Update(bson.M{"id": g.Id}, bson.M{"id": g.Id, "title": g.Title, "users": g.Users}); err != nil {
			c.JSON(400, obj{"msg": "Group dont updated", "err": err.Error()})
			return
		}

		c.JSON(200, obj{"msg": "Deleted User"})
	})
	// Sample: http://localhost:4000/api/group/plus/?First=56689358f071fe0c3c145943&second=5668936cf071fe0c3c145947
	router.GET("/api/group/plus/", func(c *gin.Context) {

		var gg groupGroup
		if err := c.Bind(&gg); err != nil {
			c.JSON(400, obj{"msg": "Wrong data format", "err": err.Error()})
			return
		}

		var g group
		var g2 group

		if err := sG.Find(bson.M{"id": bson.ObjectIdHex(gg.First)}).One(&g); err != nil {
			c.JSON(400, obj{"msg": "User not found", "err": err.Error()})
			return
		}

		if err := sG.Find(bson.M{"id": bson.ObjectIdHex(gg.Second)}).One(&g2); err != nil {
			c.JSON(400, obj{"msg": "User not found", "err": err.Error()})
			return
		}
		for i := range g2.Users {
			g.Users = append(g.Users, g2.Users[i])
		}

		if err := sG.Update(bson.M{"id": g.Id}, bson.M{"id": g.Id, "title": g.Title, "users": g.Users}); err != nil {
			c.JSON(400, obj{"msg": "Group dont updated", "err": err.Error()})
			return
		}

		c.JSON(200, obj{"msg": "Groups concatinated"})
	})
	// Sample: http://localhost:4000/api/group/minus/?First=56689358f071fe0c3c145943&second=5668936cf071fe0c3c145947
	router.GET("/api/group/minus/", func(c *gin.Context) {

		var gg groupGroup
		if err := c.Bind(&gg); err != nil {
			c.JSON(400, obj{"msg": "Wrong data format", "err": err.Error()})
			return
		}

		var g group
		var g2 group

		if err := sG.Find(bson.M{"id": bson.ObjectIdHex(gg.First)}).One(&g); err != nil {
			c.JSON(400, obj{"msg": "Group not found", "err": err.Error()})
			return
		}

		if err := sG.Find(bson.M{"id": bson.ObjectIdHex(gg.Second)}).One(&g2); err != nil {
			c.JSON(400, obj{"msg": "Group not found", "err": err.Error()})
			return
		}

		t := -1
		for i := range g2.Users {
			for j := range g.Users {
				if g.Users[j] == g2.Users[i] {
					t = j
				}
			}
			if t != -1 {
				g.Users = append(g.Users[:t], g.Users[t+1:]...)
				t = -1
			}
		}
		if err := sG.Update(bson.M{"id": g.Id}, bson.M{"id": g.Id, "title": g.Title, "users": g.Users}); err != nil {
			c.JSON(400, obj{"msg": "Group dont updated", "err": err.Error()})
			return
		}
		c.JSON(200, obj{"msg": "Group users-Group2 users [Complete!] "})
	})
	fmt.Println("Server started at http://localhost:4000")
	router.Run(":4000")
}

/*
user.create {login: "s", pass: "p", age: y}  // возвращает {id: 111} 		†
user.delete {id: x} 														†
user.findLogin {login: "s"}            // возвращает {id: 111}				†
group.create {title: "s"}             // возвращает {id: 222}
group.delete {id: x}
group.addUser {id: x, userId: y}
group.getUsers {id: x}    // возвращает [{id: 1, login: "..."}, {id: 2, login: "..."}]
group.deleteUser  {id: x, userId: y}
group.plus(id1: x, id2: x)   // добавить в первую группу юзеров второй группы
group.minus(id1: x, id2: x)  // удалить из первой группы юзеров второй группы
*/
