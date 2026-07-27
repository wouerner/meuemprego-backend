FROM golang:1.23-alpine

# Instalar dependências do sistema
RUN apk add --no-cache git ca-certificates

# Instalar Air (Hot Reload) e Swag (Gerador de Swagger/OpenAPI) em versões compatíveis
RUN go install github.com/air-verse/air@v1.52.3
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

WORKDIR /app

# Copiar arquivos do módulo Go
COPY go.mod go.sum* ./

# Garantir dependências
RUN go mod download

# Copiar o restante do código-fonte
COPY . .

# Expõe a porta da API
EXPOSE 8080

# Gerar documentação Swagger, atualizar dependências e iniciar o Air com Hot Reload
CMD ["sh", "-c", "go mod tidy && swag init -g cmd/api/main.go -o ./docs && air -c .air.toml"]
