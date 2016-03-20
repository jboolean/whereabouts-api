# Whereabouts

An api to track and view user's locations relative to a home base.
Actual locations are hidden and presented as an an extimated travel time to home and current velocity.
Used by my roommate and I to see when the other will be home. 
Can scale to many users with varying permissions.

Interactive documentation is available via a built-in Swagger instance at /docs.

## Installation

Simply install and run as a Go app. The following environment variables must be present.

### Environment variables
- `PORT`: Web port to serve on
- `DATABASE_URL`: Postgres connection url
- `PGSSL`: Postgres sslmode (disable, require)
- `WA_MAPS_API_KEY`: A Google API key with access to the Google Maps Distance Matrix API.  
- `WA_HOME`: Coordinates of the home base e.g. "40.000000,-70.000000" 
- `WA_TRAVEL_MODE`: The transit mode to use for travel time. One of "walking", "driving", "transit", "bicycling"