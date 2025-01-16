# Property Listing and Details Design using Beego Framework

## Project Overview
This project aims to design two pages, **Property Listing** and **Property Details**, using the Beego framework. These pages will allow users to search for properties by location and view detailed information about a specific property.

## Pages Overview

### 1. Property Listing Page
The Property Listing page will display a list of properties available in a specific location. The page will include the following sections:
- **Search Box**: For users to search properties by location.
- **Location Breadcrumb**: To display the location hierarchy for easy navigation.
- **Tiles Section**: This will display a list of properties as tiles, and the content will be loaded dynamically using AJAX without requiring page reloads.

**Notes**: 
- The tiles section should be populated via AJAX call to ensure smooth user interaction and prevent page reloads.

### 2. Property Details Page
The Property Details page will display detailed information about a selected property, including:
- **Property Description**
- **Property Images**
- **Property Features**
- **Location**


**Notes**:
- The Property Details page will be accessible when a user clicks on a property name or property tile from the Property Listing page.
- Ensure that all sections of the property details are visible in the design.

### 3. Home or Index Page
The Home page will serve as the entry point for users to select a destination and view all properties available in that location. Once a destination is selected, users can:
1. See a list of all properties in that location (Property Listing page).
2. Click on a property to be redirected to the Property Details page.

## Technologies Used
- **Beego Framework**: For back-end logic and routing.
- **AJAX**: For loading property tiles dynamically without page reloads.
- **HTML/CSS**: For front-end design and layout.
- **JavaScript**: For handling AJAX calls and dynamic content loading.

## Features

### Property Listing Page Features:
- **Location Search Box**: Allows users to search for properties in a specific location.
- **Location Breadcrumbs**: Displays the current location and helps users navigate back.
- **Tiles Section**: Displays properties as tiles that can be clicked to view more details.
- **AJAX Integration**: The tiles are dynamically loaded using AJAX to prevent page reloads and enhance the user experience.

### Property Details Page Features:
- **Detailed Property Information**: Displays all necessary property details (description, images, pricing, etc.).
- **Clickable Property Name**: Users can click on the property tile or name from the Property Listing page to view more details.
- **Responsive Design**: Ensures the pages are responsive and accessible on different devices.


## How to Run

1. Clone the repository:
git clone 
cd 


2. Install dependencies:
- Navigate to the project directory:
  ```
  cd rental_view
  ```
- Install Beego framework and other dependencies (if any):
  ```
  go get -u github.com/beego/beego/v2
  ```

3. Run the application:
```
  go run main.go
  ```
  or,
  ```
   be run
  ```

4. Open your browser and navigate to:
http://localhost:8090

