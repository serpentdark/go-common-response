package response

import (
	"encoding/json"
	"testing"
)

func TestErrorResponses(t *testing.T) {
	tests := []struct {
		name     string
		response *ErrorResponse
		wantCode string
	}{
		{
			name: "Bad Request - BFF",
			response: BadRequest(
				BNYBBadRequest,
				"Invalid request body",
				"req-111-222",
				ErrorIssue{Service: "Web Application", Issue: "Request body malformed"},
			),
			wantCode: "B-NYB-400",
		},
		{
			name: "Unauthorized - Auth Service",
			response: Unauthorized(
				CAUTUnauthorized,
				"Token expired",
				"req-333-444",
				ErrorIssue{Service: "Auth Service", Issue: "JWT token has expired"},
			),
			wantCode: "C-AUT-401",
		},
		{
			name: "Not Found - Listing Service",
			response: NotFound(
				CLSTNotFound,
				"Listing not found",
				"req-555-666",
				ErrorIssue{Service: "Listing Service", Issue: "Listing ID not found"},
			),
			wantCode: "C-LST-404",
		},
		{
			name: "Bad Gateway - Integration Service",
			response: BadGateway(
				XINTBadGateway,
				"External API error",
				"req-777-888",
				ErrorIssue{Service: "Integration Service", Issue: "SMS API returned an error"},
			),
			wantCode: "X-INT-502",
		},
		{
			name: "Payload Too Large - Media Service",
			response: PayloadTooLarge(
				CMEDPayloadTooLarge,
				"File too large",
				"req-999-000",
				ErrorIssue{Service: "Media Service", Issue: "File exceeds 10MB limit"},
			),
			wantCode: "C-MED-413",
		},
		{
			name: "Gateway Timeout - BFF",
			response: GatewayTimeout(
				BNYBGatewayTimeout,
				"Upstream timeout",
				"req-aaa-bbb",
				ErrorIssue{Service: "Web Application", Issue: "Upstream service timeout"},
			),
			wantCode: "B-NYB-504",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check success field is false
			if tt.response.Success != false {
				t.Errorf("Expected success to be false, got %v", tt.response.Success)
			}

			// Check error code
			if tt.response.Error.Code != tt.wantCode {
				t.Errorf("Expected code %s, got %s", tt.wantCode, tt.response.Error.Code)
			}

			// Check error has traceId
			if tt.response.Error.TraceID == "" {
				t.Error("Expected traceId to be set")
			}

			// Check error has timestamp
			if tt.response.Error.Timestamp == "" {
				t.Error("Expected timestamp to be set")
			}

			// Check error has details
			if len(tt.response.Error.Details) == 0 {
				t.Error("Expected details to be set")
			}

			// Marshal to JSON to verify structure
			jsonData, err := json.MarshalIndent(tt.response, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal response: %v", err)
			}

			t.Logf("Response JSON:\n%s\n", string(jsonData))
		})
	}
}

