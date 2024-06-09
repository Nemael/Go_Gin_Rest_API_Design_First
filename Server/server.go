// Package routes Books API.
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//	Host: localhost:8080
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//swagger:meta
package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Book represents a book with an ID, a title, an author and a quantity
// swagger:model
type Book struct {
	// The ID of the book
	// required: true
	// example: 1
	ID int `json:"id"`

	// The title of the book
	// required: true
	// example: "Safe-ish shapes for kids"
	Title string `json:"title"`

	// The author of the book
	// required: true
	// example: "Philip"
	Author string `json:"author"`

	// The quantity available of a book
	// required: true
	// example: 5
	Quantity int `json:"quantity"`
}

var (
	ConnectionString = "rest:password@tcp(localhost:3306)/books"
)

func getDB() (*sql.DB, error) {
	return sql.Open("mysql", ConnectionString)
}

//swagger:route GET /books books getbooks
//
// GetBooks returns all books.
func getBooks(c *gin.Context) {
	books := []Book{}
	db, err := getDB()
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Database connection failed."})
		return
	}
	rows, err := db.Query("SELECT id, title, author, quantity FROM books")
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Database connection failed."})
		return
	}
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Quantity)
		if err != nil {
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Database connection failed."})
		}
		books = append(books, book)
	}
	c.IndentedJSON(http.StatusOK, books) //Data sent is the books struct
}

//swagger:route GET /books/{id} books getBookById
//
// GetBookById return a book by its id. If no book is associated to this id, this route will give an error.
func getBookById(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Database connection failed."})
		return
	}
	id, err := strconv.Atoi(c.Param("id")) //Path parameter
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Id is not a valid number"}) //gin.H is a shortcut to write custom json
		return
	}
	var book Book
	row := db.QueryRow("SELECT id, title, author, quantity FROM books where id = ?", id)
	err = row.Scan(&book.ID, &book.Title, &book.Author, &book.Quantity)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."}) //gin.H is a shortcut to write custom json
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

//swagger:route POST /books books createBook
//
// CreateBook creates a book from json data and returns it. A book contains a key (int), a title (string), an author (string) and a quantity (int). All values are required.
func createBook(c *gin.Context) {
	var newBook Book
	db, err := getDB()
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Database connection failed."})
		return
	}
	if err := c.BindJSON(&newBook); err != nil { //Binds the received values to newBook
		return //In case we get an error, BindJSON gives a return response
	}
	_, err = db.Exec("INSERT INTO books (id, title, author, quantity) VALUES (?, ?, ?, ?)", newBook.ID, newBook.Title, newBook.Author, newBook.Quantity)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Error inserting into the database"})
		return
	}
	c.IndentedJSON(http.StatusCreated, newBook)
}

//swagger:route PATCH /checkout?id={id} books getBookById
//
// CheckoutBook checks out the book corresponding to the given id. Checking out a book means its quantity will be reduced by 1. If there is no copy available (quantity = 0), it is not possible to check out this book and this route will give an error.
func checkoutBook(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Database connection failed."})
		return
	}
	str_id, ok := c.GetQuery("id") //Get inline query
	id, err := strconv.Atoi(str_id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Id is not a valid number"}) //gin.H is a shortcut to write custom json
		return
	}

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	var book Book
	row := db.QueryRow("SELECT id, title, author, quantity FROM books where id = ?", id)
	err = row.Scan(&book.ID, &book.Title, &book.Author, &book.Quantity)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."}) //gin.H is a shortcut to write custom json
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."}) //gin.H is a shortcut to write custom json
		return
	}
	book.Quantity -= 1

	db.QueryRow("UPDATE books SET quantity = ? WHERE id = ?", book.Quantity, book.ID)
	c.IndentedJSON(http.StatusOK, book)
}

//swagger:route PATCH /return?id={id} books getBookById
//
// ReturnBook returns the book corresponding to the given id. Returning a book means its quantity will be increased by 1. If there is no id associated with this book, the route will give an error.
func returnBook(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Database connection failed."})
		return
	}
	str_id, ok := c.GetQuery("id") //Get inline query
	id, err := strconv.Atoi(str_id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Id is not a valid number"}) //gin.H is a shortcut to write custom json
		return
	}

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	var book Book
	row := db.QueryRow("SELECT id, title, author, quantity FROM books where id = ?", id)
	err = row.Scan(&book.ID, &book.Title, &book.Author, &book.Quantity)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."}) //gin.H is a shortcut to write custom json
		return
	}

	book.Quantity += 1

	db.QueryRow("UPDATE books SET quantity = ? WHERE id = ?", book.Quantity, book.ID)
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	// Ping database
	db, err := getDB()
	if err != nil {
		fmt.Println("Error with the database" + err.Error())
		return
	} else {
		err = db.Ping()
		if err != nil {
			fmt.Println("Error making connection tot he DB, please check credentials. The error is: " + err.Error())
			return
		}
	}
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBookById) // ':id' Setups a path parameter
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
