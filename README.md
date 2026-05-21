# go-common-response

Go module for standardized API response format with a zone-based error code system.

## Installation

```bash
go get github.com/serpentdark/go-common-response@latest
```

Or pin a specific version:

```bash
go get github.com/serpentdark/go-common-response@v0.1.5
```

## Features

- Standardized success and error response structures
- Zone-based error code system (B / C / X)
- **200 pre-defined error codes** — every service covers the full HTTP status spectrum (20 codes each)
- 20 helper functions for common HTTP errors (400, 401, 403, 404, 406, 408, 409, 410, 413, 415, 422, 423, 428, 429, 500, 501, 502, 503, 504, 507)
- Pagination support
- Async job response format

## Error Code Format

```
{ZONE}-{SERVICE}-{HTTP_STATUS}
```

Example: `C-AUT-401` = Core / Auth Service / Unauthorized

## Usage

### Success Response

```go
import "github.com/serpentdark/go-common-response"

// Simple success
resp := response.Success(map[string]string{"id": "123"})

// Paginated list
items := []Product{...}
pagination := response.NewPagination(128, 1, 20)
resp := response.SuccessList(items, pagination)

// Async job
resp := response.SuccessAsync("job_abc123", "processing", 30)
```

### Error Response

```go
// Simple error
resp := response.NotFound(
    response.CLSTNotFound,
    "Listing not found",
    "req-trace-id",
)

// Error with details
resp := response.ValidationFailed(
    response.BNYBValidationFailed,
    "Validation failed",
    "req-trace-id",
    response.ErrorIssue{Service: "Web Application", Issue: "Email invalid"},
    response.ErrorIssue{Service: "Web Application", Issue: "Password too short"},
)

// Newer helpers added in v0.1.3
resp := response.PreconditionRequired(
    response.CMEDPreconditionRequired,
    "If-Match header required for media update",
    "req-trace-id",
)

resp := response.UnsupportedMediaType(
    response.CMEDUnsupportedMediaType,
    "Only image/jpeg and image/png are accepted",
    "req-trace-id",
)

resp := response.InsufficientStorage(
    response.CMEDInsufficientStorage,
    "Storage quota exceeded",
    "req-trace-id",
)
```

## Available Helpers

| Helper | HTTP Status |
|--------|------------|
| `BadRequest` | 400 |
| `Unauthorized` | 401 |
| `Forbidden` | 403 |
| `NotFound` | 404 |
| `NotAcceptable` | 406 |
| `RequestTimeout` | 408 |
| `Conflict` | 409 |
| `Gone` | 410 |
| `PayloadTooLarge` | 413 |
| `UnsupportedMediaType` | 415 |
| `ValidationFailed` | 422 |
| `Locked` | 423 |
| `PreconditionRequired` | 428 |
| `TooManyRequests` | 429 |
| `InternalServerError` | 500 |
| `NotImplemented` | 501 |
| `BadGateway` | 502 |
| `ServiceUnavailable` | 503 |
| `GatewayTimeout` | 504 |
| `InsufficientStorage` | 507 |

All helpers share the same signature:

```go
func Helper(code, message, traceID string, details ...ErrorIssue) *ErrorResponse
```

## Error Codes

### Zone Codes
- **B** — Web Application (BFF)
- **C** — Core Services
- **X** — Integration Services

### Services & Prefixes

| Prefix | Service |
|--------|---------|
| `B-NYB` | Nayoo Web BFF |
| `B-BOB` | Backoffice BFF |
| `C-AUT` | Auth Service |
| `C-LST` | Listing Service |
| `C-MED` | Media Service |
| `C-SRC` | Search Service |
| `C-NOT` | Notification Service |
| `C-RPT` | Report Service |
| `C-PRD` | Product Service |
| `X-INT` | Integration Service |

### Codes per Service

Every service block above exposes the **same 20 HTTP status codes**, named with the pattern `<Prefix><Name>`:

```
<Prefix>BadRequest            // 400
<Prefix>Unauthorized          // 401
<Prefix>Forbidden             // 403
<Prefix>NotFound              // 404
<Prefix>NotAcceptable         // 406
<Prefix>RequestTimeout        // 408
<Prefix>Conflict              // 409
<Prefix>Gone                  // 410
<Prefix>PayloadTooLarge       // 413
<Prefix>UnsupportedMediaType  // 415
<Prefix>ValidationFailed      // 422
<Prefix>Locked                // 423
<Prefix>PreconditionRequired  // 428
<Prefix>TooManyRequests       // 429
<Prefix>InternalServerError   // 500
<Prefix>NotImplemented        // 501
<Prefix>BadGateway            // 502
<Prefix>ServiceUnavailable    // 503
<Prefix>GatewayTimeout        // 504
<Prefix>InsufficientStorage   // 507
```

For example, the Auth service exposes `CAUTBadRequest` → `"C-AUT-400"`, `CAUTPreconditionRequired` → `"C-AUT-428"`, and so on for all 20 statuses. Same shape for every prefix.

**Total: 10 services × 20 statuses = 200 codes.**

### Adding a New Service

Copy any existing service block in `response.go` and rename the prefix. Example for a hypothetical Tracking Service (`C-TRK`):

```go
// C-TRK - Tracking Service
CTRKBadRequest            = "C-TRK-400"
CTRKUnauthorized          = "C-TRK-401"
// ... (all 20 statuses)
CTRKInsufficientStorage   = "C-TRK-507"
```

## Response Examples

### Success
```json
{
  "success": true,
  "data": {
    "id": "123",
    "name": "Example"
  }
}
```

### Error
```json
{
    "success": false,
    "error": {
        "code": "B-BOB-401",
        "message": "Invalid or expired token",
        "details": [
            {
                "service": "Backoffice BFF",
                "issue": "keycloak token verify failed"
            }
        ],
        "traceId": "f3a8076c56a6938a68bbc06c2290fe7d",
        "timestamp": "2026-05-20T04:41:52Z"
    }
}
```

## Testing

```bash
go test -v
```

## Changelog

### v0.1.5
- Added `B-BOB` Backoffice BFF error code prefix with the full 20-status set.
- Updated documentation examples to distinguish Backoffice BFF (`B-BOB`) from Nayoo Web BFF (`B-NYB`).

### v0.1.3
- Expanded error codes to cover full HTTP status spectrum (180 codes total)
- Added 8 new helper functions: `NotAcceptable`, `RequestTimeout`, `Gone`, `UnsupportedMediaType`, `Locked`, `PreconditionRequired`, `NotImplemented`, `InsufficientStorage`
- Filled missing 413/502/504 on every service block
- 100% backward compatible — existing codes and helpers unchanged

## License

MIT
