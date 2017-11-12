package main

import "database/sql"
import "log"
import "fmt"
import _ "github.com/lib/pq"

// Wrapper around postgres interactions
type PostgresClient struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Db       *sql.DB
}

func (p *PostgresClient) GetDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.Dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

type Image struct {
	ImageId  int    // ID from postgres
	Filename string // Filename relative to the configured working directory path
}

// Returns a slice of absolute paths to unclassified images
func (p *PostgresClient) GetUnclassified() []*Image {
	sqlStatement := `
	SELECT image_id, filename FROM images WHERE classified IS false`

	rows, err := p.Db.Query(sqlStatement)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	ret := make([]*Image, 0)
	var (
		imageId  int
		filename string
	)
	for rows.Next() {
		if err := rows.Scan(&imageId, &filename); err != nil {
			log.Fatal(err)
		}
		img := &Image{ImageId: imageId, Filename: filename}
		ret = append(ret, img)

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return ret
}

// Saves classification for image
// TODO: This whole thing should be in a transaction
func (p *PostgresClient) SaveClassification(imageId int, classId string,
	className string, probability float64) {

	// Write classification to db
	sqlStatement := `  
	  INSERT INTO classifications (image_id, class_id, class_name, probability)
	  VALUES ($1, $2, $3, $4)
	`
	_, err := p.Db.Exec(sqlStatement, imageId, classId, className, probability)
	if err != nil {
		panic(err)
	}

	// Update classified status of image
	sqlStatement = `  
	UPDATE images SET classified=true WHERE image_id=$1
	`
	_, err = p.Db.Exec(sqlStatement, imageId)
	if err != nil {
		panic(err)
	}
}

func NewPostgresClient(pgHost string, pgPort int, pgUser string,
	pgPassword string, pgDbname string) *PostgresClient {
	p := new(PostgresClient)
	p.Host = pgHost
	p.Port = pgPort
	p.User = pgUser
	p.Password = pgPassword
	p.Dbname = pgDbname
	p.Db = p.GetDB()
	p.Db.SetMaxOpenConns(50)
	return p
}
