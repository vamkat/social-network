# Custom Types Documentation

This package provides custom types for the social network backend, ensuring type safety, validation, and consistent serialization.

## Overview

All custom types implement the `Validator` interface, which requires a `Validate() error` method. Use `ValidateStruct()` to validate entire structs containing these types.

### Validator Interface

```go
type Validator interface {
    Validate() error
}
```

### ValidateStruct Function

Validates structs by iterating over exported fields and checking those that implement the `Validator` interface.

**Behavior**:
- Calls `Validate()` on fields implementing `Validator`.
- For fields without the `validate:"nullable"` tag, zero values are treated as errors.
- Nullable fields skip validation if empty.
- Primitives are excluded except slices of custom types.
- If a field is a slice of a custom type, if a null value is found in that slice validation error is returned.

**Tags**:
- `validate:"nullable"`: Marks the field as optional; zero values are allowed and skip validation.

**Example**:

```go
type RegisterRequest struct {
    Username  customtypes.Username  `json:"username,omitempty" validate:"nullable"` // optional
    FirstName customtypes.Name      `json:"first_name,omitempty" validate:"nullable"` // optional
    LastName  customtypes.Name      `json:"last_name"` // required
    About     customtypes.About     `json:"about"`     // required
    Email     customtypes.Email     `json:"email,omitempty" validate:"nullable"` // optional
}

err := customtypes.ValidateStruct(req)
if err != nil {
    // handle validation errors
}
```

**Notes**: Slice validation code is commented out in the implementation. Unexported fields are skipped.

## Types

### Id

**Description**: Represents an encrypted ID (int64). Allows null values in JSON but encrypts to a hash using hashids.

**Validation**: Must be > 0.

**Marshal/Unmarshal**: Marshals to encrypted string using hashids (requires `ENC_KEY` env var). Unmarshals from hash back to int64.

**Usage**: For secure ID transmission in APIs. Implements `Scan()` and `Value()` methods for use in postgress database calls.

### UnsafeId

**Description**: Represents an unencrypted ID (int64).

**Validation**: Must be > 0.

**Marshal/Unmarshal**: Standard JSON marshal/unmarshal as int64.

**Usage**: For internal use where encryption is not needed. Implements `Scan()` and `Value()` methods for use in postgress database calls.

### Ids

**Description**: Slice of `Id`.

**Validation**: Non-empty slice, all IDs must be valid.

**Marshal/Unmarshal**: As []int64.

**Usage**: For lists of IDs. Implements `Scan()` and `Value()` methods for use in postgress database calls.

**Extra Features**: Implements `Unique()` method that returns a copy of type 'Ids' with only the unique entries of the given instance.

### About

**Description**: Bio or description text.

**Validation**: 3-300 characters, no control characters (except \n, \r, \t).

**Marshal/Unmarshal**: Standard string.

**Usage**: User bios.

### Audience

**Description**: Visibility level for posts/comments/events.

**Validation**: Must be one of: "everyone", "group", "followers", "selected" (case-insensitive), no control characters.

**Marshal/Unmarshal**: Standard string.

**Usage**: Content visibility.

### PostBody

**Description**: Body text for posts.

**Validation**: 3-500 characters, no control characters.

**Marshal/Unmarshal**: Standard string.

**Usage**: Post content.

### CommentBody

**Description**: Body text for comments.

**Validation**: 3-400 characters, no control characters.

**Marshal/Unmarshal**: Standard string.

**Usage**: Comment content.

### EventBody

**Description**: Body text for events.

**Validation**: 3-400 characters, no control characters.

**Marshal/Unmarshal**: Standard string.

**Usage**: Event descriptions.

### MsgBody

**Description**: Body text for messages.

**Validation**: 3-400 characters, no control characters.

**Marshal/Unmarshal**: Standard string.

**Usage**: Chat messages. Implements `Scan()` and `Value()` methods for use in postgress database calls.

### CtxKey

**Description**: Type alias for context keys to enforce naming conventions.

**Validation**: N/A.

**Marshal/Unmarshal**: N/A.

**Usage**: Context keys like `ClaimsKey`, `UserId`, etc.

### DateOfBirth

**Description**: User's date of birth.

**Validation**: Not zero, not in future, age 13-120.

**Marshal/Unmarshal**: "2006-01-02" format.

**Usage**: User profiles.

### EventDateTime

**Description**: Date and time for events.

**Validation**: Not zero, not in past, within 6 months ahead.

**Marshal/Unmarshal**: RFC3339.

**Usage**: Event scheduling.

### GenDateTime

**Description**: Generic nullable datetime.

**Validation**: Not zero.

**Marshal/Unmarshal**: RFC3339.

**Usage**: Timestamps like created_at. Implements `Scan()` and `Value()` methods for use in postgress database calls.

### Email

**Description**: Email address.

**Validation**: Valid email regex, no control characters.

**Marshal/Unmarshal**: Standard string.

**Usage**: User emails.

### Username

**Description**: Username.

**Validation**: 3-32 chars, alphanumeric + underscore, no control characters.

**Marshal/Unmarshal**: Standard string.

**Usage**: Usernames.

### Identifier

**Description**: Username or email.

**Validation**: Matches username or email regex.

**Marshal/Unmarshal**: Standard string.

**Usage**: Login identifiers.

### Name

**Description**: First or last name.

**Validation**: Min 2 chars, Unicode letters + apostrophes/hyphens/spaces, no control characters.

**Marshal/Unmarshal**: Standard string.

**Usage**: User names.

### Limit

**Description**: Pagination limit.

**Validation**: 1-500.

**Marshal/Unmarshal**: As int32.

**Usage**: API pagination. Implements `Scan()` and `Value()` methods for use in postgress database calls.

### Offset

**Description**: Pagination offset.

**Validation**: >= 0.

**Marshal/Unmarshal**: As int32.

**Usage**: API pagination. Implements `Scan()` and `Value()` methods for use in postgress database calls.

### Password

**Description**: Plain password.

**Validation**: 8-64 chars, requires uppercase, lowercase, digit, symbol, no control characters.

**Marshal/Unmarshal**: Marshals to "********".

**Usage**: Password input.

### HashedPassword

**Description**: Hashed password.

**Validation**: Non-empty, no control characters.

**Marshal/Unmarshal**: Marshals to "********".

**Usage**: Stored passwords.

### SearchTerm

**Description**: Search query term.

**Validation**: Min 2 chars, alphanumeric + spaces/hyphens, no control characters.

**Marshal/Unmarshal**: Standard string.

**Usage**: Search inputs.

### Title

**Description**: Title for groups/chats.

**Validation**: 1-50 chars, no control characters.

**Marshal/Unmarshal**: Standard string.

**Usage**: Group/chat titles.