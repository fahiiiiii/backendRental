// database/database.go
package database

import (
    "fmt"
    "log"

    "github.com/beego/beego/v2/client/orm"
    _ "github.com/lib/pq" // PostgreSQL driver
    beego "github.com/beego/beego/v2/server/web"
)

// InitializeDatabase sets up the database connection
func InitializeDatabase() {
    // Retrieve database configuration
    dbDriver, err := beego.AppConfig.String("dbdriver")
    if err != nil || dbDriver == "" {
        log.Println("Database driver not specified, defaulting to postgres")
        dbDriver = "postgres"
    }

    dbHost, err := beego.AppConfig.String("dbhost")
    if err != nil {
        log.Fatal("Database host not specified in app.conf")
    }

    dbUser, err := beego.AppConfig.String("dbuser")
    if err != nil {
        log.Fatal("Database user not specified in app.conf")
    }

    dbPassword, err := beego.AppConfig.String("dbpassword")
    if err != nil {
        log.Fatal("Database password not specified in app.conf")
    }

    dbName, err := beego.AppConfig.String("dbname")
    if err != nil {
        log.Fatal("Database name not specified in app.conf")
    }

    dbPort, err := beego.AppConfig.String("dbport")
    if err != nil {
        log.Fatal("Database port not specified in app.conf")
    }

    // Construct the connection string
    dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    // Register the PostgreSQL driver explicitly
    err = orm.RegisterDriver("postgres", orm.DRPostgres)
    if err != nil {
        log.Fatalf("Failed to register PostgreSQL driver: %v", err)
    }

    // Register the database
    err = orm.RegisterDataBase("default", "postgres", dataSource)
    if err != nil {
        log.Fatalf("Failed to register database: %v", err)
    }

    // Enable ORM debug mode (optional)
    orm.Debug = true

    // Log successful connection
    log.Println("Database connection established successfully")

    // Auto create tables (optional)
    err = orm.RunSyncdb("default", false, true)
    if err != nil {
        log.Fatalf("Failed to sync database schema: %v", err)
    }
}