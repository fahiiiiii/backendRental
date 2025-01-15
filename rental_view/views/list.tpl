<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vacation Rentals & Homes</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
            font-family: Arial, sans-serif;
        }

        body {
            background-color: #f5f5f5;
            padding: 20px;
        }

        .header {
            background-color: #003580;
            padding: 15px 20px;
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            z-index: 100;
        }

        .search-container {
            display: flex;
            gap: 10px;
            max-width: 1200px;
            margin: 0 auto;
        }

        .search-input {
            flex: 1;
            padding: 8px 15px;
            border: none;
            border-radius: 4px;
        }

        .date-input {
            width: 150px;
            padding: 8px 15px;
            border: none;
            border-radius: 4px;
        }

        .guests-input {
            width: 100px;
            padding: 8px 15px;
            border: none;
            border-radius: 4px;
        }

        .search-button {
            padding: 8px 20px;
            background-color: #00a699;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }

        .listings-container {
            max-width: 1200px;
            margin: 80px auto 20px;
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
            gap: 25px;
            padding: 20px;
        }

        .listing-card {
            background: white;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            transition: transform 0.2s;
            position: relative;
            aspect-ratio: 1/1.2;
            display: flex;
            flex-direction: column;
        }

        .listing-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        }

        .image-container {
            position: relative;
            width: 100%;
            height: 60%;
            overflow: hidden;
        }

        .listing-image {
            width: 100%;
            height: 100%;
            object-fit: cover;
            transition: transform 0.3s ease-in-out;
        }

        .image-container:hover .listing-image {
            transform: scale(1.1);
        }
          .favorite-button {
            position: absolute;
            top: 20px;
            right: 20px;
            font-size: 24px;
            cursor: pointer;
            transition: all 0.3s ease;
            color: #ccc;
        }

        .favorite-button:hover {
            transform: scale(1.1);
        }

        .favorite-button.active {
            color: #ff4d4d;
            animation: pulse 0.3s ease;
        }
        .favorite-button {
            position: absolute;
            top: 10px;
            right: 10px;
            background: white;
            border: none;
            border-radius: 50%;
            width: 32px;
            height: 32px;
            display: flex;
            align-items: center;
            justify-content: center;
            cursor: pointer;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }

        .favorite-button:hover {
            background: #f8f8f8;
        }

        .listing-details {
            padding: 15px;
            flex: 1;
            display: flex;
            flex-direction: column;
        }

        .listing-price {
            font-size: 1.2em;
            font-weight: bold;
            color: #333;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .price-value {
            color: #00a699;
        }

        .listing-title {
            margin: 8px 0;
            font-size: 1em;
            color: #484848;
            line-height: 1.4;
        }

        .listing-amenities {
            display: flex;
            gap: 10px;
            margin: 8px 0;
            color: #717171;
            font-size: 0.9em;
        }

        .rating {
            display: flex;
            align-items: center;
            gap: 5px;
            margin-top: auto;
            padding-top: 10px;
            border-top: 1px solid #eee;
        }

        .rating-star {
            color: #00a699;
        }

        .review-count {
            color: #717171;
            font-size: 0.9em;
        }

        .view-button {
            display: block;
            width: 100%;
            padding: 12px;
            background-color: #00a699;
            color: white;
            text-align: center;
            text-decoration: none;
            border-radius: 6px;
            margin-top: 15px;
            font-weight: bold;
            transition: background-color 0.2s;
        }

        .view-button:hover {
            background-color: #008578;
        }

        .price-badge {
            position: absolute;
            top: 10px;
            left: 10px;
            background: rgba(0,0,0,0.7);
            color: white;
            padding: 5px 10px;
            border-radius: 4px;
            font-size: 0.9em;
        }
        .breadcrumb {
    padding: 10px 0;
    color: #666;
    font-size: 14px;
    margin-bottom: 20px;
}

.breadcrumb span:last-child {
    color: #00a699;
    font-weight: 500;
}

.error {
    color: #ff4d4d;
    text-align: center;
    padding: 20px;
    background: #fff;
    border-radius: 8px;
    margin: 20px 0;
}

        @media (max-width: 768px) {
            .search-container {
                flex-direction: column;
            }
            
            .date-input, .guests-input {
                width: 100%;
            }

            .listings-container {
                grid-template-columns: 1fr;
                padding: 10px;
            }

            .listing-card {
                aspect-ratio: 1/1.4;
            }
        }
    </style>
    <script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
</head>

