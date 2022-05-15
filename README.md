# Capture Maps Backend API

This API captures detailed images of an area
using [Google Maps Static API](https://developers.google.com/maps/documentation/maps-static/overview)

## Setup

1. Create `.env` in the project root and store Google Maps Static API
   `API_KEY=<YOUR API_KEY>`
2. Run project
   `go run src/main.go`

## Sample request

```http request
POST /print

{
    "lat": 30.316963, 
    "lng": 78.032560,
    "zoom": 15,
    "radius": 5
}
```

## Sample Response

![sample response](./assets/response.png)

