# Important
- At first, the API is not connected to a database, it is only in-memory. I then connected it to a mysql database.
- `go mod init [name_of_folder]` to init a project in go

# Dependencies needed
- Setup 'Dependency tracking', `go mod init example/Rest_API_Go_Gin`to create go.mod file
- `go get github.com/gin-gonic/gin`
- ```import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)```


# Struct setup
- Name needs to start with a capital letter
- Make the struct serializable, to convert it to json in both ways. Upper case names makes it an 'exported field' (public), the json version we want it in lowercase
    - You can make it serializable by adding ``json:"field"`` after the field

# Gin details
- `func getBooks(c *gin.Context) {c.IndentedJSON(http.StatysOK, books)}`
    - The context is related to everything needed in the context of our request

# Curl commands
- `curl localhost:8080/books` to get books
- `curl localhost:8080/books --header "Content-Type: application/json" -d @data.json --request "POST"` to add books from data.json
- `curl localhost:8080/checkout?id=2 --request "PATCH"` to checkout a book

# Swagger
- `swagger generate spec -o ./swagger.json`
- `swagger serve -F=swagger swagger.json`

# Things to remember
- I can use keep-alive for http client to re-use the connection within the timeout

# Still need to make:
- /Go client
- Add API Keys/Oauth2?
- /Connect to mysql database
- Swagger implementation
- Swagger is not completely implemented, still having troube with parameters and the "try it out" button
- Upload it all on github/gitlab
- UNDERSTAND THE DIFFERENCE BETWEEN RESTFUL AND GRPC
- Fuse all my .md documents
- Clean my REST API Server code (a lot of repeating code can be but into functions)
- Move my REST API Server code to MVC structure (model/view/controller) model will have the book struct, controller will have the many methods to get the data, and view will have the router setting, I guess?
- Clean my REST API Client
- It would be possible to use concurrency for better API calls and database access (threads?)
