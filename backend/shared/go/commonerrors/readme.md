# commonerrors

`commonerrors` provides a **structured, opinionated error model** for Go backend services, with a strong focus on:

* **Stable error classification**
* **Explicit error wrapping**
* **Stack trace for current error origin**
* **Clear separation between internal diagnostics and public messages**
* **Idiomatic integration with `errors.Is` / `errors.As`**
* **First-class gRPC status mapping**

This package is designed for **internal service use**, especially in gRPC-based systems, where consistent error semantics and debuggability matter more than minimal error strings.

---

## Design Goals

* Separate **classification** (what kind of error) from **cause** (what actually happened)
* Preserve the **full causal chain** of errors
* Capture **stack traces only at the point of failure**
* Support Go’s standard error inspection (`errors.Is`, `errors.As`)
* Provide a single, consistent mapping to **gRPC status codes**
* Avoid string-based error handling and ad-hoc conventions

---

## Error Classification

Error *classes* are represented as **sentinel `error` values**, similar to `os` and gRPC errors:

```go
ErrNotFound
ErrInvalidArgument
ErrInternal
ErrUnavailable
// ...
```

These sentinels are used as **error classes**, not typically returned directly.

Each error class maps 1:1 to a gRPC `codes.Code` via internal maps.

---

## The `Error` Type

```go
type Error struct {
	class      error  // Classification (never nil)
	input     string // Input / context at the wrapping site
	stack     string // Stack trace (captured only at origin)
	err       error  // Wrapped cause
	publicMsg string // Message safe for clients
}
```

### Invariants

* `class` is **guaranteed to be non-nil**
* Stack traces are captured **only at the root cause**
* Wrapping preserves the **original causal chain**
* Errors are fully compatible with `errors.Is` and `errors.As`

---

## Error String Representation

Calling `Error()` (or logging the error) produces a **recursive, flattened string** containing:

* The error class
* The input/context at each wrapping level
* The stack trace (only the origin func per linked error)
* The full wrapped error chain

```go
log.Println(err)
```

This prints the **entire causal chain**. This behavior is **intentional** and meant for **internal logging and debugging**.

---

## Creating Errors

### Creating a new error (origin point)

```go
err := New(ErrInternal, sqlErr, "query users")
```

* Captures a stack trace
* Sets the error class
* Records input/context

---

### Wrapping an existing error

```go
err = Wrap(ErrNotFound, err, "GetUser")
```

* Preserves the original stack trace
* Updates or inherits the error class
* Adds contextual input

---

### Wrapping without changing the class

```go
err = Wrap(nil, err, "repository layer")
```

* Keeps the existing classification
* Adds context only

#### It is not recommended to use a `commonerrors` type as a wraped error.

---

### Public-facing messages

```go
err := Wrap(ErrUnauthenticated, tokenErr, "validate token").
	WithPublic("authentication required")
```

* `publicMsg` is **never used internally**
* It is only consumed by transport adapters (e.g. gRPC)

---

## Error Inspection

### Classification (`errors.Is`)

```go
if errors.Is(err, ErrNotFound) {
	// handle not found
}
```
The `Is()` custom method only checks the outer most error `commonerrors.class`. The intention is each error to effectivelly contain only one clasification.


### Accessing the custom error (`errors.As`)

```go
var e *commonerrors.Error
if errors.As(err, &e) {
	log.Println(e)
}
```

---

## Root Cause Extraction

Retrieve the **original underlying error message**:

```go
cause := GetSource(err)
```

This walks the unwrap chain until the final error.

---

## gRPC Integration

### Convert to gRPC code

```go
code := ToGRPCCode(err)
```

Behavior:

* Context cancellation & deadlines are handled explicitly
* Domain errors map via internal code maps
* Fallback is `codes.Unknown`

---

### Convert to gRPC status error

```go
return GRPCStatus(err)
```

Behavior:

* Existing gRPC status errors are propagated
* Most outer error class prevails over nested error classes
* Context errors are mapped first
* Domain errors use `publicMsg` when present
* Empty public messages fall back to a safe default
* Internal details are never exposed

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

* **Create errors at failure boundaries**
* **Wrap errors at layer boundaries**
* Use `WithPublic` only for client-safe messages

---

## Non-Goals

The following are intentionally out of scope:

* HTTP-specific error formatting
* Localization / i18n
* Automatic logging
* Error mutation after creation
* Structured logging output

These can be layered on top without changing the core API.

---

## Summary

`commonerrors` provides a **predictable, explicit, and debuggable** error model that:

* Scales cleanly across service boundaries
* Preserves full causal context
* Keeps public and internal concerns separate
* Plays nicely with Go’s standard error tooling

It favors **clarity and correctness** over minimalism, making it well-suited for production backend systems.

