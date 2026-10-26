# HTTP testing guide

This repository's HTTP tests live in `task_test.go`. The filename ends in
`_test.go`, so Go's test runner discovers the test functions in it.

## What is being tested?

`TestGetUserE2E` tests the GET user endpoint through the application's real
`http.ServeMux` and handler:

- `GET /users/1` should return `200 OK` and a body containing `Abhyshek`.
- `GET /users/3` should return `404 Not Found` and a body containing
  `No such profile exists`.

The tests use the initial in-memory sample data from `main.go`. They do not
need the application to be running separately.

## How the test works

1. **`newMux()`** builds the same router used by `main()`. This ensures the
   route patterns and handlers are exercised together.
2. **`httptest.NewServer(newMux())`** starts a temporary local HTTP server for
   the duration of the test.
3. **`server.Client()`** provides an HTTP client configured for that test server.
4. **The test table** is a slice of structs. Each row stores a case name, URL
   path, expected status code, and expected text in the response body.
5. **`for _, tt := range tests`** loops over the cases. `tt` refers to the
   current case.
6. **`t.Run(tt.name, ...)`** runs a named subtest for each row. A failure
   identifies the scenario that failed.
7. **Assertions** compare the actual status and body with the expected values.

This is an in-process integration test over HTTP. It tests the router and
handler together, but does not test deployment, TLS, or an external database.

## Run all tests

From the repository root:

```bash
go test -v
```

## Run only the GET test

```bash
go test -run '^TestGetUserE2E$' -v
```

## Run only one subtest

Run only the existing-user case:

```bash
go test -run '^TestGetUserE2E/existing_user$' -v
```

Run only the missing-user case:

```bash
go test -run '^TestGetUserE2E/missing_user$' -v
```

Go normalizes spaces in subtest names to underscores in the reported test name.

## Notes

- `server.URL` is the temporary server's base URL, such as
  `http://127.0.0.1:54321`; the port is assigned for the test.
- `tt.path` is only a path, such as `/users/1`. Concatenating them creates the
  full URL for the request.
- `defer server.Close()` shuts down the temporary server when the parent test
  finishes.
- The test uses the existing global sample map. Avoid running tests that mutate
  that shared map in parallel until the application is refactored to give each
  test an isolated store.
