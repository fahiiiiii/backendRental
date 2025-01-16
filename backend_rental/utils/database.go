package utils

import (
    "fmt"
    "log"
    "github.com/beego/beego/v2/server/web"
    "github.com/beego/beego/v2/client/orm"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    _ "github.com/lib/pq"
    "backend_rental/models"
)

var db *gorm.DB

// InitDatabaseFromConfig initializes both GORM and Beego ORM
func InitDatabaseFromConfig() error {
    // Load configuration from app.conf
    err := web.LoadAppConfig("ini", "conf/app.conf")
    if err != nil {
        return fmt.Errorf("failed to load config: %v", err)
    }

    // Retrieve the database configuration values
    dbUser, err := web.AppConfig.String("dbuser")
    if err != nil {
        return fmt.Errorf("dbuser not found in config: %v", err)
    }
    
    dbPassword, err := web.AppConfig.String("dbpassword")
    if err != nil {
        return fmt.Errorf("dbpassword not found in config: %v", err)
    }

    dbHost, err := web.AppConfig.String("dbhost")
    if err != nil {
        return fmt.Errorf("dbhost not found in config: %v", err)
    }

    dbPort, err := web.AppConfig.String("dbport")
    if err != nil {
        return fmt.Errorf("dbport not found in config: %v", err)
    }

    dbName, err := web.AppConfig.String("dbname")
    if err != nil {
        return fmt.Errorf("dbname not found in config: %v", err)
    }

    // Initialize GORM
    gormDSN := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
        dbUser, dbPassword, dbName, dbHost, dbPort)
    
    db, err = gorm.Open(postgres.Open(gormDSN), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("failed to connect to database with GORM: %v", err)
    }

    // Initialize Beego ORM
    orm.RegisterDriver("postgres", orm.DRPostgres)
    
    beegoConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        dbUser, dbPassword, dbHost, dbPort, dbName)
    
    err = orm.RegisterDataBase("default", "postgres", beegoConnStr)
    if err != nil {
        return fmt.Errorf("failed to register database with Beego ORM: %v", err)
    }

    // Register models with Beego ORM
    // orm.RegisterModel(new(models.Location))
    orm.RegisterModel(
        new(models.Location),
        // new(models.PropertyDetails),
        // Add other models here
    )

    // Optional: Create tables if they don't exist
    err = orm.RunSyncdb("default", false, true)
    if err != nil {
        return fmt.Errorf("failed to sync database: %v", err)
    }

    // Use force=true to recreate tables (use with caution in production)
    err = orm.RunSyncdb("default", false, true)
    if err != nil {
        return fmt.Errorf("failed to sync database: %v", err)
    }

    // Set up connection pool parameters for Beego ORM
    orm.SetMaxIdleConns("default", 10)
    orm.SetMaxOpenConns("default", 100)

    log.Println("Database connected successfully (GORM and Beego ORM)")
    return nil
}

// GetDB returns the initialized GORM database instance
func GetDB() *gorm.DB {
    return db
}

// GetBeegoOrm returns a new Beego ORM object
func GetBeegoOrm() orm.Ormer {
    return orm.NewOrm()
}

// package utils

// import (
//     "log"
//     "github.com/beego/beego/v2/server/web"
//     "gorm.io/gorm"
//     "gorm.io/driver/postgres"
// )

// var db *gorm.DB


// // InitDatabaseFromConfig loads config and initializes the database connection
// func InitDatabaseFromConfig() error {
//     // Load configuration from app.conf
//     err := web.LoadAppConfig("ini", "conf/app.conf")
//     if err != nil {
//         return err
//     }

//     // Retrieve the database configuration values, handling errors
//     dbUser, err := web.AppConfig.String("dbuser")
//     if err != nil {
//         return err
//     }
    
//     dbPassword, err := web.AppConfig.String("dbpassword")
//     if err != nil {
//         return err
//     }

//     dbHost, err := web.AppConfig.String("dbhost")
//     if err != nil {
//         return err
//     }

//     dbPort, err := web.AppConfig.String("dbport")
//     if err != nil {
//         return err
//     }

//     dbName, err := web.AppConfig.String("dbname")
//     if err != nil {
//         return err
//     }

//     // Build the correct Data Source Name (DSN)
//     dsn := "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=disable"

//     // Initialize the database connection using the DSN
//     db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
//     if err != nil {
//         return err
//     }

//     log.Println("Database connected successfully")
//     return nil
// }

// // GetDB returns the initialized database instance
// func GetDB() *gorm.DB {
//     return db
// }
