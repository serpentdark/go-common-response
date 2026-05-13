# go-common-response

Go module for standardized API response format with error code system.

## Installation

```bash
go get github.com/serpentdark/go-common-response
```

## Features

- Standardized success and error response structures
- Zone-based error code system (B, C, X)
- Pagination support
- Async job response format
- Pre-defined error codes for all services

## Error Code Format

```
{ZONE}-{SERVICE}-{CATEGORY}
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
```

## Error Codes

### Zone Codes
- **B** - Web Application (BFF)
- **C** - Core Services
- **X** - Integration Services

### Services
- `B-NYB` - Web Application
- `C-AUT` - Auth Service
- `C-LST` - Listing Service
- `C-MED` - Media Service
- `C-SRC` - Search Service
- `C-NOT` - Notification Service
- `C-RPT` - Report Service
- `X-INT` - Integration Service

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
    "code": "C-LST-404",
    "message": "Listing not found",
    "details": [
      {
        "service": "Listing Service",
        "issue": "Listing ID not found"
      }
    ],
    "traceId": "req-555-666",
    "timestamp": "2026-04-10T08:16:40Z"
  }
}
```

## Testing

```bash
go test -v
```

## License

MIT
