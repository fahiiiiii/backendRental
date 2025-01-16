# Property Listing Project

This project is a full-stack web application for listing rental properties in the United Arab Emirates. 

## Setup the Project

1. **Clone the repository**:
   ```sh
   git clone https://github.com/fahiiiiii/backend_rental
   cd backend_rental


***Running the Project***
Run the project:

Recommended:
    ```sh
    go run main.go

Alternative:
    ```sh
    bee run

This will start the server, fetch data, create tables, and store data into the database.

***Viewing the Database***
Access the Docker container:

    ```sh
    docker exec -it property-listing-db-1 /bin/bash
    
Connect to the PostgreSQL database:

    ```sh
    psql -U fahimah -d rental_db

Check the tables:
 SQL
\dt

Run sample queries:

SQL
SELECT PropertyName, City, Country FROM locations LIMIT 10;
-- or,
SELECT * FROM "location" LIMIT 10;
SELECT * FROM RentalProperty;
SELECT * FROM property_details;


Dependencies
Ensure you have the following dependencies installed:

Go
Bee
Docker
