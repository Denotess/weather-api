# weather-api

Simple Gin API that fetches weather data from Visual Crossing and returns the
current temperature and conditions. Responses are cached in Redis by normalized
location.

## Requirements
- Go
- Redis (for caching)

## Configuration
Create `cmd/api/.env`:
```
URL=https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/{LOCATION}?unitGroup=metric&key={KEY}&contentType=json
KEY=your_api_key
ADDRESS=localhost:6379
PASSWORD=
DATABASE=0
```

`URL` must include `{LOCATION}` and `{KEY}`. Location is supplied per request.

## Run
```
go run ./cmd/api
```

## Usage
```
curl -X POST http://localhost:8080/weather \
  -H "Content-Type: application/json" \
  -d '{"location":"New York, NY"}'
```

Response:
```
{"temp":12.3,"conditions":"Partially cloudy"}
```

## Cache behavior
Cache key format is `weather:<normalized-location>` with a 30-minute TTL. If
Redis is unavailable, the API will still call the upstream service and return
fresh data.
