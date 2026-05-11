## Como Executar o Projeto

### Requisitos

- Docker
- Docker Compose

### Passos para rodar a aplicação

1. Clone o repositório.
2. Na raiz do projeto, execute:
    ```bash
    docker-compose up --build
    ```
3. A aplicação estará disponível em `http://localhost:8080`.

### Como Testar

Para rodar os testes unitários via Docker, execute:

```bash
docker-compose run app go test ./...
```

### Configuração

As configurações podem ser ajustadas no arquivo `.env` ou diretamente no `docker-compose.yaml`.

| Variável               | Descrição                                    | Padrão |
| ---------------------- | -------------------------------------------- | ------ |
| PORT                   | Porta do servidor                            | 8080   |
| DEFAULT_IP_LIMIT       | Limite de requisições por segundo para IP    | 10     |
| DEFAULT_TOKEN_LIMIT    | Limite de requisições por segundo para Token | 100    |
| BLOCK_DURATION_SECONDS | Tempo de bloqueio em segundos                | 300    |
| REDIS_HOST             | Host do Redis                                | redis  |
| REDIS_PORT             | Porta do Redis                               | 6379   |

### Estratégias de Limitação

O projeto utiliza o padrão **Strategy** para persistência. Atualmente, a implementação utiliza Redis (`internal/storage/redis.go`), mas pode ser facilmente estendida implementando a interface `PersistenceStrategy` em `internal/limiter/strategy.go`.

### Precedência

Se um header `API_KEY` for fornecido, o limite configurado para Tokens será aplicado, sobrepondo o limite por IP.
