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
        "/v1/signup": {
            "post": {
                "description": "Signup",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Signup",
                "parameters": [
                    {
                        "description": "SignupRequestBody",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/signup.SignupRequestBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/signup.SignupResponseBody"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "response.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string",
                    "example": "ErrorCode"
                },
                "message": {
                    "type": "string",
                    "example": "error message"
                }
            }
        },
        "signup.SignupRequestBody": {
            "type": "object",
            "required": [
                "user"
            ],
            "properties": {
                "user": {
                    "$ref": "#/definitions/signup.SignupRequestBodyUser"
                }
            }
        },
        "signup.SignupRequestBodyAccount": {
            "type": "object",
            "required": [
                "currency",
                "name",
                "password"
            ],
            "properties": {
                "currency": {
                    "type": "string",
                    "enum": [
                        "JPY"
                    ],
                    "example": "JPY"
                },
                "name": {
                    "type": "string",
                    "maxLength": 10,
                    "minLength": 1,
                    "example": "For work"
                },
                "password": {
                    "type": "string",
                    "example": "1234"
                }
            }
        },
        "signup.SignupRequestBodyUser": {
            "type": "object",
            "required": [
                "account",
                "email",
                "name",
                "password"
            ],
            "properties": {
                "account": {
                    "$ref": "#/definitions/signup.SignupRequestBodyAccount"
                },
                "email": {
                    "type": "string",
                    "example": "sato@example.com"
                },
                "name": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 1,
                    "example": "Sato Taro"
                },
                "password": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 8,
                    "example": "password"
                }
            }
        },
        "signup.SignupResponseBody": {
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
        "signup.SignupResponseBodyAccount": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number",
                    "example": 0
                },
                "currency": {
                    "type": "string",
                    "example": "JPY"
                },
                "id": {
                    "type": "string",
                    "example": "01J9R7YPV1FH1V0PPKVSB5C7LE"
                },
                "name": {
                    "type": "string",
                    "example": "For work"
                },
                "updatedAt": {
                    "type": "string",
                    "example": "2021-08-01T00:00:00Z"
                }
            }
        },
        "signup.SignupResponseBodyUser": {
            "type": "object",
            "properties": {
                "account": {
                    "$ref": "#/definitions/signup.SignupResponseBodyAccount"
                },
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
	Title:            "Pocgo API",
	Description:      "This is a sample server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}