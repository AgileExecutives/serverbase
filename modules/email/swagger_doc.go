// Package email provides pre-written swagger documentation for the email
// endpoints.  The spec is registered with the server's DocRegistry during
// module Initialize so that swagger.MergeAndRegister picks it up.
package email

// EmailSwaggerJSON is the Swagger 2.0 spec for the email module endpoints.
const EmailSwaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "Email API",
    "description": "Email management and notification endpoints.",
    "version": "1.0.0"
  },
  "basePath": "/api/v1",
  "paths": {
    "/emails": {
      "get": {
        "security": [{"BearerAuth": []}],
        "summary": "List emails",
        "description": "Returns a paginated list of emails. Optionally filter by status.",
        "produces": ["application/json"],
        "tags": ["emails"],
        "operationId": "getEmails",
        "parameters": [
          {"name": "page",   "in": "query", "type": "integer", "description": "Page number (default 1)"},
          {"name": "limit",  "in": "query", "type": "integer", "description": "Items per page (default 10)"},
          {"name": "status", "in": "query", "type": "string",  "description": "Filter by status: pending | sent | delivered | failed"}
        ],
        "responses": {
          "200": {"description": "Paginated email list"},
          "401": {"description": "Unauthorized"},
          "500": {"description": "Internal server error"}
        }
      }
    },
    "/emails/{id}": {
      "get": {
        "security": [{"BearerAuth": []}],
        "summary": "Get email by ID",
        "produces": ["application/json"],
        "tags": ["emails"],
        "operationId": "getEmail",
        "parameters": [
          {"name": "id", "in": "path", "required": true, "type": "integer", "description": "Email ID"}
        ],
        "responses": {
          "200": {"description": "Email object"},
          "400": {"description": "Invalid ID"},
          "401": {"description": "Unauthorized"},
          "404": {"description": "Email not found"},
          "500": {"description": "Internal server error"}
        }
      }
    },
    "/emails/send": {
      "post": {
        "security": [{"BearerAuth": []}],
        "summary": "Send an email",
        "description": "Queues and dispatches an email asynchronously.",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "tags": ["emails"],
        "operationId": "sendEmail",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": ["to", "subject"],
              "properties": {
                "to":       {"type": "string", "example": "user@example.com"},
                "from":     {"type": "string", "example": "noreply@example.com"},
                "subject":  {"type": "string", "example": "Welcome"},
                "body":     {"type": "string", "description": "Plain-text body"},
                "html_body": {"type": "string", "description": "HTML body"}
              }
            }
          }
        ],
        "responses": {
          "201": {"description": "Email queued"},
          "400": {"description": "Invalid request"},
          "401": {"description": "Unauthorized"},
          "500": {"description": "Internal server error"}
        }
      }
    },
    "/emails/stats": {
      "get": {
        "security": [{"BearerAuth": []}],
        "summary": "Get email statistics",
        "produces": ["application/json"],
        "tags": ["emails"],
        "operationId": "getEmailStats",
        "responses": {
          "200": {
            "description": "Email counts by status",
            "schema": {
              "type": "object",
              "properties": {
                "total":     {"type": "integer"},
                "pending":   {"type": "integer"},
                "sent":      {"type": "integer"},
                "delivered": {"type": "integer"},
                "failed":    {"type": "integer"}
              }
            }
          },
          "401": {"description": "Unauthorized"}
        }
      }
    },
    "/emails/latest-emails": {
      "get": {
        "summary": "Get latest mock emails",
        "description": "Only available when MOCK_EMAIL=true. Returns emails captured by the mock email service.",
        "produces": ["application/json"],
        "tags": ["emails"],
        "operationId": "getLatestMockEmails",
        "responses": {
          "200": {"description": "List of captured mock emails"},
          "503": {"description": "Mock email not enabled"}
        }
      }
    }
  },
  "securityDefinitions": {
    "BearerAuth": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    }
  },
  "definitions": {}
}`
