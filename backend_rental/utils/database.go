// utils/database.go
package utils

import (
    "fmt"
    "log"
    "os"
    "time"
    "github.com/beego/beego/v2/client/orm"
    _ "github.com/lib/pq"
    beego "github.com/beego/beego/v2/server/web"
    "backend_rental/models"
)

var modelsRegistered = false

func InitDatabase() {
    // Ensure models are only registered once
    if modelsRegistered {
        return
    }

    // Get database configuration from environment variables first, fall back to config file
    dbHost := os.Getenv("DB_HOST")
    if dbHost == "" {
        dbHost = beego.AppConfig.DefaultString("dbhost", "db")
    }

    dbPort := os.Getenv("DB_PORT")
    if dbPort == "" {
        dbPort = beego.AppConfig.DefaultString("dbport", "5432")
    }

    dbUser := os.Getenv("DB_USER")
    if dbUser == "" {
        dbUser = beego.AppConfig.DefaultString("dbuser", "")
    }

    dbPassword := os.Getenv("DB_PASSWORD")
    if dbPassword == "" {
        dbPassword = beego.AppConfig.DefaultString("dbpassword", "")
    }

    dbName := os.Getenv("DB_NAME")
    if dbName == "" {
        dbName = beego.AppConfig.DefaultString("dbname", "")
    }

    // Log configuration (excluding password)
    log.Printf("Database Configuration:")
    log.Printf("Host: %s", dbHost)
    log.Printf("Port: %s", dbPort)
    log.Printf("User: %s", dbUser)
    log.Printf("Database: %s", dbName)

    // Construct the connection string
    connectionString := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName,
    )

    // Register the database driver
    if err := orm.RegisterDriver("postgres", orm.DRPostgres); err != nil {
        log.Fatalf("Database driver registration failed: %v", err)
    }

    // Register the database with retry logic
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        err := orm.RegisterDataBase("default", "postgres", connectionString)
        if err == nil {
            log.Println("Successfully connected to database")
            break
        }
        if i == maxRetries-1 {
            log.Fatalf("Failed to register database after %d attempts: %v", maxRetries, err)
        }
        log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
        // Wait for 5 seconds before retrying
        time.Sleep(5 * time.Second)
    }

    // Register models only once
    orm.RegisterModel(
        new(models.Location),
        // Add other models here
    )
    modelsRegistered = true

    // Sync database schema
    if err := orm.RunSyncdb("default", false, true); err != nil {
        log.Fatalf("Database schema sync failed: %v", err)
    }

    // Optional: Enable ORM debug mode in development
    if beego.AppConfig.DefaultString("runmode", "dev") == "dev" {
        orm.Debug = true
    }

    log.Println("Database initialized successfully!")
}
