# Matchmaking Engine

A high-performance in-memory matchmaking engine for dating applications, built with Go.

## Project Description

This system implements a real-time matchmaking service that uses geohashing for efficient location-based filtering and pre-computed match scores for instant retrieval. Users can register profiles and get their top 5 most compatible matches based on age similarity, shared interests, and location proximity.

## Architecture Decisions

### Data Structure Design
- **In-memory storage**: Uses maps for O(1) profile lookups
- **Geohash quadrants**: Groups users by location for efficient proximity filtering  
- **Pre-computed scores**: Match scores calculated on registration for instant retrieval
- **Exclusion sets**: Maintains excluded users (matched/blocked/disliked) separately
- **Gender preference filtering**: Supports user-defined gender preferences for matching
- **Pagination support**: Allows fetching matches in pages for better scalability

### Pre-computation Strategy
When a new profile registers:
1. Calculate geohash (precision 5) for location-based indexing
2. Score against all profiles in same and adjacent quadrants
3. Store sorted match lists for both new and existing profiles
4. Result: O(1) match retrieval, O(n) registration where n = users in quadrant

### Scoring Algorithm
- Age similarity: 30 points max (reduced by 2 points per year difference)
- Shared interests: 40 points max (10 points per common interest)
- Location proximity: 30 points max (reduced by 1 point per 5km)
- Total: 100 points max

## Running the Project

### Prerequisites
- Go 1.21 or higher
- Docker (optional)

### Quick Start
```bash
# Install dependencies
make deps

# Run the server
make run

# Run tests
make test

# Seed sample data (server must be running)
make seed
```

### Docker
```bash
# Build image
make docker-build

# Run container
make docker-run
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
  "interests": ["music", "art", "travel"],
  "looking_for": "M"
}
```

### PUT /profiles
Update an existing user profile.

Request: Same as POST /profiles

### GET /match/:id
Get matches for a user with pagination.

Query Parameters:
- `limit`: Number of results (default: 5)
- `offset`: Pagination offset (default: 0)

Response:
```json
{
  "matches": [
    {
      "profile_id": "user456",
      "score": 85.5
    }
  ],
  "total": 15,
  "limit": 5,
  "offset": 0
}
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

### What I Would Do Differently
1. **Database Persistence**: Add PostgreSQL/MongoDB for data durability
2. **Caching Layer**: Implement Redis for distributed caching
3. **Horizontal Scaling**: Shard by geohash prefix for multiple instances
4. **Async Processing**: Use message queues for match computation
5. **Monitoring**: Add Prometheus metrics and health checks
6. **Authentication**: Implement JWT-based auth
7. **Rate Limiting**: Protect endpoints from abuse
8. **Batch Operations**: Optimize bulk profile updates

### Performance Optimizations
- Use goroutines for parallel match computation
- Implement connection pooling
- Add request timeouts and circuit breakers
- Use protobuf for internal communication

## Testing

The project includes unit tests for:
- Geohash calculation and quadrant assignment
- Match scoring algorithm
- Exclusion filtering
- Profile registration and retrieval
- Distance calculation

Run tests with: `go test ./...`