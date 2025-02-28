# Event Management API Documentation

## Base URL
`http://localhost:8080/api/v1`

## Authentication
The API uses JWT (JSON Web Token) for authentication. For protected endpoints, include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Error Handling
All error responses follow this format:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z",
  "code": 400,
  "message": "Error message description"
}
```

Common error codes:
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 409: Conflict/Duplicate Entry
- 413: Entity Too Large
- 500: Internal Server Error

## Rate Limiting
Protected endpoints have rate limiting. Headers in the response will include:
- `X-RateLimit-Limit`: Maximum number of requests allowed in the time window
- `X-RateLimit-Remaining`: Number of requests remaining in the current window
- `X-RateLimit-Reset`: Unix timestamp when the rate limit resets

---

# Authentication Endpoints

## Register User

Creates a new user account.

**URL**: `/auth/register`  
**Method**: `POST`  
**Auth Required**: No

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

**Success Response**:
- **Code**: 201 Created
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z"
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid input)
- **Code**: 409 Conflict (User already exists)

## Login

Authenticates a user and returns a JWT token.

**URL**: `/auth/login`  
**Method**: `POST`  
**Auth Required**: No

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid input)
- **Code**: 401 Unauthorized (Invalid credentials)

---

# User Endpoints

## Get User Profile

Retrieves the profile of the authenticated user.

**URL**: `/users/profile`  
**Method**: `GET`  
**Auth Required**: Yes

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z",
  "data": {
    "id": "uuid-string",
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": "user",
    "profile_image": "https://example.com/images/profile.jpg",
    "phone_number": "123-456-7890",
    "organization": "Company Name",
    "bio": "User biography text",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-02-01T00:00:00Z",
    "last_login_at": "2025-02-28T00:00:00Z"
  }
}
```

**Error Responses**:
- **Code**: 401 Unauthorized

## Update User Profile

Updates the profile information of the authenticated user.

**URL**: `/users/profile`  
**Method**: `PUT`  
**Auth Required**: Yes

**Request Body**:
```json
{
  "full_name": "John Updated Doe",
  "phone_number": "987-654-3210",
  "organization": "New Company",
  "bio": "Updated biography text"
}
```

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z"
}
```

**Error Responses**:
- **Code**: 400 Bad Request
- **Code**: 401 Unauthorized

## Change Password

Updates the password for the authenticated user.

**URL**: `/users/password`  
**Method**: `PUT`  
**Auth Required**: Yes

**Request Body**:
```json
{
  "current_password": "currentpassword123",
  "new_password": "newpassword456"
}
```

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z"
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid input or same password)
- **Code**: 401 Unauthorized
- **Code**: 403 Forbidden (Incorrect current password)

## Upload Profile Image

Uploads a profile image for the authenticated user.

**URL**: `/users/profile-image`  
**Method**: `PUT`  
**Auth Required**: Yes  
**Content-Type**: `multipart/form-data`

**Request Body**:
- `image`: The image file (JPEG, PNG, or GIF, max 10MB)

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z"
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid file)
- **Code**: 401 Unauthorized
- **Code**: 413 Request Entity Too Large (File too large)

---

# Category Endpoints

## List Categories

Retrieves all categories.

**URL**: `/categories`  
**Method**: `GET`  
**Auth Required**: No

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z",
  "data": [
    {
      "id": "uuid-string",
      "name": "Conference",
      "description": "Professional gathering for discussion",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    },
    {
      "id": "uuid-string",
      "name": "Workshop",
      "description": "Hands-on learning session",
      "created_at": "2025-01-02T00:00:00Z",
      "updated_at": "2025-01-02T00:00:00Z"
    }
  ]
}
```

## Get Category

Retrieves a specific category by ID.

