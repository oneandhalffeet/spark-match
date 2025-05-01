# Matchmaking Engine

## Project Description

a high-performance in-memory matchmaking engine for a dating app. The system should allow users to register profiles and retrieve the top 5 most compatible matches in real time.A combination of geohashing, scoring logic, and smart memory structures to make the system fast and extensible, even without a real database.

## Architecture Decisions

### Data Structure Design
- **In-memory storage**: Uses maps for O(1) profile lookups
- **Geohash quadrants**: Groups users by location for efficient proximity filtering  
- **Pre-computed scores**: Match scores calculated on registration for instant retrieval
- **Exclusion sets**: Maintains excluded users (matched/blocked/disliked) separately

## Running the Project

### Prerequisites
- Go 1.21 or higher
- Docker (optional)

### Local Development
```bash
# Install dependencies
go mod download

# Run the server
go run main.go

# Run tests
go test ./...
```

### Docker
```bash
# Build image
docker build -t matchmaking-engine .

# Run container
docker run -p 8080:8080 matchmaking-engine
```

## API Documentation

### POST /profiles
Register a new user profile.

Request:
```json
{
  "id": "user123",
  "age": 28,
  "gender": "F",
  "location": { "lat": 13.7563, "lon": 100.5018 },
  "interests": ["music", "art", "travel"]
}
```

### GET /match/:id
Get top 5 matches for a user.

Response:
```json
[
  {
    "profile_id": "user456",
    "score": 85.5
  },
  ...
]
```

### GET /seed?count=100
Generate dummy data for testing.

Response:
```json
{
  "message": "Successfully created 100 profiles",
  "count": 100
}
```

## Production Considerations

## Testing

The project includes unit tests for:
- Geohash calculation and quadrant assignment
- Match scoring algorithm
- Exclusion filtering
- Profile registration and retrieval
- Distance calculation

Run tests with: `go test ./...`