# MinIO API

Este proyecto expone una API REST profesional para administrar buckets y objetos en MinIO con:

- Gin (alto rendimiento)
- Validaciones robustas
- Logging estructurado (zap)
- Respuestas estandarizadas
- API Key
- Swagger UI

## Endpoints

- `GET /health`
- `GET /v1/buckets`
- `GET /v1/buckets/{bucket}/objects`
- `POST /v1/buckets/{bucket}/objects`
- `DELETE /v1/buckets/{bucket}/objects/{object}`

## Swagger

1. Instalar swag: `go install github.com/swaggo/swag/cmd/swag@latest`
2. Generar docs: `make swagger`
3. Iniciar API: `make run`
4. Visitar: `http://localhost:8080/swagger/index.html`

## Ejemplo de respuesta estándar

```json
{
  "success": true,
  "code": "OK",
  "message": "Objetos listados",
  "data": [],
  "request_id": "a1b2c3"
}
```
