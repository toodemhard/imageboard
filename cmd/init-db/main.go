package main

import "github.com/toodemhard/imageboard/pkg/imageboard"

func main() {
	db := imageboard.ConnectToDB()
	imageboard.InitSchema(db)
}
