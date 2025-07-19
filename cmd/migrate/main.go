package main

import (
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql" // Renamed alias for clarity
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/keshav78-78/ECOM/config"
	"github.com/keshav78-78/ECOM/db" // Assuming this package exists and NewMySQLStorage returns *sql.DB
)

func main() {
	// Step 1: Create a database connection (*sql.DB)
	// I am assuming your db.NewMySQLStorage function returns a standard *sql.DB object.
	// The configuration is now using the correct mysql.Config from the go-sql-driver.
	dbConn, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// We need to ping the database to ensure the connection is alive before proceeding.
	if err := dbConn.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	// Step 2: Create a new migrate driver instance from the database connection.
	// The function is WithInstance, and it takes the *sql.DB object.
	driver, err := mysqlMigrate.WithInstance(dbConn, &mysqlMigrate.Config{})
	if err != nil {
		log.Fatal("Could not create migrate driver:", err)
	}

	// Step 3: Create a new migrate instance.
	// The function name is NewWithDatabaseInstance.
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations", // Source of your migration files
		"mysql",                         // The name of the database
		driver,                          // The migrate driver we just created
	)
	if err != nil {
		log.Fatal("Could not create migrate instance:", err)
	}

	// Step 4: Check for command-line arguments safely.
	// We should check if an argument was actually provided.
	if len(os.Args) < 2 {
		log.Fatal("Expected 'up' or 'down' command")
	}
	cmd := os.Args[1] // The command is the first argument after the program name.

	if cmd == "up" {
		log.Println("Running migrations up...")
		// m.Up() applies all available migrations.
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("An error occurred while running migration up:", err)
		}
		log.Println("Migrations ran successfully!")
	}

	if cmd == "down" {
		log.Println("Running migrations down...")
		// The method is Down(), with a capital 'D'.
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("An error occurred while running migration down:", err)
		}
		log.Println("Migrations rolled back successfully!")
	}
}
