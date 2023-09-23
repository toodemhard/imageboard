package imageboard

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Thread struct {
	Thread_id int
	Title     string
	Comment   string
}

type Reply struct {
	Reply_id  int
	Thread_id int
	Comment   string
}

func ConnectToDB() *sqlx.DB {
	return sqlx.MustConnect("pgx", os.Getenv("DATABASE_URL"))
}

func InitSchema(db *sqlx.DB) {
	schema := `CREATE TABLE threads (
        thread_id bigserial primary key,
        title varchar(50),
        comment varchar(500)
        );
        CREATE TABLE replies(
            reply_id bigserial primary key,
            thread_id bigserial references threads(thread_id),
            comment varchar(500)
        );`

	db.MustExec(schema)
}

func createPost(db *sqlx.DB, title string, comment string) {
	thread := `insert into threads(title, comment) values ($1,$2)`
	res, err := db.Exec(thread, title, comment)
	_ = res
	if err != nil {
		log.Println(err)
	}
}

func createReply(db *sqlx.DB, thread_id int, comment string) {
	reply := `insert into replies(thread_id, comment) values ($1, $2)`
	res, err := db.Exec(reply, thread_id, comment)
	_ = res
	if err != nil {
		log.Println(err)
	}
}

func getAllPosts(db *sqlx.DB) {
	threads := []Thread{}
	err := db.Select(&threads, "select * from threads")
	if err != nil {
		log.Println(err)
	}
	for _, p := range threads {
		fmt.Println(p)
	}
}

func getPostReplies(db *sqlx.DB, thread_id int) {
	rows, err := db.Queryx(`select * from replies where thread_id=$1`, thread_id)
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
