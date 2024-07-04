# Time Tracker API

## Description

Time Tracker API is designed to manage users and their tasks. The API provides capabilities for creating, updating, and deleting users, as well as managing user tasks.

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/your-repo/time-tracker.git
    ```

2. Navigate to the project directory:
    ```sh
    cd time-tracker
    ```

3. Install dependencies:
    ```sh
    go mod download
    ```

4. Set up your database and apply migrations.

5. Start the server:
    ```sh
    go run cmd/tracker/main.go
    ```

## Usage

### Endpoints

#### Get Users
- **URL:** `/users`
- **Method:** `GET`
- **Query Parameters:**
  - `passport_number` (string): Passport number.
  - `pass_serie` (string): Passport series.
  - `surname` (string): Surname.
  - `name` (string): Name.
  - `patronymic` (string): Patronymic.
  - `address` (string): Address.
  - `page` (int, required): Page number.
  - `page_size` (int, required): Number of items per page.
- **Responses:**
  - `200 OK`: List of users.
  - `400 Bad Request`: Invalid query parameters.
  - `404 Not Found`: No users found.
  - `500 Internal Server Error`: Server error.

#### Get User Tasks
- **URL:** `/users/tasks`
- **Method:** `GET`
- **Query Parameters:**
  - `user_id` (int, required): User ID.
  - `start_date` (string, required): Start date.
  - `end_date` (string, required): End date.
- **Responses:**
  - `200 OK`: List of user tasks.
  - `400 Bad Request`: Invalid query parameters.
  - `500 Internal Server Error`: Server error.

#### Create User
- **URL:** `/create`
- **Method:** `POST`
- **Request Body:**
  ```json
  {
    "passport_number": "string"
  }
  ```
- **Responses:**
  - `200 OK`: User created.
  - `400 Bad Request`: Invalid query parameters.
  - `500 Internal Server Error`: Server error.

#### Delete User
- **URL:** `/users/:id`
- **Method:** `DELETE`
- **Path Parameters:**
  - `id` (int): User ID.
- **Responses:**
  - `200 OK`: User deleted.
  - `400 Bad Request`: Invalid parameters.
  - `404 Not Found`: User not found.
  - `500 Internal Server Error`: Server error.

#### Update User
- **URL:** `/users/:id`
- **Method:** `PUT`
- **Path Parameters:**
  - `id` (int): User ID.
- **Request Body:**
  ```json
  {
    "surname": "string",
    "name": "string",
    "patronymic": "string",
    "address": "string"
  }
  ```
- **Responses:**
  - `200 OK`: User updated.
  - `400 Bad Request`: Invalid parameters.
  - `404 Not Found`: User not found.
  - `500 Internal Server Error`: Server error.

#### Start Task
- **URL:** `/tasks/start`
- **Method:** `POST`
- **Request Body:**
  ```json
  {
    "user_id": "int",
    "description": "string"
  }
  ```
- **Responses:**
  - `200 OK`: Task started.
  - `400 Bad Request`: Invalid request body.
  - `500 Internal Server Error`: Server error.
 
#### Stop Task
- **URL:** `/tasks/:id/stop`
- **Method:** `POST`
- **Path Parameters:**
  - `id` (int): Task ID.
- **Responses:**
  - `200 OK`: Task stopped.
  - `400 Bad Request`: Invalid task ID.
  - `500 Internal Server Error`: Server error.
 
## Swagger Specification

The Swagger specification:
 - [Swagger Specification English](swagger-spec-eng.yaml)
 - [Swagger Specification Russian](swagger-spec-ru.yaml)

To view and interact with the API documentation, you can use [Online Swagger UI](https://editor.swagger.io/). Copy code from file and paste to Editor.
