# API Documentation

## Base URL
`http://localhost:8080`

## Endpoints

### Create Profile
**POST** `/profiles`

Register a new user profile in the system.

#### Request Body
```json
{
  "id": "string",
  "age": "integer",
  "gender": "string",
  "location": {
    "lat": "float",
    "lon": "float"
  },
  "interests": ["string"]
}
```

#### Response
- **200 OK**: Profile created successfully
```json
{
  "id": "user123",
  "age": 28,
  "gender": "F",
  "location": { "lat": 13.7563, "lon": 100.5018 },
  "interests": ["music", "art", "travel"]
}
```
- **400 Bad Request**: Invalid input data
- **500 Internal Server Error**: Server error

---

### Get Matches
**GET** `/match/{id}`

Retrieve top 5 compatible matches for a user.

#### Parameters
- `id` (path): User ID to get matches for

#### Response
- **200 OK**: Matches retrieved successfully
```json
[
  {
    "profile_id": "user456",
    "score": 85.5
  },
  {
    "profile_id": "user789",
    "score": 78.2
  }
]
```
- **404 Not Found**: User profile not found

---

### Seed Data
**GET** `/seed`

Generate dummy profiles for testing.

#### Query Parameters
- `count` (optional): Number of profiles to generate (default: 100)

#### Response
- **200 OK**: Data seeded successfully
```json
{
  "message": "Successfully created 100 profiles",
  "count": 100
}
```
- **500 Internal Server Error**: Error generating data

## Data Models

### Profile
| Field | Type | Description |
|-------|------|-------------|
| id | string | Unique user identifier |
| age | integer | User's age (18+) |
| gender | string | User's gender (M/F) |
| location | Location | User's coordinates |
| interests | string[] | List of user interests |

### Location
| Field | Type | Description |
|-------|------|-------------|
| lat |