#!/bin/bash

# Генерация кода из OpenAPI спецификации

# Генерация серверного кода для gorilla/mux
echo "Generating server code..."
oapi-codegen -generate gorilla-server -package api ./api/openapi.yaml > ./internal/server/http/api/server.gen.go

# Генерация моделей
echo "Generating types..."
oapi-codegen -generate types -package api ./api/openapi.yaml > ./internal/server/http/api/types.gen.go

# Генерация клиента
echo "Generating client..."
oapi-codegen -generate client -package api ./api/openapi.yaml > ./internal/server/http/api/client.gen.go

# Генерация клиента
echo "Generating spec..."
oapi-codegen -generate spec -package api ./api/openapi.yaml > ./internal/server/http/api/spec.gen.go

echo "Код успешно сгенерирован из OpenAPI спецификации"