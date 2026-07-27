# Runter Backend API - Go + Docker

Backend profissional em Go com arquitetura em camadas, suporte Docker/Docker-Compose, banco de dados PostgreSQL com GORM, autenticação via JWT com hash bcrypt e versionamento de API (`/api/v1`).

---

## 🛠️ Tecnologias e Bibliotecas Utilizadas

- **Linguagem**: Go 1.22+
- **Framework/Router**: [Chi v5](github.com/go-chi/chi/v5) (Leve, 100% net/http compatível, suporte nativo a sub-roteamento e middleware)
- **Banco de Dados**: PostgreSQL 16
- **ORM**: [GORM](gorm.io/gorm) com driver Postgres e Auto-Migrate
- **Autenticação**: JWT (`golang-jwt/jwt/v5`) & Criptografia de senha com `bcrypt`
- **Orquestração**: Docker & Docker Compose (Multi-stage build com imagem final Alpine ultra leve)
- **Middlewares**: Logger, Recoverer, CORS, Auth JWT

---

## 📂 Arquitetura do Projeto (Standard Go Layout)

```text
backend/
├── cmd/
│   └── api/
│       └── main.go           # Ponto de entrada do servidor e roteamento v1
├── internal/
│   ├── config/               # Leitura de variáveis de ambiente (.env)
│   ├── database/             # Conexão com PostgreSQL & Migrações GORM
│   ├── domain/               # Entidades, DTOs e Interfaces
│   ├── handler/              # Handlers HTTP (Controllers)
│   ├── middleware/           # Middleware JWT de autorização
│   ├── repository/           # Camada de persistência / Acesso ao banco
│   └── service/              # Regras de negócio, geração de JWT e Bcrypt
├── .env / .env.example       # Variáveis de ambiente
├── Dockerfile                # Build multi-etapas (Go builder + Alpine runner)
├── docker-compose.yml        # Serviços App Go e Banco PostgreSQL
├── go.mod                    # Módulo Go e dependências
└── README.md                 # Documentação do projeto
```

---

## 🚀 Como Executar com Docker

### 1. Iniciar os Containers (API + PostgreSQL)

```bash
docker compose up -d --build
```

O Docker Compose irá:
1. Subir o container do PostgreSQL (`runter_postgres`) na porta `5432`.
2. Aguardar a verificação de saúde (`healthcheck`) do banco de dados.
3. Compilar a imagem Go e iniciar a API (`runter_api`) na porta `8080`.

### 2. Verificar Logs

```bash
docker compose logs -f app
```

### 3. Encerrar os Containers

```bash
docker compose down -v
```

---

## 📌 Endpoints da API (`/api/v1`)

### 1. Verificação de Saúde (Público)
- **GET** `/api/v1/health`
  - Resposta: `200 OK`

### 2. Registro de Usuário (Público)
- **POST** `/api/v1/auth/register`
  - Body:
    ```json
    {
      "name": "Desenvolvedor Go",
      "email": "dev@exemplo.com",
      "password": "senha_segura_123"
    }
    ```

### 3. Login de Usuário (Público)
- **POST** `/api/v1/auth/login`
  - Body:
    ```json
    {
      "email": "dev@exemplo.com",
      "password": "senha_segura_123"
    }
    ```
  - Resposta:
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": 1,
        "name": "Desenvolvedor Go",
        "email": "dev@exemplo.com",
        "created_at": "2026-07-24T22:00:00Z",
        "updated_at": "2026-07-24T22:00:00Z"
      }
    }
    ```

### 4. Perfil do Usuário Autenticado (Protegido por JWT)
- **GET** `/api/v1/users/me`
  - Header: `Authorization: Bearer <seu_token_jwt>`
  - Resposta:
    ```json
    {
      "id": 1,
      "name": "Desenvolvedor Go",
      "email": "dev@exemplo.com",
      "created_at": "2026-07-24T22:00:00Z"
    }
    ```