func TestErrorWithMultipleDetails(t *testing.T) {
	response := ValidationFailed(
		BNYBValidationFailed,
		"Validation failed",
		"req-multi-123",
		ErrorIssue{Service: "Web Application", Issue: "Email format invalid"},
		ErrorIssue{Service: "Web Application", Issue: "Password too short"},
		ErrorIssue{Service: "Web Application", Issue: "Username already exists"},
	)

	if len(response.Error.Details) != 3 {
		t.Errorf("Expected 3 details, got %d", len(response.Error.Details))
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	t.Logf("Validation Error Response:\n%s\n", string(jsonData))
}

func TestErrorWithoutDetails(t *testing.T) {
	response := InternalServerError(
		CLSTInternalServerError,
		"Internal server error",
		"req-no-details",
	)

	if response.Error.Details != nil && len(response.Error.Details) > 0 {
		t.Error("Expected details to be empty")
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	t.Logf("Error Without Details:\n%s\n", string(jsonData))
}

func TestErrorIssueErrAndCodeFields(t *testing.T) {
	resp := Unauthorized(
		BBOBUnauthorized,
		"Invalid or expired token",
		"f3a8076c56a6938a68bbc06c2290fe7d",
		NewIssue("Backoffice BFF", "keycloak token verify failed",
			"jwt: token signature is invalid", ""),
	)

	if resp.Error.Code != "B-BOB-401" {
		t.Errorf("want code B-BOB-401, got %s", resp.Error.Code)
	}
	if len(resp.Error.Details) != 1 {
		t.Fatalf("want 1 detail, got %d", len(resp.Error.Details))
	}
	if resp.Error.Details[0].Err != "jwt: token signature is invalid" {
		t.Errorf("want err string preserved, got %q", resp.Error.Details[0].Err)
	}

	out, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Issue with err field:\n%s", string(out))
}

func TestWrapUpstreamChain(t *testing.T) {
	// Imagine: BFF (B-BOB) calls Listing (C-LST) which calls Report (C-RPT)
	// Report blew up first, Listing wrapped it, and now BFF wraps Listing.

	reportErr := InternalServerError(
		CRPTInternalServerError,
		"failed to query monthly report",
		"trace-rpt-1",
		NewIssue("Report Service", "mongo query timeout",
			"context deadline exceeded", ""),
	)

	// Listing wraps Report
	listingErr := WrapUpstream(
		CLSTInternalServerError,
		"downstream report failed",
		"Listing Service",
		"trace-lst-1",
		reportErr,
	)

	// BFF wraps Listing → this is what the FE receives
	bffErr := WrapUpstream(
		BBOBInternalServerError,
		"upstream service error",
		"Backoffice BFF",
		"f3a8076c56a6938a68bbc06c2290fe7d",
		listingErr,
	)

	if bffErr.Error.Code != "B-BOB-500" {
		t.Errorf("want outer code B-BOB-500, got %s", bffErr.Error.Code)
	}
	wantChain := []string{"B-BOB-500", "C-LST-500", "C-RPT-500"}
	if len(bffErr.Error.Chain) != len(wantChain) {
		t.Fatalf("want chain %v, got %v", wantChain, bffErr.Error.Chain)
	}
	for i, hop := range wantChain {
		if bffErr.Error.Chain[i] != hop {
			t.Errorf("chain[%d] want %s got %s", i, hop, bffErr.Error.Chain[i])
		}
	}

	out, _ := json.MarshalIndent(bffErr, "", "  ")
	t.Logf("Three-level chained error:\n%s", string(out))
}

func TestWithChainHelper(t *testing.T) {
	resp := InternalServerError(BBOBInternalServerError, "boom", "trace-1").
		WithChain("B-BOB-500", "C-LST-500", "C-RPT-500")

	if got := len(resp.Error.Chain); got != 3 {
		t.Fatalf("want 3 hops, got %d", got)
	}
	if resp.Error.Chain[2] != "C-RPT-500" {
		t.Errorf("want last hop C-RPT-500, got %s", resp.Error.Chain[2])
	}
}

func TestServiceFromCode(t *testing.T) {
	cases := map[string]string{
		"B-BOB-401": "B-BOB",
		"C-LST-500": "C-LST",
		"C-RPT-404": "C-RPT",
		"X-INT-502": "X-INT",
		"":          "",
		"weird":     "weird",
	}
	for in, want := range cases {
		if got := serviceFromCode(in); got != want {
			t.Errorf("serviceFromCode(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAllServiceErrorCodes(t *testing.T) {
	services := map[string][]string{
		"BFF": {
			BNYBBadRequest, BNYBUnauthorized, BNYBForbidden,
			BNYBNotFound, BNYBConflict, BNYBValidationFailed,
			BNYBTooManyRequests, BNYBInternalServerError,
			BNYBServiceUnavailable, BNYBGatewayTimeout,
		},
		"Auth": {
			CAUTBadRequest, CAUTUnauthorized, CAUTForbidden,
			CAUTNotFound, CAUTConflict, CAUTValidationFailed,
			CAUTTooManyRequests, CAUTInternalServerError,
			CAUTServiceUnavailable,
		},
		"Listing": {
			CLSTBadRequest, CLSTUnauthorized, CLSTForbidden,
			CLSTNotFound, CLSTConflict, CLSTValidationFailed,
			CLSTTooManyRequests, CLSTInternalServerError,
			CLSTServiceUnavailable,
		},
		"Media": {
			CMEDBadRequest, CMEDUnauthorized, CMEDForbidden,
			CMEDNotFound, CMEDConflict, CMEDPayloadTooLarge,
			CMEDValidationFailed, CMEDTooManyRequests,
			CMEDInternalServerError, CMEDServiceUnavailable,
		},
		"Integration": {
			XINTBadRequest, XINTUnauthorized, XINTForbidden,
			XINTNotFound, XINTValidationFailed, XINTTooManyRequests,
			XINTInternalServerError, XINTBadGateway,
			XINTServiceUnavailable, XINTGatewayTimeout,
		},
		"Realtime Search Sync Worker": {
			WRSSBadRequest, WRSSNotFound, WRSSConflict,
			WRSSValidationFailed, WRSSTooManyRequests,
			WRSSInternalServerError, WRSSServiceUnavailable,
			WRSSGatewayTimeout,
		},
		"Product Worker": {
			WPRDBadRequest, WPRDNotFound, WPRDConflict,
			WPRDValidationFailed, WPRDTooManyRequests,
			WPRDInternalServerError, WPRDServiceUnavailable,
			WPRDGatewayTimeout,
		},
		"Listing Worker": {
			WLSTBadRequest, WLSTNotFound, WLSTConflict,
			WLSTValidationFailed, WLSTTooManyRequests,
			WLSTInternalServerError, WLSTServiceUnavailable,
			WLSTGatewayTimeout,
		},
	}

	t.Log("Testing all service error codes:")
	for service, codes := range services {
		t.Logf("\n%s Service - %d error codes defined", service, len(codes))
		for _, code := range codes {
			if code == "" {
				t.Errorf("Empty error code in %s service", service)
			}
		}
	}
}
