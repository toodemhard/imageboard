package imageboard

import (
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Thread struct {
	Thread_id int
	Title     string `form:"title"`
	Comment   string `form:"comment"`
}

type Reply struct {
	Reply_id  int
	Thread_id int
	Comment   string `form:"comment"`
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

func CreateThread(db *sqlx.DB, thread Thread) {
	cmd := `INSERT INTO threads(title, comment) VALUES ($1,$2)`
	res, err := db.Exec(cmd, thread.Title, thread.Comment)
	_ = res
	if err != nil {
		log.Println(err.Error())
	}
}

func CreateReply(db *sqlx.DB, reply Reply) {
	cmd := `INSERT INTO replies(thread_id, comment) VALUES ($1, $2)`
	res, err := db.Exec(cmd, reply.Thread_id, reply.Comment)
	_ = res
	if err != nil {
		log.Println(err)
	}
}

func queryAllThreads(db *sqlx.DB) ([]Thread, error) {
	threads := []Thread{}
	err := db.Select(&threads, "SELECT * FROM threads")
	if err != nil {
		return threads, err
	}
	return threads, nil
}

func queryThread(db *sqlx.DB, thread_id int64) (Thread, error) {
	thread := Thread{}
	err := db.Get(&thread, `SELECT thread_id,title,comment FROM threads WHERE thread_id=$1`, thread_id)
	if err != nil {
		return thread, err
	}
	return thread, nil
}

func queryThreadReplies(db *sqlx.DB, thread_id int) ([]Reply, error) {
	replies := []Reply{}
	rows, err := db.Queryx(`SELECT * FROM replies WHERE thread_id=$1`, thread_id)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		reply := Reply{}
		err := rows.StructScan(&reply)
		if err != nil {
			return replies, err
		}
		replies = append(replies, reply)
	}
	return replies, nil
}
