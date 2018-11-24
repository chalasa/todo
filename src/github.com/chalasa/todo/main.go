package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware" //

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

func main() {

	// Echo instance
	e := echo.New()

	// YML
	// mongo:
	//   host:

	// MONGO_HOST
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	mogoHost := viper.GetString("mongo.host")
	mogoUser := viper.GetString("mongo.user")
	mogoPass := viper.GetString("mongo.pass")
	port := ":" + viper.GetString("mongo.port")

	connectStr := fmt.Sprintf("%v:%v@%v", mogoUser, mogoPass, mogoHost)
	session, err := mgo.Dial(connectStr)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	h := &handler{
		m: session,
	}
	// Middleware
	e.Use(middleware.Logger())
	// Create TODO -> func create()
	e.POST("/todos", h.create)
	e.GET("/todos", h.list)
	e.GET("/todos/:id", h.view)
	e.PUT("/todos/:id", h.update)
	e.DELETE("/todos/:id", h.delete)

	// Start server
	e.Logger.Fatal(e.Start(port))
}

// Model
type todo struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Topic string        `json:"topic" bson:"topic" `
	Done  bool          `json:"done" bson:"done"`
}

type handler struct {
	m *mgo.Session
}

func (h *handler) create(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}

	t.ID = bson.NewObjectId()
	res := session.DB("workshop").C("cl-todos")
	if err2 := res.Insert(t); err2 != nil {
		return err2
	}

	return c.JSON(http.StatusOK, t)
}

// * pointer to handler
func (h *handler) list(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()
	var ts []todo //slide

	res := session.DB("workshop").C("cl-todos")
	if err2 := res.Find(nil).All(&ts); err2 != nil {
		return err2
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *handler) view(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()
	id := bson.ObjectIdHex(c.Param("id"))

	var t todo //slide
	res := session.DB("workshop").C("cl-todos")
	if err2 := res.FindId(id).One(&t); err2 != nil {
		return err2
	}

	return c.JSON(http.StatusOK, t)
}

func (h *handler) update(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()
	id := bson.ObjectIdHex(c.Param("id"))

	var t todo //slide
	res := session.DB("workshop").C("cl-todos")
	if err2 := res.FindId(id).One(&t); err2 != nil {
		return err2
	}

	t.Done = true

	if err2 := res.UpdateId(id, &t); err2 != nil {
		return err2
	}

	return c.JSON(http.StatusOK, t)
}

func (h *handler) delete(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()
	id := bson.ObjectIdHex(c.Param("id"))

	res := session.DB("workshop").C("cl-todos")
	if err2 := res.RemoveId(id); err2 != nil {
		return err2
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": "success",
	})
}

// 13.250.119.252
// root : example
// db workshop
