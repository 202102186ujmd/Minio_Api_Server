package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
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
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "schemes": ["http"],
    "paths": {}
}`

// SwaggerInfo guarda metadata editable en runtime.
var SwaggerInfo = &swag.Spec{
	Version:          "1.1.0",
	Host:             "localhost:8080",
	BasePath:         "/v1",
	Title:            "MinIO Management API",
	Description:      "API profesional para gestión de MinIO (CRUD de objetos), con validaciones, logging, métricas y respuestas estandarizadas.",
	InfoInstanceName: swag.Name,
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
