package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/repo"
)

// var schema = `
// CREATE TABLE person (
//     first_name text,
//     last_name text,
//     email text
// );

// CREATE TABLE place (
//     country text,
//     city text NULL,
//     telcode integer
// )`

// type Person struct {
// 	FirstName string `db:"first_name"`
// 	LastName  string `db:"last_name"`
// 	Email     string
// }

// type Place struct {
// 	Country string
// 	City    sql.NullString
// 	TelCode int
// }

func Example() {
	// this Pings the database trying to connect
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("postgres", "user=admin password=CHANGEME dbname=blog sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	// db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO posts (id, title, content) VALUES ($1, $2, $3)", 1, "Moiron", "jmoiron@jmoiron.net")
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	content := "This is a test"
	tx.NamedExec("INSERT INTO posts (id, title, content) VALUES (:ID, :Title, :Content)", &repo.PostRecord{
		ID:      2,
		Title:   "Example post",
		Content: sql.NullString{String: content, Valid: true},
	})
	tx.Commit()

	// Query the database, storing results in a []Person (wrapped in []interface{})
	posts := []repo.PostRecord{}
	db.Select(&posts, "SELECT * FROM posts ORDER BY id ASC")
	// jason, john := posts[0], posts[1]

	for _, post := range posts {
		fmt.Printf("%#v\n", post)
	}

	// fmt.Printf("%#v\n%#v", posts[0], posts[1])
	// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}
	// Person{FirstName:"John", LastName:"Doe", Email:"johndoeDNE@gmail.net"}

	// You can also get a single result, a la QueryRow
	// jason = Person{}
	// err = db.Get(&jason, "SELECT * FROM person WHERE first_name=$1", "Jason")
	// fmt.Printf("%#v\n", jason)
	// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}

	// if you have null fields and use SELECT *, you must use sql.Null* in your struct
	// places := []Place{}
	// err = db.Select(&places, "SELECT * FROM place ORDER BY telcode ASC")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// usa, singsing, honkers := places[0], places[1], places[2]

	// fmt.Printf("%#v\n%#v\n%#v\n", usa, singsing, honkers)
	// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
	// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
	// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}

	// Loop through rows using only one struct
	// place := Place{}
	// rows, err := db.Queryx("SELECT * FROM place")
	// for rows.Next() {
	// 	err := rows.StructScan(&place)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// 	fmt.Printf("%#v\n", place)
	// }
	// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
	// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}
	// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}

	// Named queries, using `:name` as the bindvar.  Automatic bindvar support
	// which takes into account the dbtype based on the driverName on sqlx.Open/Connect
	// _, err = db.NamedExec(`INSERT INTO person (first_name,last_name,email) VALUES (:first,:last,:email)`,
	// 	map[string]interface{}{
	// 		"first": "Bin",
	// 		"last":  "Smuth",
	// 		"email": "bensmith@allblacks.nz",
	// 	})

	// Selects Mr. Smith from the database
	// rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:fn`, map[string]interface{}{"fn": "Bin"})

	// Named queries can also use structs.  Their bind names follow the same rules
	// as the name -> db mapping, so struct fields are lowercased and the `db` tag
	// is taken into consideration.
	// rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:first_name`, jason)

	// batch insert

	// batch insert with structs
	// personStructs := []Person{
	// 	{FirstName: "Ardie", LastName: "Savea", Email: "asavea@ab.co.nz"},
	// 	{FirstName: "Sonny Bill", LastName: "Williams", Email: "sbw@ab.co.nz"},
	// 	{FirstName: "Ngani", LastName: "Laumape", Email: "nlaumape@ab.co.nz"},
	// }

	// _, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
	//     VALUES (:first_name, :last_name, :email)`, personStructs)

	// // batch insert with maps
	// personMaps := []map[string]interface{}{
	// 	{"first_name": "Ardie", "last_name": "Savea", "email": "asavea@ab.co.nz"},
	// 	{"first_name": "Sonny Bill", "last_name": "Williams", "email": "sbw@ab.co.nz"},
	// 	{"first_name": "Ngani", "last_name": "Laumape", "email": "nlaumape@ab.co.nz"},
	// }

	// _, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
	//     VALUES (:first_name, :last_name, :email)`, personMaps)
}