<body>
    <header class="header">
        <div class="search-container">
            <input type="text" class="search-input" placeholder="Dubai, Dubai, United Arab Emirates">
            <input type="date" class="date-input">
            <input type="number" class="guests-input" placeholder="Guests">
            <button class="search-button">Search</button>
        </div>
    </header>

     <div id="listings-container" class="listings-container">
        <!-- Listings will be dynamically loaded here -->
    </div>

    <script>
        $(document).ready(function() {
            // AJAX call to fetch properties
            $.ajax({
                url: '/v1/property/list',
                method: 'GET',
                success: function(properties) {
                    var listingsContainer = $('#listings-container');
                    listingsContainer.empty(); // Clear existing listings

                    properties.forEach(function(property) {
                        var amenitiesHtml = property.amenities.join(' • ');
                        var listingCard = `
                            <div class="listing-card">
                                <div class="image-container">
                                    <img src="${property.image}" alt="${property.title}" class="listing-image">
                                    <div class="price-badge">${property.badgeText}</div>
                                    <button class="favorite-button">♡</button>
                                </div>
                                <div class="listing-details">
                                    <div class="listing-price">
                                        <span class="price-value">From $${property.price.toLocaleString()}</span>
                                        <span class="rating-star">★ ${property.rating}</span>
                                    </div>
                                    <h3 class="listing-title">${property.title}</h3>
                                    <div class="listing-amenities">
                                        <span>${property.bedrooms} Bedrooms</span>
                                        <span>•</span>
                                        <span>${amenitiesHtml}</span>
                                    </div>
                                    <div class="rating">
                                        <span class="review-count">${property.reviewCount} reviews</span>
                                    </div>
                                    <a href="#" class="view-button">View Availability</a>
                                </div>
                                 <a href="/property/details?id=${property.id}" class="view-button">View Details</a
                            </div>
                        `;
                        listingsContainer.append(listingCard);
                    });
                },
                error: function(xhr, status, error) {
                    console.error("Error fetching properties:", error);
                }
            });

            // Optional: Search functionality
            $('.search-button').on('click', function() {
                // Implement search logic here
                console.log('Search clicked');
            });
        });








        document.addEventListener('DOMContentLoaded', function() {
    // Dropdown item click handler
    $(document).on('click', '.dropdown-item', function(e) {
        e.preventDefault();
        const cityName = $(this).text();
        const cityId = $(this).data('city-id');
        
        // Update search input with selected city
        $('.search-input').val(cityName);
        
        // Close dropdown
        $('#dropdown').hide();
        
        // Update breadcrumb
        updateBreadcrumb(cityName);
        
        // Load properties for selected city
        loadPropertiesForCity(cityId);
    });

    function updateBreadcrumb(cityName) {
        const breadcrumb = `
            <div class="breadcrumb mb-4">
                <span>Home</span> > 
                <span>Properties</span> > 
                <span>${cityName}</span>
            </div>
        `;
        $('#listings-container').before(breadcrumb);
    }

    function loadPropertiesForCity(cityId) {
        $.ajax({
            url: '/all-properties',
            method: 'GET',
            data: { city_id: cityId },
            success: function(response) {
                const listingsContainer = $('#listings-container');
                listingsContainer.empty();

                response.forEach(function(property) {
                    const listingCard = `
                        <div class="listing-card">
                            <div class="image-container">
                                <img src="${property.image}" alt="${property.title}" class="listing-image">
                                <div class="price-badge">${property.badgeText}</div>
                                <button class="favorite-button">♡</button>
                            </div>
                            <div class="listing-details">
                                <div class="listing-price">
                                    <span class="price-value">From $${property.price.toLocaleString()}</span>
                                    <span class="rating-star">★ ${property.rating}</span>
                                </div>
                                <h3 class="listing-title">${property.title}</h3>
                                <div class="listing-amenities">
                                    <span>${property.bedrooms} Bedrooms</span>
                                    <span>•</span>
                                    <span>${property.amenities.join(' • ')}</span>
                                </div>
                                <div class="rating">
                                    <span class="review-count">${property.reviewCount} reviews</span>
                                </div>
                                <a href="/property/details?id=${property.id}" class="view-button">View Details</a>
                            </div>
                        </div>
                    `;
                    listingsContainer.append(listingCard);
                });
            },
            error: function(xhr, status, error) {
                console.error("Error loading properties:", error);
                $('#listings-container').html('<p class="error">Error loading properties. Please try again.</p>');
            }
        });
    }
});
    </script>
</body>
</html>



