package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware" //

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// func main() {
// 	// fmt.Println(name.Name)
// 	// name.Show()
// }

// func hello(w http.ResponseWriter, r *http.Request) {
// 	io.WriteString(w, "Hello world!")
// }

func main() {
	// // localhost:8000/bar
	// http.HandleFunc("/bar", hello)
	// http.ListenAndServe(":8000", nil)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// Route => handler
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	// Create TODO -> func create()
	e.POST("/todos", create)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Model
type todo struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Topic string        `json:"topic" bson:"topic" `
	Done  bool          `json:"done" bson:"done"`
}

func create(c echo.Context) error {
	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}

	session, err := mgo.Dial("root:example@13.250.119.252")

	if err != nil {
		return err
	}

	res := session.DB("workshop").C("cl-todos")
	if err2 := res.Insert(t); err != nil {
		return err2
	}
	return c.JSON(http.StatusOK, t)
}

// 13.250.119.252
// root : example
// db workshop
