package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// "gorm.io/gorm"
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // Docker service name of the PostgreSQL container
	port     = 5433         // Internal PostgreSQL port inside the container
	user     = "admin"      // From environment POSTGRES_USER
	password = "password"   // From environment POSTGRES_PASSWORD
	dbname   = "mydatabase" // From environment POSTGRES_DB
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			// IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			// ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful: true, // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&Book{}, &Publisher{}, &Author{}, &AuthorBook{})

	// ขาสร้าง
	publisher := Publisher{
		Details: "Publisher Details",
		Name:    "Publisher Name",
	}
	_ = createPublisher(db, &publisher)

	// Example data for a new author
	author := Author{
		Name: "Author Name",
	}
	_ = createAuthor(db, &author)

	// // Example data for a new book with an author
	book := Book{
		Name:        "Book Title",
		Author:      "Book Author",
		Description: "Book Description",
		PublisherID: publisher.ID,     // Use the ID of the publisher created above
		Authors:     []Author{author}, // Add the created author
	}
	_ = createBookWithAuthor(db, &book, []uint{author.ID})

	// ขาเรียก

	// Example: Get a book with its publisher
	bookWithPublisher, err := getBookWithPublisher(db, 1) // assuming a book with ID 1
	if err != nil {
		// Handle error
	}

	// Example: Get a book with its authors
	bookWithAuthors, err := getBookWithAuthors(db, 1) // assuming a book with ID 1
	if err != nil {
		// Handle error
	}

	// Example: List books of a specific author
	authorBooks, err := listBooksOfAuthor(db, 1) // assuming an author with ID 1
	if err != nil {
		// Handle error
	}

	fmt.Println(bookWithPublisher)
	fmt.Println(bookWithAuthors)
	fmt.Println(authorBooks)
}
