package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "API profesional para gestión de MinIO (CRUD de objetos), con validaciones, logging, métricas y respuestas estandarizadas.",
        "title": "MinIO Management API",
        "termsOfService": "https://example.com/terms",
        "contact": {
            "name": "Soporte API",
            "url": "https://example.com/support",
            "email": "soporte@example.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "1.1.0"
    },
    "host": "localhost:8080",
    "basePath": "/v1",
    "schemes": ["http"],
    "paths": {},
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "X-API-Key",
            "in": "header"
        }
    }
}`

func init() {
	swag.Register(swag.Name, &s{})
}

type s struct{}

func (s *s) ReadDoc() string {
	return docTemplate
}
