// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/me": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint returns the profile of the authenticated user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User API"
                ],
                "summary": "Read My Profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/me.ReadMyProfileResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    }
                }
            }
        },
        "/api/v1/me/accounts": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint creates a new account.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Account API"
                ],
                "summary": "Create Account",
                "parameters": [
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/accounts.CreateAccountRequestBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/accounts.CreateAccountResponse"
                        }
                    },
                    "400": {
                        "description": "Validation Failed or Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    }
                }
            }
        },
        "/api/v1/me/accounts/{account_id}/transactions": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint retrieves the transaction history of the specified account.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transaction API"
                ],
                "summary": "List Transactions",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Account ID to be operated.",
                        "name": "account_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The start date for filtering transactions (format: YYYYMMDD).",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The end date for filtering transactions (format: YYYYMMDD).",
                        "name": "to",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Comma-separated transaction types to filter by. Valid values are DEPOSIT, WITHDRAW, and TRANSFER. If not specified, all transaction types are included.",
                        "name": "operationTypes",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The sorting order of transactions based on transactionAt. Valid values are ASC or DESC. Defaults to DESC.",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "The maximum number of transaction histories per page. Can be specified between 1 and 100.",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "The page number for paginated results.",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/transactions.ListTransactionsResponse"
                        }
                    },
                    "400": {
                        "description": "Validation Failed or Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint executes a transaction (deposit, withdraw, or transfer) for the specified account.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transaction API"
                ],
                "summary": "Execute Transaction",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Account ID to be operated.",
                        "name": "account_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/transactions.ExecuteTransactionRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/transactions.ExecuteTransactionResponse"
                        }
                    },
                    "400": {
                        "description": "Validation Failed or Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    }
                }
            }
        },
        "/api/v1/signin": {
            "post": {
                "description": "This endpoint authenticates the user using their email and password, and issues an access token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication API"
                ],
                "summary": "Signin",
                "parameters": [
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/signin.SigninRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/signin.SigninResponse"
                        }
                    },
                    "400": {
                        "description": "Validation Failed or Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    }
                }
            }
        },
        "/api/v1/signup": {
            "post": {
                "description": "This endpoint creates a new user and issues an access token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication API"
                ],
                "summary": "Signup",
                "parameters": [
                    {
                        "description": "Request Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/signup.SignupRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/signup.SignupResponse"
                        }
                    },
                    "400": {
                        "description": "Validation Failed or Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ValidationProblemDetail"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ProblemDetail"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "accounts.CreateAccountRequestBody": {
            "type": "object",
            "properties": {
                "currency": {
                    "description": "The currency for the account. Supported values are JPY or USD.",
                    "type": "string",
                    "example": "JPY"
                },
                "name": {
                    "description": "The name of the account. Must be 3-20 characters long.",
                    "type": "string",
                    "example": "For work"
                },
                "password": {
                    "description": "A 4-digit password for securing the account.",
                    "type": "string",
                    "example": "1234"
                }
            }
        },
        "accounts.CreateAccountResponse": {
            "type": "object",
            "properties": {
                "balance": {
                    "description": "The current balance of the account.",
                    "type": "number",
                    "example": 0
                },
                "currency": {
                    "description": "The currency for the account.",
                    "type": "string",
                    "example": "JPY"
                },
                "id": {
                    "description": "The ID of the account.",
                    "type": "string",
                    "example": "01J9R7YPV1FH1V0PPKVSB5C7LE"
                },
                "name": {
                    "description": "The name of the account.",
                    "type": "string",
                    "example": "For work"
                },
                "updatedAt": {
                    "description": "The date and time the account was last updated.",
                    "type": "string",
                    "example": "2021-08-01T00:00:00Z"
                }
            }
        },
        "me.ReadMyProfileResponse": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "sato@example.com"
                },
                "id": {
                    "type": "string",
                    "example": "01J9R7YPV1FH1V0PPKVSB5C8FW"
                },
                "name": {
                    "type": "string",
                    "example": "Sato Taro"
                }
            }
        },
        "response.ProblemDetail": {
            "type": "object",
            "properties": {
                "detail": {
                    "type": "string",
                    "example": "Error detail message"
                },
                "instance": {
                    "type": "string",
                    "example": "/path/to/resource"
                },
                "status": {
                    "type": "integer",
                    "example": 400
                },
                "title": {
                    "type": "string",
                    "example": "Error title"
                },
                "type": {
                    "type": "string",
                    "example": "https://example.com/probs/error-title"
                }
            }
        },
        "response.ValidationError": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string",
                    "example": "field name"
                },
                "message": {
                    "type": "string",
                    "example": "error message"
                }
            }
        },
        "response.ValidationProblemDetail": {
            "type": "object",
            "properties": {
                "detail": {
                    "type": "string",
                    "example": "Error detail message"
                },
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.ValidationError"
                    }
                },
                "instance": {
                    "type": "string",
                    "example": "/path/to/resource"
                },
                "status": {
                    "type": "integer",
                    "example": 400
                },
                "title": {
                    "type": "string",
                    "example": "Error title"
                },
                "type": {
                    "type": "string",
                    "example": "https://example.com/probs/error-title"
                }
            }
        },
        "signin.SigninRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "description": "The email address of the user, used for login.",
                    "type": "string",
                    "example": "sato@example.com"
                },
                "password": {
                    "description": "The password associated with the email address, required for login. Must be 8-20 characters long.",
                    "type": "string",
                    "example": "password"
                }
            }
        },
        "signin.SigninResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string",
                    "example": "eyJhb..."
                }
            }
        },
        "signup.SignupRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "description": "The email address of the user, used for login.",
                    "type": "string",
                    "example": "sato@example.com"
                },
                "name": {
                    "description": "The name of the user. Must be 3-20 characters long.",
                    "type": "string",
                    "example": "Sato Taro"
                },
                "password": {
                    "description": "The password associated with the email address, required for login. Must be 8-20 characters long.",
                    "type": "string",
                    "example": "password"
                }
            }
        },
        "signup.SignupResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string",
                    "example": "eyJhb..."
                },
                "user": {
                    "$ref": "#/definitions/signup.SignupResponseBodyUser"
                }
            }
        },
        "signup.SignupResponseBodyUser": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "sato@example.com"
                },
                "id": {
                    "type": "string",
                    "example": "01J9R7YPV1FH1V0PPKVSB5C8FW"
                },
                "name": {
                    "type": "string",
                    "example": "Sato Taro"
                }
            }
        },
        "transactions.ExecuteTransactionRequestBody": {
            "type": "object",
            "properties": {
                "amount": {
                    "description": "The transaction amount.",
                    "type": "number",
                    "example": 1000
                },
                "currency": {
                    "description": "The currency of the transaction. Supported values are JPY and USD.",
                    "type": "string",
                    "example": "JPY"
                },
                "operationType": {
                    "description": "Specifies the type of transaction. Valid values are DEPOSIT, WITHDRAW, or TRANSFER.",
                    "type": "string",
                    "example": "DEPOSIT"
                },
                "password": {
                    "description": "The account password.",
                    "type": "string",
                    "example": "1234"
                },
                "receiverAccountId": {
                    "description": "Required for TRANSFER operations. Represents the recipient account ID.",
                    "type": "string",
                    "example": "01J9R8AJ1Q2YDH1X9836GS9D87"
                }
            }
        },
        "transactions.ExecuteTransactionResponse": {
            "type": "object",
            "properties": {
                "accountId": {
                    "type": "string",
                    "example": "01J9R7YPV1FH1V0PPKVSB5C8FW"
                },
                "amount": {
                    "type": "number",
                    "example": 1000
                },
                "currency": {
                    "type": "string",
                    "example": "JPY"
                },
                "id": {
                    "type": "string",
                    "example": "01J9R8AJ1Q2YDH1X9836GS9E89"
                },
                "operationType": {
                    "type": "string",
                    "example": "DEPOSIT"
                },
                "receiverAccountId": {
                    "type": "string",
                    "example": "01J9R8AJ1Q2YDH1X9836GS9D87"
                },
                "transactionAt": {
                    "type": "string",
                    "example": "2024-03-20T15:00:00Z"
                }
            }
        },
        "transactions.ListTransactionsResponse": {
            "type": "object",
            "properties": {
                "total": {
                    "type": "integer",
                    "example": 1
                },
                "transactions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/transactions.ListTransactionsTransaction"
                    }
                }
            }
        },
        "transactions.ListTransactionsTransaction": {
            "type": "object",
            "properties": {
                "accountId": {
                    "type": "string",
                    "example": "01J9R7YPV1FH1V0PPKVSB5C8FW"
                },
                "amount": {
                    "type": "number",
                    "example": 1000
                },
                "currency": {
                    "type": "string",
                    "example": "JPY"
                },
                "id": {
                    "type": "string",
                    "example": "01J9R8AJ1Q2YDH1X9836GS9E89"
                },
                "operationType": {
                    "type": "string",
                    "example": "DEPOSIT"
                },
                "receiverAccountId": {
                    "type": "string",
                    "example": "01J9R8AJ1Q2YDH1X9836GS9D87"
                },
                "transactionAt": {
                    "type": "string",
                    "example": "2024-03-20T15:00:00Z"
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
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "pocgo API",
	Description:      "This is a sample server. Please enter your token in the format: \"Bearer <token>\" in the Authorization header.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
