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

## User Pix

### Create Pix Key
- **URL**: `/api/user/pix`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "pix_key": "user@example.com"
  }
  ```
- **Response**:
  - `201 Created`: UserPix object
  - `500 Internal Server Error`: DB error

### List Pix Keys
- **URL**: `/api/user/pix`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**:
  - `200 OK`: List of UserPix objects
  - `500 Internal Server Error`: DB error

### Get Pix Key
- **URL**: `/api/user/pix/:id`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**:
  - `200 OK`: UserPix object
  - `404 Not Found`: Pix key not found or not owned by user

### Delete Pix Key
- **URL**: `/api/user/pix/:id`
- **Method**: `DELETE`
- **Headers**: `Authorization: Bearer <token>`
- **Response**:
  - `200 OK`: `{"message": "Pix key deleted successfully"}`
  - `500 Internal Server Error`: DB error

## User Day Off

### Create Day Off
- **URL**: `/api/user/dayoff`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "init_hour": "2023-10-27T08:00:00Z",
    "end_hour": "2023-10-27T17:00:00Z",
    "repeat": true,
    "repeat_type": "weekly",
    "repeat_value": "10"
  }
  ```
- **Response**:
  - `201 Created`: UserDayOff object (First instance)

### List Day Offs
- **URL**: `/api/user/dayoff`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**:
  - `200 OK`: List of UserDayOff objects

### Update Day Off
- **URL**: `/api/user/dayoff/:id?mode=[single|future|all]`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "init_hour": "2023-10-27T09:00:00Z",
    "end_hour": "2023-10-27T18:00:00Z"
  }
  ```
- **Response**:
  - `200 OK`: `{"message": "Day off updated successfully"}`

### Delete Day Off
- **URL**: `/api/user/dayoff/:id?mode=[single|future|all]`
- **Method**: `DELETE`
- **Headers**: `Authorization: Bearer <token>`
- **Response**:
  - `200 OK`: `{"message": "Day off deleted successfully"}`
