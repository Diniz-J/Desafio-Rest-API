# TODO API em Go

API REST para gerenciamento de tarefas (TODO list) com Go e MySQL.

## Tecnologias

- **Go 1.23+**
- **MySQL 8.0**
- **Gorilla Mux** (roteamento HTTP)
- **Docker & Docker Compose** (containerização)

## Arquitetura

Projeto organizado em camadas:

```
cmd/
  main.go                 - Entry point da aplicação
internal/
  handler/               - Camada HTTP (recebe requests, retorna responses)
  service/               - Lógica de negócio e validações
  repository/            - Acesso a dados (queries SQL)
  model/                 - Structs e constantes
  config/                - Configurações (database, etc)
```

## Requisitos

- Go 1.23 ou superior
- Docker e Docker Compose
- Git

## Como Rodar

### 1. Clone o repositório

```bash
git clone https://github.com/Diniz-J/Desafio-Rest-API.git
cd Desafio-Rest-API
```

### 2. Configure as variáveis de ambiente

```bash
cp .env.example .env
```

Edite o arquivo `.env` com suas credenciais:

```env
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=tasks_db
MYSQL_USER=tasks_user
MYSQL_PASSWORD=tasks_password

DB_HOST=main_db
DB_PORT=3306
DB_USER=tasks_user
DB_PASSWORD=tasks_password
DB_NAME=tasks_db
```

### 3. Suba o banco de dados

```bash
docker-compose up -d
```

Aguarde o container ficar saudável (healthcheck passa):

```bash
docker-compose ps
```

### 4. Crie a tabela de tasks

```bash
docker exec -i main_db mysql -uroot -prootpassword tasks_db < migrations/001_create_tasks_table.sql
```

### 5. Instale as dependências Go

```bash
go mod download
```

### 6. Rode a aplicação

```bash
go run cmd/main.go
```

A API estará disponível em: `http://localhost:8080`

## Endpoints

### POST /api/v1/tasks
Cria uma nova tarefa.

**Request:**
```json
{
  "title": "Comprar leite",
  "description": "2% gordura"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid-gerado",
  "title": "Comprar leite",
  "description": "2% gordura",
  "status": "pending",
  "priority": "medium",
  "created_at": "2026-02-03T10:00:00Z",
  "updated_at": "2026-02-03T10:00:00Z"
}
```

### GET /api/v1/tasks
Lista todas as tarefas. Aceita filtro por status.

**Query params:**
- `status` (opcional): `pending` ou `completed`

**Exemplos:**
```bash
# Listar todas
curl http://localhost:8080/api/v1/tasks

# Filtrar apenas pendentes
curl http://localhost:8080/api/v1/tasks?status=pending
```

**Response:** `200 OK`
```json
[
  {
    "id": "uuid-1",
    "title": "Comprar leite",
    "status": "pending",
    ...
  },
  {
    "id": "uuid-2",
    "title": "Estudar Go",
    "status": "completed",
    ...
  }
]
```

### GET /api/v1/tasks/{id}
Busca uma tarefa específica por ID.

**Response:** `200 OK` ou `404 Not Found`
```json
{
  "id": "uuid-1",
  "title": "Comprar leite",
  "description": "2% gordura",
  "status": "pending",
  "priority": "medium",
  "created_at": "2026-02-03T10:00:00Z",
  "updated_at": "2026-02-03T10:00:00Z"
}
```

### PUT /api/v1/tasks/{id}
Atualiza uma tarefa existente.

**Request:**
```json
{
  "title": "Comprar leite desnatado",
  "description": "0% gordura",
  "status": "completed",
  "priority": "high"
}
```

**Validações:**
- `status`: deve ser `pending` ou `completed`
- `priority`: deve ser `low`, `medium` ou `high`

**Response:** `200 OK` ou `404 Not Found`

### DELETE /api/v1/tasks/{id}
Deleta uma tarefa.

**Response:** `204 No Content` ou `404 Not Found`

### PATCH /api/v1/tasks/{id}/complete
Marca uma tarefa como concluída (atalho para não precisar enviar PUT completo).

**Response:** `200 OK` ou `404 Not Found`
```json
{
  "id": "uuid-1",
  "title": "Comprar leite",
  "status": "completed",
  ...
}
```

## Testes

Execute os testes unitários:

```bash
go test ./internal/service/... -v
```

Cobertura de testes:

```bash
go test ./internal/service/... -cover
```

## Docker

### Build da imagem

```bash
docker build -t todo-api .
```

### Rodar container

```bash
docker run -p 8080:8080 --env-file .env todo-api
```

## Estrutura do Banco de Dados

### Tabela: tasks

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | VARCHAR(36) PRIMARY KEY | UUID da tarefa |
| title | VARCHAR(255) NOT NULL | Título da tarefa |
| description | TEXT | Descrição detalhada (opcional) |
| status | ENUM('pending','completed') | Status atual |
| priority | ENUM('low','medium','high') | Prioridade |
| created_at | TIMESTAMP | Data de criação |
| updated_at | TIMESTAMP | Data da última atualização |
| deleted_at | TIMESTAMP NULL | Soft delete (não implementado) |

## Desenvolvimento

### Comandos úteis

```bash
# Rodar aplicação com hot reload (requer air)
air

# Ver logs do banco
docker-compose logs -f main_db

# Acessar banco diretamente
docker exec -it main_db mysql -uroot -p tasks_db

# Parar containers
docker-compose down

# Limpar volumes (CUIDADO: apaga dados)
docker-compose down -v
```

### Convenções de código

- Mensagens de commit seguem [Conventional Commits](https://www.conventionalcommits.org/)
- Código em inglês, documentação em português
- Testes unitários obrigatórios para camada de service
- Separação clara de responsabilidades entre camadas

## Git Flow

```
main           - Branch de produção (código estável)
develop        - Branch de integração
feature/*      - Features em desenvolvimento
bugfix/*       - Correções de bugs
```

## Licença

Projeto educacional - Mentoria Backend Go

## Autor

**Rodrigo Junior** ([@Diniz-J](https://github.com/Diniz-J))

Mentor: Andre Abreu ([@andreabreu76](https://github.com/andreabreu76))
