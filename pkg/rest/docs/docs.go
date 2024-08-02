// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Project Korrel8r",
            "url": "https://github.com/korrel8r/korrel8r"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "https://github.com/korrel8r/korrel8r/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/domains": {
            "get": {
                "summary": "Get name, configuration and status for each domain.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/Domain"
                            }
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/domains/{domain}/classes": {
            "get": {
                "summary": "Get class names and descriptions for a domain.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Domain name",
                        "name": "domain",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Classes"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/graphs/goals": {
            "post": {
                "summary": "Create a correlation graph from start objects to goal queries.",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "include rules in graph edges",
                        "name": "rules",
                        "in": "query"
                    },
                    {
                        "description": "search from start to goal classes",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Goals"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Graph"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/graphs/neighbours": {
            "post": {
                "summary": "Create a neighbourhood graph around a start object to a given depth.",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "include rules in graph edges",
                        "name": "rules",
                        "in": "query"
                    },
                    {
                        "description": "search from neighbours",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Neighbours"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Graph"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/lists/goals": {
            "post": {
                "summary": "Create a list of goal nodes related to a starting point.",
                "parameters": [
                    {
                        "description": "search from start to goal classes",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Goals"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/Node"
                            }
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/objects": {
            "get": {
                "summary": "Execute a query, returns a list of JSON objects.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "query string",
                        "name": "query",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            }
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "Classes": {
            "description": "Classes is a map from class names to a short description.",
            "type": "object",
            "additionalProperties": {
                "type": "string"
            }
        },
        "Domain": {
            "description": "Domain configuration information.",
            "type": "object",
            "properties": {
                "name": {
                    "description": "Name of the domain.",
                    "type": "string"
                },
                "stores": {
                    "description": "Stores configured for the domain.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Store"
                    }
                }
            }
        },
        "Edge": {
            "description": "Directed edge in the result graph, from Start to Goal classes.",
            "type": "object",
            "properties": {
                "goal": {
                    "description": "Goal is the class name of the goal node.",
                    "type": "string",
                    "example": "domain:class"
                },
                "rules": {
                    "description": "Rules is the set of rules followed along this edge.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Rule"
                    },
                    "x-omitempty": true
                },
                "start": {
                    "description": "Start is the class name of the start node.",
                    "type": "string"
                }
            }
        },
        "Goals": {
            "type": "object"
        },
        "Graph": {
            "description": "Graph resulting from a correlation search.",
            "type": "object",
            "properties": {
                "edges": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Edge"
                    }
                },
                "nodes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Node"
                    }
                }
            }
        },
        "Neighbours": {
            "description": "Starting point for a neighbours search.",
            "type": "object",
            "properties": {
                "depth": {
                    "description": "Max depth of neighbours graph.",
                    "type": "integer"
                },
                "start": {
                    "$ref": "#/definitions/Start"
                }
            }
        },
        "Node": {
            "type": "object",
            "properties": {
                "class": {
                    "description": "Class is the full class name in \"DOMAIN:CLASS\" form.",
                    "type": "string",
                    "example": "domain:class"
                },
                "count": {
                    "description": "Count of results found for this class, after de-duplication.",
                    "type": "integer"
                },
                "queries": {
                    "description": "Queries yielding results for this class.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/QueryCount"
                    }
                }
            }
        },
        "QueryCount": {
            "description": "Query run during a correlation with a count of results found.",
            "type": "object",
            "properties": {
                "count": {
                    "description": "Count of results or -1 if the query was not executed.",
                    "type": "integer"
                },
                "query": {
                    "description": "Query for correlation data.",
                    "type": "string"
                }
            }
        },
        "Rule": {
            "type": "object",
            "properties": {
                "name": {
                    "description": "Name is an optional descriptive name.",
                    "type": "string"
                },
                "queries": {
                    "description": "Queries generated while following this rule.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/QueryCount"
                    }
                }
            }
        },
        "Start": {
            "type": "object"
        },
        "Store": {
            "description": "Store is a map of name:value attributes used to connect to a store.",
            "type": "object",
            "additionalProperties": {
                "type": "string"
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "v1alpha1",
	Host:             "localhost:8080",
	BasePath:         "/api/v1alpha1",
	Schemes:          []string{"https", "http"},
	Title:            "REST API",
	Description:      "REST API for the Korrel8r correlation engine.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
