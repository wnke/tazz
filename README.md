# tazz

Spec-first HTTP driver for device CRUD using OpenAPI and generated server interfaces.

## API

The OpenAPI source of truth lives in `internal/driver/http/openapi.yaml`.

### Device fields

- `url` (URI)
- `username`
- `password` (write-only, never returned)
- `power` (`on`, `off`, `standby`)

## Code generation

This repository is configured for `oapi-codegen`.

```sh
just gen-api
```

Generated code target is `internal/driver/http/gen/*`.

## Tests

```sh
just test
```
