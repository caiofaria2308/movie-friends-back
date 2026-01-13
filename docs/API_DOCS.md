# API Documentation

## Authentication provided by `accounts_router`

### Register
- **URL**: `/auth/register`
- **Method**: `POST`
- **Content-Type**: `application/json`
- **Body**:
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "role": "user" 
  }
  ```
- **Response**:
  - `200 OK`: `{"message": "User created successfully", "user_id": 1}`
  - `400 Bad Request`: Validation error
  - `500 Internal Server Error`: DB error

### Login
- **URL**: `/auth/login`
- **Method**: `POST`
- **Content-Type**: `application/json`
- **Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "password123"
  }
  ```
- **Response**:
  - `200 OK`: `{"token": "eyJhbG..."}`
  - `401 Unauthorized`: Invalid credentials

## Protected Routes
To access protected routes, include the token in the `Authorization` header:
`Authorization: Bearer <token>`

### User Profile
- **URL**: `/api/user/profile`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**:
  - `200 OK`: `{"name": "John Doe", "role": "user"}`
  - `401 Unauthorized`: Invalid token
  - `404 Not Found`: User not found

### Admin Dashboard
- **URL**: `/api/admin/dashboard`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>` (User must have `admin` role)
- **Response**: `{"message": "Welcome Admin"}`
