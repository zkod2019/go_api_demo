package main

import (
	"errors"
	"net/http" // built into go

	"github.com/gin-gonic/gin"
)

type book struct {
	// names are uppercase so that they can be exported (or viewed) by outside modules
	// json names are lowercase ("when serializing w/ json, convert this field to lowercase")
	ID       string `json:"id"` // using json for api (returning and recieving)
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{ // list ("slice") of books. Data structure to represent all books (instead of using DB)
	{ID: "1", Title: "The Things They Carried", Author: "Tim O'Brien", Quantity: 4},
	{ID: "2", Title: "Sharp Objects", Author: "Gillian Flynn", Quantity: 2},
	{ID: "3", Title: "Reasons the Breathe", Author: "Rebecca Donovan", Quantity: 3},
}

func getBooks(c *gin.Context) { // return json obj of all books
	// gin.Context = all the info about the request & allows you to return a response
	c.IndentedJSON(http.StatusOK, books) // indented JSON = nicely formatted JSON
	// ^ w/ that we are serializing the book struct & a slice of the []book.
}

// contex --> 4 func ()
func createBook(c *gin.Context) { // gin.Context has query parameters, data payload, headers, and anything else u have
	var newBook book                             // creating a new book of type book
	if err := c.BindJSON(&newBook); err != nil { // bind json to the new book (its pointer to new book so we directly modify field values)
		return // will return the error that BindJSON gets
	}
	// append to the book slive
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook) // return the book that we just created w/ the status code
}

func bookById(c *gin.Context) {
	id := c.Param("id") // param is path parameter, ex. "/books/2" (where 2 is the id)
	// ^ will access whatever is associated w/ the id paramter passed from the router.GET("/books/:id")
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."}) // gin.H maps string to any interface
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	// if id of book isnt valid ^ return error
	for i, b := range books {
		// loop through all books
		if b.ID == id {
			return &books[i], nil
			// return pointer so we can modify attributes of the book or fields of struct from diff func
		}
	}
	return nil, errors.New("book not found")
}

// check out a book reduces quantity by one (doesn't let checkout if quantity=0)
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing ID query parameter"})
		return
	}

	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."}) // gin.H maps string to any interface
		return
	}
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."}) // gin.H maps string to any interface
		return
	}
	book.Quantity -= 1 // i fall conditions are true, wee reduce quantity
	c.IndentedJSON(http.StatusOK, book)
	// for i, b := range books {
	// 	if b.ID == id {
	// 		&book[i].Quantity - 1
	// 	}
	// }
}

// return a book add a quantity back
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing ID query parameter"})
		return
	}

	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."}) // gin.H maps string to any interface
		return
	}
	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

// curl localhost:8080/books --include --header "Content-Type: application/json" -d @body.json --request "POST"
// ^we r curling to localhost..., will include header "application/json" --> defining data we r sending. -d = data. @ = file. request type = POST
func main() {
	router := gin.Default()            // creating gin router
	router.GET("/books", getBooks)     // creating endpoint --> when localhost:8080/books is opened, getBooks func will be called
	router.GET("/books/:id", bookById) //":id" will defualt to string type
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")

	// GET = getting info
	// POST = adding info, or creating something new
	// PATCH = updating something

}
