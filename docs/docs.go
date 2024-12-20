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
        "/balance": {
            "get": {
                "description": "Page where you can get the balance for a client",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Logic"
                ],
                "summary": "Get Balance",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Page where you can log into the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication/Authorization"
                ],
                "summary": "Login Page",
                "responses": {
                    "200": {
                        "description": "bearer_token",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/metrics": {
            "get": {
                "description": "Page for fetching custom metrics about the server's operation",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Metrics"
                ],
                "summary": "Metrics Page",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/registration": {
            "post": {
                "description": "Page where you can register in the server's database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication/Authorization"
                ],
                "summary": "Registration Page",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/transactions": {
            "get": {
                "description": "Page where you can get a list of transactions for the user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Logic"
                ],
                "summary": "List of Transactions",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            },
            "post": {
                "description": "Page where you can add a new transaction to the list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Logic"
                ],
                "summary": "Add Transaction",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{"http", "https"},
	Title:            "Payment System API",
	Description:      "This is a sample server for Payment System",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