**URL**: `/categories/{id}`  
**Method**: `GET`  
**Auth Required**: No

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z",
  "data": {
    "id": "uuid-string",
    "name": "Conference",
    "description": "Professional gathering for discussion",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

**Error Response**:
- **Code**: 404 Not Found

## Create Category

Creates a new category.

**URL**: `/categories`  
**Method**: `POST`  
**Auth Required**: Yes

**Request Body**:
```json
{
  "name": "Seminar",
  "description": "Educational presentation on a specific topic"
}
```

**Success Response**:
- **Code**: 201 Created
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z"
}
```

**Error Responses**:
- **Code**: 400 Bad Request
- **Code**: 401 Unauthorized
- **Code**: 409 Conflict (Duplicate category name)

## Update Category

Updates an existing category.

**URL**: `/categories/{id}`  
**Method**: `PUT`  
**Auth Required**: Yes

**Request Body**:
```json
{
  "name": "Updated Seminar",
  "description": "Updated description for the seminar category"
}
```

**Success Response**:
- **Code**: 201 Created
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z"
}
```

**Error Responses**:
- **Code**: 400 Bad Request
- **Code**: 401 Unauthorized
- **Code**: 404 Not Found

## Delete Category

Deletes a specific category.

**URL**: `/categories/{id}`  
**Method**: `DELETE`  
**Auth Required**: Yes

**Success Response**:
- **Code**: 204 No Content
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z"
}
```

**Error Responses**:
- **Code**: 401 Unauthorized
- **Code**: 404 Not Found

---

# Event Endpoints

## List Events

Retrieves a paginated list of events.

**URL**: `/events`  
**Method**: `GET`  
**Auth Required**: No

**Query Parameters**:
- `page`: Page number (default: 1)
- `page_size`: Number of items per page (default: 10, max: 100)
- `sort_by`: Field to sort by (options: title, start_date, end_date, created_at)
- `sort_dir`: Sort direction (asc or desc, default: desc)

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z",
  "data": [
    {
      "id": "uuid-string",
      "title": "Tech Conference 2025",
      "description": "Annual technology conference",
      "start_date": "2025-06-15T09:00:00Z",
      "end_date": "2025-06-17T18:00:00Z",
      "creator_id": "user-uuid-string",
      "category_id": "category-uuid-string",
      "tags": [
        {
          "id": "tag-uuid-string",
          "name": "technology"
        }
      ],
      "files": [
        {
          "id": "file-uuid-string",
          "event_id": "event-uuid-string",
          "file_name": "schedule.pdf",
          "file_type": "application/pdf",
          "file_url": "https://example.com/files/schedule.pdf",
          "created_at": "2025-01-15T00:00:00Z"
        }
      ],
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-10T00:00:00Z"
    }
  ]
}
```

## Search Events

Searches for events with various filters.

**URL**: `/events/search`  
**Method**: `GET`  
**Auth Required**: No

**Query Parameters**:
- `query`: Search term in title and description
- `start_date`: Filter events starting after this date (RFC3339)
- `end_date`: Filter events ending before this date (RFC3339)
- `creator`: Filter by creator ID
- `page`: Page number (default: 1)
- `page_size`: Number of items per page (default: 10, max: 100)
- `sort_by`: Field to sort by (options: title, start_date, end_date, created_at)
- `sort_dir`: Sort direction (asc or desc)

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z",
  "data": {
    "events": [
      {
        "id": "uuid-string",
        "title": "Tech Conference 2025",
        "description": "Annual technology conference",
        "start_date": "2025-06-15T09:00:00Z",
        "end_date": "2025-06-17T18:00:00Z",
        "creator_id": "user-uuid-string",
        "category_id": "category-uuid-string",
        "tags": [],
        "files": [],
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-10T00:00:00Z"
      }
    ],
    "total_count": 5,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

## Get Event

Retrieves a specific event by ID.

**URL**: `/events/{id}`  
**Method**: `GET`  
**Auth Required**: No

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z",
  "data": {
    "id": "uuid-string",
    "title": "Tech Conference 2025",
    "description": "Annual technology conference",
    "start_date": "2025-06-15T09:00:00Z",
    "end_date": "2025-06-17T18:00:00Z",
    "creator_id": "user-uuid-string",
    "category_id": "category-uuid-string",
    "tags": [
      {
        "id": "tag-uuid-string",
        "name": "technology"
      }
    ],
    "files": [
      {
        "id": "file-uuid-string",
        "event_id": "event-uuid-string",
        "file_name": "schedule.pdf",
        "file_type": "application/pdf",
        "file_url": "https://example.com/files/schedule.pdf",
        "created_at": "2025-01-15T00:00:00Z"
      }
    ],
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-10T00:00:00Z"
  }
}
```

**Error Response**:
- **Code**: 404 Not Found

## Create Event

Creates a new event.

**URL**: `/events`  
**Method**: `POST`  
**Auth Required**: Yes

**Request Body**:
```json
{
  "title": "Tech Workshop 2025",
  "description": "Hands-on technology workshop",
  "category_id": "category-uuid-string",
  "start_date": "2025-07-15T09:00:00Z",
  "end_date": "2025-07-15T17:00:00Z"
}
```

**Success Response**:
- **Code**: 201 Created
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z"
}
```

