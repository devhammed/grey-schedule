# Schedule Management System

A schedule management system with a React frontend and a Go backend.

- Backend: Go server with in-memory storage, HTTP API, and gRPC API
- Frontend: React + TypeScript app using React Query
- Concurrency handling and conflict prevention for bookings

## Structure

- `backend/` — Go backend (HTTP API, gRPC API, Storage)
- `frontend/` — React + TS frontend (Vite)

## Prerequisites

- Go 1.22+
- Node.js 18+ and npm

## Running the Backend

From the `backend` directory:

```shell 
go run ./cmd/server
```

This starts the gRPC server on port 8081 and the HTTP server on port 8080.

### Environment:

- `GRPC_ADDR`: gRPC address to bind (default :8081)
- `HTTP_ADDR`: HTTP address to bind (default :8080)

### GRPC API

**Proto syntax:** `proto3`
**Package:** `schedule`
**Go package:** `github.com/devhammed/grey-schedule/backend/proto;schedulepb`

---

#### Messages

##### `Appointment`
Represents a single appointment.

| Field      | Type   | Description          |
|------------|--------|--------------------|
| `id`       | string | Unique appointment ID |
| `title`    | string | Appointment title    |
| `start`    | string | Start time (RFC3339) |
| `end`      | string | End time (RFC3339)   |
| `createdAt`| string | Creation timestamp (RFC3339) |

---

##### `CreateAppointmentRequest`
Request message for creating an appointment.

| Field   | Type   | Description          |
|---------|--------|--------------------|
| `title` | string | Appointment title    |
| `start` | string | Start time (RFC3339) |
| `end`   | string | End time (RFC3339)   |

---

##### `CreateAppointmentResponse`
Response message for creating an appointment.

| Field       | Type        | Description       |
|-------------|------------|-----------------|
| `appointment` | `Appointment` | The created appointment |

---

##### `ListAppointmentsRequest`
Request message to list all appointments (empty message).

---

##### `ListAppointmentsResponse`
Response message for listing appointments.

| Field           | Type        | Description           |
|-----------------|------------|---------------------|
| `appointments`  | repeated `Appointment` | List of appointments |

---

##### `DeleteAppointmentRequest`
Request message to delete an appointment.

| Field | Type   | Description        |
|-------|--------|------------------|
| `id`  | string | ID of the appointment to delete |

---

##### `DeleteAppointmentResponse`
Response message for deleting an appointment (empty message).

---

#### Service: `AppointmentService`
Defines the gRPC service for managing appointments.

| RPC Method             | Request Type               | Response Type                | Description                     |
|------------------------|---------------------------|------------------------------|---------------------------------|
| `CreateAppointment`    | `CreateAppointmentRequest` | `CreateAppointmentResponse`  | Creates a new appointment       |
| `ListAppointments`     | `ListAppointmentsRequest`  | `ListAppointmentsResponse`   | Returns all appointments        |
| `DeleteAppointment`    | `DeleteAppointmentRequest` | `DeleteAppointmentResponse`  | Deletes an appointment by ID    |

---

#### Example grpcurl Commands

**Create an appointment**
```shell
grpcurl -plaintext -d '{"title":"Team Meeting","start":"2025-10-14T17:00:00Z","end":"2025-10-14T18:00:00Z"}' localhost:8081 schedule.AppointmentService/CreateAppointment
```

**List appointments**
```shell
grpcurl -plaintext -d '{}' localhost:8081 schedule.AppointmentService/ListAppointments
```

**Delete an appointment**
```shell
grpcurl -plaintext -d '{"id":"550e8400-e29b-41d4-a716-446655440000"}' localhost:8081 schedule.AppointmentService/DeleteAppointment
```


### HTTP API

All endpoints use JSON for request and response bodies.

#### Endpoints

##### Health Check
**GET** `/health`

Returns the health status of the server.

**Response:**
- **200 OK**: `ok`

---

##### Create Appointment
**POST** `/api/appointments`

Creates a new appointment. Validates the request and checks for time conflicts with existing appointments.

**Request Body:**
```json
{
  "title": "string",
  "start": "2025-10-14T17:45:30Z",
  "end": "2025-10-15T17:45:30Z"
}
```

**Request Fields:**

| Field   | Type   | Required | Description                                  |
|---------|--------|----------|----------------------------------------------|
| `title` | string | Yes      | Appointment title (1-255 characters)         |
| `start` | string | Yes      | Start time in RFC3339 format                 |
| `end`   | string | Yes      | End time in RFC3339 format (must be after start) |

**Response:**

- **201 Created**
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Team Meeting",
    "start": "2025-10-14T17:45:30Z",
    "end": "2025-10-15T17:45:30Z",
    "createdAt": "2025-10-14T16:00:00Z"
  }
  ```

- **400 Bad Request**
    - Invalid JSON body
    - Invalid RFC3339 timestamp format
    - Title is empty or exceeds 255 characters
    - Start time is not before end time

  ```json
  {
    "error": "title is required and must be 1-255 characters"
  }
  ```
  ```json
  {
    "error": "invalid time range: start must be before end"
  }
  ```
  ```json
  {
    "error": "start and end must be RFC3339 timestamps"
  }
  ```

- **409 Conflict**
    - Time conflict with existing appointment

  ```json
  {
    "error": "time conflict with existing appointment"
  }
  ```

**Example:**
```bash
curl -X POST http://localhost:8080/api/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Meeting",
    "start": "2025-10-14T17:00:00Z",
    "end": "2025-10-14T18:00:00Z"
  }'
```

---

##### List Appointments
**GET** `/api/appointments`

Returns all appointments sorted by start time (and creation time if start times are equal).

**Response:**

- **200 OK**
  ```json
  [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Team Meeting",
      "start": "2025-10-14T17:00:00Z",
      "end": "2025-10-14T18:00:00Z",
      "createdAt": "2025-10-14T16:00:00Z"
    },
    {
      "id": "660f9511-f2ac-52e5-b827-557766551111",
      "title": "Client Call",
      "start": "2025-10-15T10:00:00Z",
      "end": "2025-10-15T11:00:00Z",
      "createdAt": "2025-10-14T16:30:00Z"
    }
  ]
  ```

**Example:**
```bash
curl http://localhost:8080/api/appointments
```

---

##### Delete Appointment
**DELETE** `/api/appointments/{id}`

Deletes an appointment by ID.

**URL Parameters:**

| Parameter | Type   | Description                    |
|-----------|--------|--------------------------------|
| `id`      | string | UUID of the appointment to delete |

**Response:**

- **204 No Content**
    - Appointment successfully deleted

- **404 Not Found**
  ```json
  {
    "error": "appointment not found"
  }
  ```

- **500 Internal Server Error**
  ```json
  {
    "error": "internal error"
  }
  ```

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/appointments/550e8400-e29b-41d4-a716-446655440000
```

---

#### Appointment Object

| Field       | Type   | Description                           |
|-------------|--------|---------------------------------------|
| `id`        | string | UUID of the appointment               |
| `title`     | string | Appointment title (1-255 characters)  |
| `start`     | string | Start time in RFC3339 format          |
| `end`       | string | End time in RFC3339 format            |
| `createdAt` | string | Creation timestamp in RFC3339 format  |

### Testing Concurrency

A Go unit test demonstrates safe concurrent booking:

```shell
go test ./internal/store -run TestConcurrentCreateConflict -v
```

## Running the Frontend

From the `frontend` directory:

```shell
npm install
npm run dev
```

You should be able to access the app at http://localhost:5173.

## Environment:
- `VITE_API_BASE`: base URL for the HTTP API (default http://localhost:8080)

