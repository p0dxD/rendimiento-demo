# rendimiento-demo

Una página pequeña para probar [rendimiento.ai](https://github.com/p0dxD/rendimiento.ai) desde cero.

- `GET /`: la página (cambie el saludo con `SALUDO`)
- `GET /healthz`: la comprobación de salud
- `GET /api/info`: lo mismo en JSON

No tiene Dockerfile ni `rendimiento.yaml` a propósito: rendimiento detecta Go y los propone en la solicitud de incorporación.

```bash
go test ./...
go run .   # http://localhost:8080
```
