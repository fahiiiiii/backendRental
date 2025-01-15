
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Property Details</title>
    <script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
        }

        body {
            padding: 20px;
            background-color: #f5f5f5;
        }

        .listing-container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            overflow: hidden;
        }

        .header {
            padding: 15px 20px;
            border-bottom: 1px solid #eee;
        }

        .property-title {
            font-size: 1.5rem;
            color: #2c3e50;
            margin-bottom: 10px;
        }

        .rating-info {
            display: flex;
            align-items: center;
            gap: 15px;
            margin-bottom: 10px;
        }

        .stars {
            color: #ffd700;
        }

        .rating-score {
            background: #003580;
            color: white;
            padding: 4px 8px;
            border-radius: 4px;
            font-weight: bold;
        }

        .property-info {
            display: flex;
            gap: 15px;
            color: #666;
            font-size: 0.9rem;
        }

        .gallery {
            display: grid;
            grid-template-columns: 2fr 1fr;
            gap: 4px;
            padding: 4px;
        }

        .main-image {
            grid-row: span 2;
            height: 400px;
            position: relative;
        }

        .side-images {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 4px;
        }

        .image {
            width: 100%;
            height: 198px;
            object-fit: cover;
            border-radius: 4px;
        }

        .main-image img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            border-radius: 4px;
        }

        .description {
            padding: 20px;
            color: #444;
            line-height: 1.6;
        }

        .view-more {
            position: absolute;
            bottom: 20px;
            left: 20px;
            background: rgba(255,255,255,0.9);
            padding: 8px 16px;
            border-radius: 4px;
            text-decoration: none;
            color: #003580;
            font-weight: 500;
        }

        @media (max-width: 768px) {
            .gallery {
                grid-template-columns: 1fr;
            }

            .main-image {
                height: 300px;
            }

            .side-images {
                grid-template-columns: 1fr;
            }

            .image {
                height: 200px;
            }

            .property-info {
                flex-direction: column;
                gap: 5px;
            }
        }
    </style>
</head>
<body>
    <div id="property-details-container" class="listing-container">
        <!-- Details will be dynamically loaded here -->
    </div>

    <script>
        $(document).ready(function() {
            // Extract property ID from URL (assuming format like /property/details?id=123)
            const urlParams = new URLSearchParams(window.location.search);
            const propertyId = urlParams.get('id');

            // AJAX call to fetch property details
            $.ajax({
                url: '/v1/property/details',
                method: 'GET',
                data: { id: propertyId },
                success: function(property) {
                    const container = $('#property-details-container');
                    
                    // Render header
                    const headerHtml = `
                        <div class="header">
                            <h1 class="property-title">${property.title}</h1>
                            <div class="rating-info">
                                <div class="stars">★★★★★</div>
                                <div class="rating-score">${property.rating.toFixed(1)}</div>
                                <div>(${property.reviewCount} Reviews)</div>
                            </div>
                            <div class="property-info">
                                <span>${property.bedrooms} Bedroom</span>
                                <span>${property.bathrooms} Bathroom</span>
                                <span>${property.guests} Guests</span>
                            </div>
                        </div>
                    `;

                    // Render gallery
                    const mainImage = property.images[0];
                    const sideImages = property.images.slice(1, 5);
                    const galleryHtml = `
                        <div class="gallery">
                            <div class="main-image">
                                <img src="${mainImage}" alt="Main Property Image">
                                <a href="#" class="view-more">View More Photos</a>
                            </div>
                            <div class="side-images">
                                ${sideImages.map(img => `
                                    <img class="image" src="${img}" alt="Property Image">
                                `).join('')}
                            </div>
                        </div>
                    `;

                    // Render description
                    const descriptionHtml = `
                        <div class="description">
                            <h2>${property.bedrooms} Bedroom Property in Dubai</h2>
                            <p>${property.description}</p>
                        </div>
                    `;

                    // Combine all sections
                    container.html(headerHtml + galleryHtml + descriptionHtml);
                },
                error: function(xhr, status, error) {
                    console.error("Error fetching property details:", error);
                    $('#property-details-container').html(`
                        <div class="error">
                            Unable to load property details. Please try again later.
                        </div>
                    `);
                }
            });

            // Optional: View More Photos functionality
            $(document).on('click', '.view-more', function(e) {
                e.preventDefault();
                // Implement photo gallery modal or lightbox
                alert('Full photo gallery to be implemented');
            });
        });
    </script>
</body>
</html>

