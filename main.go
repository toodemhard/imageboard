package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Post struct {
	Post_id int
	Title   string
	Content string
}

type Reply struct {
	Reply_id int
	Post_id  int
	Content  string
}

func initSchema(db *sqlx.DB) {
	schema := `CREATE TABLE posts (
        post_id BIGSERIAL PRIMARY KEY,
        title VARCHAR(50),
        content VARCHAR(500)
        );
        CREATE TABLE replies(
            reply_id BIGSERIAL PRIMARY KEY,
            post_id BIGSERIAL references posts(post_id),
            content VARCHAR(500)
        );`

	db.MustExec(schema)
}

func createPost(db *sqlx.DB, title string, content string) {
	post := `insert into posts(title, content) values ($1,$2)`
	res, err := db.Exec(post, title, content)
	_ = res
	if err != nil {
		log.Println(err)
	}
}

func createReply(db *sqlx.DB, post_id int, content string) {
	reply := `insert into replies(post_id, content) values ($1, $2)`
	res, err := db.Exec(reply, post_id, content)
	_ = res
	if err != nil {
		log.Println(err)
	}
}

func getAllPosts(db *sqlx.DB) {
	posts := []Post{}
	err := db.Select(&posts, "select * from posts")
	if err != nil {
		log.Println(err)
	}
	for _, p := range posts {
		fmt.Println(p)
	}
}

func getPostReplies(db *sqlx.DB, post_id int) {
	rows, err := db.Queryx(`select * from replies where post_id=$1`, post_id)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		reply := Reply{}
		err := rows.StructScan(&reply)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(reply)
	}
}

func main() {
	db := sqlx.MustConnect("pgx", os.Getenv("DATABASE_URL"))
	_ = db
}