**Error Responses**:
- **Code**: 400 Bad Request
- **Code**: 401 Unauthorized
- **Code**: 404 Not Found (Category not found)

## Update Event

Updates an existing event.

**URL**: `/events/{id}`  
**Method**: `PUT`  
**Auth Required**: Yes

**Request Body**:
```json
{
  "title": "Updated Tech Workshop 2025",
  "description": "Updated workshop description",
  "start_date": "2025-07-16T09:00:00Z",
  "end_date": "2025-07-16T17:00:00Z"
}
```

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z"
}
```

**Error Responses**:
- **Code**: 400 Bad Request
- **Code**: 401 Unauthorized
- **Code**: 403 Forbidden (Not the event creator)
- **Code**: 404 Not Found

## Delete Event

Deletes a specific event.

**URL**: `/events/{id}`  
**Method**: `DELETE`  
**Auth Required**: Yes

**Success Response**:
- **Code**: 204 No Content
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z"
}
```

**Error Responses**:
- **Code**: 401 Unauthorized
- **Code**: 403 Forbidden (Not the event creator)
- **Code**: 404 Not Found

## Upload File to Event

Uploads a file attachment to an event.

**URL**: `/events/{id}/upload`  
**Method**: `POST`  
**Auth Required**: Yes  
**Content-Type**: `multipart/form-data`

**Request Body**:
- `file`: The file to upload (max 10MB)

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-02-28T12:34:56.789Z"
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid file)
- **Code**: 401 Unauthorized
- **Code**: 404 Not Found (Event not found)
- **Code**: 413 Request Entity Too Large (File too large)

---

# Health Endpoints

## Health Check

Check the overall health of the API.

**URL**: `/health`  
**Method**: `GET`  
**Auth Required**: No

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "status": "UP",
  "components": {
    "database": {
      "name": "database",
      "status": "UP",
      "details": {
        "open_connections": 5,
        "in_use": 1,
        "idle": 4,
        "max_open_connections": 100
      }
    },
    "redis": {
      "name": "redis",
      "status": "UP",
      "details": {
        "info": {
          "Used Memory": "1.25MB",
          "Used Memory RSS": "3.45MB",
          "Memory Fragmentation Ratio": "2.76",
          "Max Memory": "0B",
          "Max Memory Policy": "noeviction"
        }
      }
    },
    "memory": {
      "name": "memory",
      "status": "UP",
      "details": {
        "alloc": 8388608,
        "total_alloc": 16777216,
        "sys": 33554432,
        "num_gc": 10
      }
    },
    "disk": {
      "name": "disk",
      "status": "UP",
      "details": {
        "path": "."
      }
    }
  },
  "timestamp": "2025-02-28T12:34:56.789Z",
  "version": "1.0.0",
  "uptime": "24h0m0s"
}
```

**Error Response**:
- **Code**: 503 Service Unavailable (If any component is down)

## Liveness Check

Simple check to determine if the application is running.

**URL**: `/health/liveness`  
**Method**: `GET`  
**Auth Required**: No

**Success Response**:
- **Code**: 200 OK
- **Content**: `OK`

## Readiness Check

Check if the application is ready to handle requests.

**URL**: `/health/readiness`  
**Method**: `GET`  
**Auth Required**: No

**Success Response**:
- **Code**: 200 OK (All components are ready)
- **Content**: Same as health check

**Error Response**:
- **Code**: 503 Service Unavailable (If any component is not ready)
