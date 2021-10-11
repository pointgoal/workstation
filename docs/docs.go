// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "PointGoal team",
            "url": "https://github.com/pointgoal/workstation",
            "email": "lark@pointgoal.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/org": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organization"
                ],
                "summary": "List organizations",
                "operationId": "1",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ListOrgResponse"
                        }
                    }
                }
            },
            "put": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organization"
                ],
                "summary": "Create organization",
                "operationId": "3",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Organization name",
                        "name": "orgName",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.CreateOrgResponse"
                        }
                    }
                }
            }
        },
        "/v1/org/{orgId}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organization"
                ],
                "summary": "Get organization",
                "operationId": "2",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Organization Id",
                        "name": "orgId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.GetOrgResponse"
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organization"
                ],
                "summary": "Update organization",
                "operationId": "5",
                "parameters": [
                    {
                        "description": "Organization",
                        "name": "org",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.UpdateOrgRequest"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "Organization Id",
                        "name": "orgId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.UpdateOrgResponse"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organization"
                ],
                "summary": "Delete organization",
                "operationId": "4",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Organization Id",
                        "name": "orgId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.DeleteOrgResponse"
                        }
                    }
                }
            }
        },
        "/v1/proj": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "project"
                ],
                "summary": "List projects",
                "operationId": "6",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Organization Id",
                        "name": "orgId",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ListProjResponse"
                        }
                    }
                }
            },
            "put": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "project"
                ],
                "summary": "create project",
                "operationId": "8",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Organization Id",
                        "name": "orgId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "Project",
                        "name": "project",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.CreateProjRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.CreateProjResponse"
                        }
                    }
                }
            }
        },
        "/v1/proj/{projId}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "project"
                ],
                "summary": "Get project",
                "operationId": "7",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Project Id",
                        "name": "projId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.GetProjResponse"
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "project"
                ],
                "summary": "update project",
                "operationId": "10",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Project Id",
                        "name": "projId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Project",
                        "name": "project",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.UpdateProjRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.UpdateProjResponse"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "project"
                ],
                "summary": "delete project",
                "operationId": "9",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Project Id",
                        "name": "projId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.DeleteProjResponse"
                        }
                    }
                }
            }
        },
        "/v1/source": {
            "put": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "source"
                ],
                "summary": "create source",
                "operationId": "11",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Project Id",
                        "name": "projId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "Source",
                        "name": "source",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.CreateSourceRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.CreateSourceResponse"
                        }
                    }
                }
            }
        },
        "/v1/source/{sourceId}": {
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "source"
                ],
                "summary": "delete source",
                "operationId": "12",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Source Id",
                        "name": "sourceId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.DeleteProjResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.CreateOrgResponse": {
            "type": "object",
            "properties": {
                "orgId": {
                    "type": "integer"
                }
            }
        },
        "controller.CreateProjRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "controller.CreateProjResponse": {
            "type": "object",
            "properties": {
                "orgId": {
                    "type": "integer"
                },
                "projId": {
                    "type": "integer"
                }
            }
        },
        "controller.CreateSourceRequest": {
            "type": "object",
            "properties": {
                "repository": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "controller.CreateSourceResponse": {
            "type": "object",
            "properties": {
                "projId": {
                    "type": "integer"
                },
                "sourceId": {
                    "type": "integer"
                }
            }
        },
        "controller.DeleteOrgResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "boolean"
                }
            }
        },
        "controller.DeleteProjResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "boolean"
                }
            }
        },
        "controller.GetOrgResponse": {
            "type": "object",
            "properties": {
                "org": {
                    "$ref": "#/definitions/controller.Org"
                }
            }
        },
        "controller.GetProjResponse": {
            "type": "object",
            "properties": {
                "proj": {
                    "$ref": "#/definitions/controller.Proj"
                }
            }
        },
        "controller.ListOrgResponse": {
            "type": "object",
            "properties": {
                "orgList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controller.Org"
                    }
                }
            }
        },
        "controller.ListProjResponse": {
            "type": "object",
            "properties": {
                "projList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controller.Proj"
                    }
                }
            }
        },
        "controller.Org": {
            "type": "object",
            "properties": {
                "meta": {
                    "$ref": "#/definitions/repository.Org"
                },
                "projIds": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "controller.Proj": {
            "type": "object",
            "properties": {
                "meta": {
                    "$ref": "#/definitions/repository.Proj"
                }
            }
        },
        "controller.UpdateOrgRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "controller.UpdateOrgResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "boolean"
                }
            }
        },
        "controller.UpdateProjRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "controller.UpdateProjResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "boolean"
                }
            }
        },
        "repository.Org": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "repository.Proj": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "orgId": {
                    "type": "integer"
                },
                "source": {
                    "$ref": "#/definitions/repository.Source"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "repository.Source": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "projId": {
                    "type": "integer"
                },
                "repository": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Workstation",
	Description: "This is workstation backend with rk-boot.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
