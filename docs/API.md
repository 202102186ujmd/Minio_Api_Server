# MinIO API

Este proyecto expone una API REST profesional para administrar buckets y objetos en MinIO con:

- Gin (alto rendimiento)
- Validaciones robustas
- Logging estructurado (zap)
- Respuestas estandarizadas
- API Key
- Swagger UI
- Métricas Prometheus

## Swagger

Para generar endpoints en Swagger debes ejecutar:

```bash
swag init -g cmd/api/main.go -o docs
```

Luego inicia la API y visita `http://localhost:8080/swagger/index.html`.

## Endpoints

- `GET /health`
- `GET /health/live`
- `GET /health/ready`
- `GET /metrics`
- `GET /v1/buckets`
- `GET /v1/buckets/{bucket}/objects`
- `POST /v1/buckets/{bucket}/objects`
- `GET /v1/buckets/{bucket}/objects/{object}/download`
- `GET /v1/buckets/{bucket}/objects/{object}/presign`
- `POST /v1/buckets/{bucket}/objects/{object}/copy`
- `POST /v1/buckets/{bucket}/objects/{object}/move`
- `DELETE /v1/buckets/{bucket}/objects/{object}`

## Métricas

Prometheus expone métricas en `GET /metrics` (sin API key).

## PM2 (producción)

1. `make build-linux`
2. `pm2 start ecosystem.config.js`

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
