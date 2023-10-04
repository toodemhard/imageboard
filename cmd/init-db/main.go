package main

import "github.com/toodemhard/imageboard/internal/imageboard"

func main() {
	db := imageboard.ConnectToDB()
	imageboard.InitSchema(db)
}
