# <B>GOLANG REST API</b>

## <b>Example for request:</b> 
### <b>Create</b>
`curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"The Go Programming Language","author":"Alan Donovan","year":2015}'`

### <b>Get All</b>
`curl http://localhost:8080/books`

### <b>Get One</b>
`curl http://localhost:8080/books/1`

### <b>Update</b>
`curl -X PUT http://localhost:8080/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","author":"New Author","year":2025}'`

### <b>Delete</b>
`curl -X DELETE http://localhost:8080/books/1`