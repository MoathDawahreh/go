# Middleware & Path Parameter Context Implementation

## Overview

This implementation demonstrates **two key architectural patterns** for production microservices:

1. **Middleware Injection at Different Levels** - Using Chi's Route grouping to apply middleware at different route levels
2. **Path Parameter Context** - Validating and extracting path parameters in middleware, making them available through context

---

## Files Created/Updated

### New Middleware Files

#### 1. `/internal/middleware/auth.go`

```go
// AuthMiddleware checks for authorization header
// Applied to ALL user and media routes
func AuthMiddleware(next http.Handler) http.Handler
```

- Validates authorization header
- Passes authenticated user info to context

#### 2. `/internal/middleware/path_params.go`

```go
// ValidateIDMiddleware validates and extracts the ID path parameter
// Applied ONLY to routes with {id}
func ValidateIDMiddleware(next http.Handler) http.Handler
```

- Validates that ID is a valid integer
- **Stores validated ID in context** for handler use
- Prevents invalid IDs from reaching handlers

#### 3. `/internal/middleware/logging.go`

```go
// LoggingMiddleware logs HTTP request details and response time
func LoggingMiddleware(next http.Handler) http.Handler
```

- Logs all requests with timestamps
- Measures and logs response duration

---

## How It Works

### Route Structure with Middleware

```
/users (All user routes)
├─ middleware.LoggingMiddleware    ← Applied to ALL /users routes
├─ middleware.AuthMiddleware       ← Applied to ALL /users routes
│
├─ POST / (CreateUser)             ← Validation happens at this level
├─ GET /  (GetAllUsers)            ← Validation happens at this level
│
└─ /{id} (ID-specific routes)
   ├─ middleware.ValidateIDMiddleware ← Applied ONLY to /{id} routes
   ├─ GET / (GetUser)              ← ID already validated & in context
   ├─ PUT / (UpdateUser)           ← ID already validated & in context
   └─ DELETE / (DeleteUser)        ← ID already validated & in context
```

### Path Parameter Context Flow

**Before (old way):**

```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")              // ← Manual extraction
    id, err := strconv.Atoi(idStr)              // ← Manual validation
    if err != nil {
        http.Error(w, "Invalid user ID", ...)   // ← Repeated in every handler
        return
    }
    // ... handler logic
}
```

**After (with context):**

```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    // ← Already validated in middleware
    id, ok := r.Context().Value("userID").(int)
    if !ok {
        http.Error(w, "Invalid user ID", ...)
        return
    }
    // ... handler logic (much cleaner!)
}
```

---

## Benefits

### 1. **Middleware Isolation**

- `LoggingMiddleware` and `AuthMiddleware` run for ALL routes
- `ValidateIDMiddleware` runs ONLY for routes with `{id}`
- No middleware pollution in handlers

### 2. **Single Responsibility**

- `ValidateIDMiddleware` handles validation once
- All handlers that need `{id}` automatically get validated ID via context
- Validation logic is centralized

### 3. **Context Safety**

- Type-safe extraction: `r.Context().Value("userID").(int)`
- Middleware ensures value exists before handler runs
- No nil pointer issues

### 4. **DRY (Don't Repeat Yourself)**

- Validation logic written once in middleware
- 5+ handlers benefit from it automatically
- Easy to add new handlers with same pattern

### 5. **Scalable for Real Microservices**

- Add authentication middleware
- Add rate limiting middleware
- Add request logging middleware
- All without modifying handlers

---

## Production-Ready Patterns

### Adding Global Middleware

```go
r.Route("/users", func(r chi.Router) {
    r.Use(middleware.LoggingMiddleware)  // Global
    r.Use(middleware.AuthMiddleware)     // Global
    r.Route("/{id}", func(r chi.Router) {
        r.Use(middleware.ValidateIDMiddleware)  // Scoped to /{id}
        // routes here...
    })
})
```

### Middleware Composition

```go
r.Route("/api", func(r chi.Router) {
    r.Use(middleware.LoggingMiddleware)
    r.Use(middleware.AuthMiddleware)
    r.Use(middleware.RateLimiting)        // Easy to add!

    r.Route("/admin", func(r chi.Router) {
        r.Use(middleware.AdminAuthMiddleware)  // Extra middleware for admin
        // admin routes...
    })
})
```

---

## Testing Benefits

### Testing Middleware in Isolation

```go
func TestValidateIDMiddleware(t *testing.T) {
    req := httptest.NewRequest("GET", "/users/invalid", nil)
    w := httptest.NewRecorder()

    handler := middleware.ValidateIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := r.Context().Value("userID")
        // assert id is in context
    }))

    handler.ServeHTTP(w, req)
}
```

### Testing Handlers Without Middleware Overhead

```go
func TestGetUser(t *testing.T) {
    req := httptest.NewRequest("GET", "/users/1", nil)
    ctx := context.WithValue(req.Context(), "userID", 1)
    req = req.WithContext(ctx)

    // Handler runs without middleware dependencies
    h := NewUserHandler(service)
    h.GetUser(w, req)
}
```

---

## Summary

This architecture is **production-grade** because it:

✅ Separates concerns (middleware ≠ handlers)
✅ Prevents code duplication (validation written once)
✅ Scales easily (add new handlers, they inherit middleware)
✅ Maintains clean request flow (context-driven data passing)
✅ Supports testing (middleware testable independently)
✅ Follows Chi best practices (Route grouping with Use)

Perfect foundation for a real microservice! 🎯
