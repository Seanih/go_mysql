package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	dbConnection()

	albums, err := getAlbumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Albums found: %v\n", albums)

	album, err := getAlbumByID(1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("This is the album you've been looking for: %v\n", album)

	// id, err := addAlbum(Album{Title: "Bad", Artist: "Michael Jackson", Price: 999.99})

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("New Album ID: %v\n", id)

	allAlbums, err := getAllAlbums()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("all albums found:\n%v", allAlbums)

}

func dbConnection() {
	// load env variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DB_USERNAME"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	// connect to db
	db, err = sql.Open("mysql", cfg.FormatDSN())

	// handles errors; do something better for production
	if err != nil {
		log.Fatal(err)
	}

	// verify if db connection is live; reconnect if not
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
}

// getAlbumsByArtist queries for albums that have the specified artist name.
func getAlbumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM albums WHERE artist = ?", name)

	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	// Defer closing rows so that any resources it holds will be released when the function exits
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}

		albums = append(albums, alb)
	}

	// checks for any errors that may have occured during the loop
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func getAlbumByID(id uint64) (Album, error) {
	var alb Album

	row := db.QueryRow("SELECT * FROM albums WHERE id = ?", id)

	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
	}

	return alb, nil
}

func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO albums (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}

func getAllAlbums() ([]Album, error) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM albums")

	if err != nil {
		return nil, fmt.Errorf("there was an error fetching the albums: %v", err)
	}

	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("error getting albums: %v", err)
		}

		albums = append(albums, alb)
	}

	return albums, nil
}
