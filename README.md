# Go Profile CRUD API

A small REST API built with Go's standard library. This project is a hands-on exercise in HTTP routing, JSON request/response handling, CRUD operations, validation, generated IDs, and protecting shared in-memory state with a mutex.

## Requirements

- Go 1.22 or newer

No third-party dependencies are required.

## Run locally

```bash
git clone https://github.com/Abhyshekbhalaji/go-profile-api.git
cd go-profile-api
go run .
```

The server listens on `http://localhost:8080`.

## Endpoints

| Method | Path | Description | Success status |
| --- | --- | --- | --- |
| GET | `/users/{userId}` | Fetch one profile | `200 OK` |
| POST | `/users` | Create a profile with a generated ID | `201 Created` |
| PUT | `/users/{userId}` | Replace an existing profile | `200 OK` |
| DELETE | `/users/{userId}` | Delete a profile | `204 No Content` |

### Profile JSON

```json
{
  "name": "Alex",
  "username": "alex_dev",
  "posts": 2,
  "reels": 0
}
```

### Try it with curl

Fetch the seeded profile:

```bash
curl -i http://localhost:8080/users/1
```

Create a profile:

```bash
curl -i -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alex","username":"alex_dev","posts":2,"reels":0}'
```

The create response uses `201 Created` and includes a `Location` header for the newly created resource. IDs are generated sequentially while the server is running.

Replace an existing profile:

```bash
curl -i -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alex","username":"alex_dev","posts":3,"reels":1}'
```

Delete a profile:

```bash
curl -i -X DELETE http://localhost:8080/users/1
```

## Go concepts practiced

### 1. Method-aware HTTP routing

`http.NewServeMux()` registers routes such as `GET /users/{userId}`. Go 1.22+ supports HTTP method patterns and path wildcards. Inside a handler, `r.PathValue("userId")` reads the wildcard value.

### 2. Handlers and HTTP responses

A handler receives an `http.ResponseWriter` and an `*http.Request`. The writer sets response headers, writes status codes, and sends response data. `http.Error` sends an error response and status code.

### 3. JSON encoding and decoding

- Struct tags such as `json:"username"` define JSON field names.
- `json.NewDecoder(r.Body).Decode(&profile)` reads a JSON request into a struct. The pointer lets the decoder populate the struct.
- `decoder.DisallowUnknownFields()` rejects JSON fields that do not match the struct.
- `json.NewEncoder(w).Encode(profile)` serializes a value as JSON in the response.
- Unknown-field checks and decoding do not replace business validation; the API separately checks required strings and non-negative counters.

### 4. Map lookups and generated IDs

A map lookup can return a value and an existence flag: `profile, exists := usersSampleDb[userID]`. The comma-ok result distinguishes a missing ID from a stored value. `strconv.Itoa(nextID)` converts the integer counter into a string map key.

### 5. Shared state and locking

The map is shared by concurrent HTTP handlers, so reads and writes need synchronization.

- `sync.RWMutex` allows multiple readers or one writer.
- `RLock()` / `RUnlock()` protect map reads.
- `Lock()` / `Unlock()` protect writes and the ID counter.
- The PUT handler checks whether the ID exists and updates the map while holding the same write lock, so another request cannot change the map between those operations.
- Avoid attempting to acquire the same non-reentrant mutex again while already holding it; that can deadlock.

### 6. CRUD semantics and status codes

- **Create:** POST validates the JSON, stores a profile, and returns `201 Created`.
- **Read:** GET returns `200 OK`, or `404 Not Found` if the ID is absent.
- **Update:** PUT replaces the full profile. It is not a partial update.
- **Delete:** DELETE returns `204 No Content`, or `404 Not Found` if the ID is absent.
- Invalid JSON or invalid profile values produce `400 Bad Request`.

## Current limitations and next steps

- Data is stored in memory, so profiles disappear when the process stops.
- PATCH (partial updates), persistent storage, authentication, and pagination are not implemented.
- Automated endpoint tests with `net/http/httptest` are a useful next step. Tests can create an HTTP request, pass it to the handler or server, and inspect the resulting status, headers, and body without opening a real network port.
- Further hardening could reject trailing JSON values, limit request-body size, and improve error logging.

## Project goal

This is a learning project focused on understanding the mechanics behind a small Go HTTP API before adding production concerns.
