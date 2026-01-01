# commonerrors

`commonerrors` provides a structured error model for Go services with:

* **Stable error classification** (sentinel errors)
* **Error wrapping** with preserved root cause
* **Public vs internal messages**
* **Idiomatic integration with `errors.Is` / `errors.As`**
* **First-class gRPC status mapping**

This package is designed for backend services (especially gRPC) that need **consistent error semantics across layers**.

---

## Design Goals

* Separate **classification** (what kind of error) from **cause** (what actually happened)
* Preserve the **original underlying error**
* Support Go’s standard error inspection (`errors.Is`, `errors.As`)
* Provide a single source of truth for **gRPC status codes**
* Avoid string-based error handling

---

## Error Classification

Error *classes* are represented as **sentinel `error` values**, similar to gRPC and `os` errors:

```go
ErrNotFound
ErrInvalidArgument
ErrInternal
ErrUnavailable
// ...
```

These are used as **error codes**, not returned directly in most cases.

They map 1:1 to gRPC status codes.

---

## The `Error` Type

```go
type Error struct {
	Code      error  // Classification (never nil)
	Err       error  // Wrapped cause
	Msg       string // Internal context
	PublicMsg string // Message safe for clients
	Source    error  // Original underlying error
}
```

### Invariants

* `Code` is **guaranteed to be non-nil**
* `Source` always points to the **first underlying error**
* Errors are **fully compatible** with `errors.Is` and `errors.As`

---

## Creating Errors

### Wrapping an error

```go
err := Wrap(ErrNotFound, sql.ErrNoRows, "user lookup failed")
```

### Wrapping without changing the code

```go
err = Wrap(nil, err, "additional context")
```

### Public-facing messages

```go
err := Wrap(ErrUnauthenticated, tokenErr).
	WithPublic("authentication required")
```

---

## Error Inspection

### Classification (`errors.Is`)

```go
if errors.Is(err, ErrNotFound) {
	// handle not found
}
```

### Accessing the custom error

```go
var e *commonerrors.Error
if errors.As(err, &e) {
	log.Println(e.Msg)
}
```

---

## Root Cause Extraction

Retrieve the **original underlying error**:

```go
cause := GetSource(err)
```

This works even across multiple wraps.

---

## gRPC Integration

### Convert to gRPC code

```go
code := ToGRPCCode(err)
```

### Convert to gRPC status error

```go
return GRPCStatus(err)
```

Behavior:

* Context cancellation & deadlines are handled
* Domain errors map via `errorToGRPC`
* Public messages are used when available
* Unknown errors fall back to `codes.Unknown`

---

## Example

```go
func GetUser(id string) error {
	user, err := repo.Find(id)
	if err != nil {
		return Wrap(ErrNotFound, err, "GetUser").
			WithPublic("user not found")
	}
	return nil
}
```

---

## Best Practices

* **Always wrap errors at boundaries**
* Use `ErrUnknown` only as a fallback
* Prefer `errors.Is` over string comparison
* Use `PublicMsg` only for client-safe text
* Do not construct `Error` manually — use `Wrap`

---

## Non-Goals

* Stack traces
* Localization
* HTTP-specific error formatting
* Automatic logging

These are intentionally left to higher layers.

---

## Summary

`commonerrors` provides a **minimal, predictable, and Go-idiomatic** error model that scales cleanly across service boundaries while remaining easy to reason about.

If you want:

* immutability guarantees
* HTTP helpers
* middleware examples
* test utilities

those can be layered on without changing the core API.
