# Gootel - Go OpenTelemetry Demo

Uma aplicação de demonstração em Go com instrumentação completa usando OpenTelemetry para traces e métricas.

## Arquitetura

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────┐
│   Go App    │────▶│  OTel Collector  │────▶│   Jaeger    │
│  (Traces &  │     │                  │     │  (Traces)   │
│   Metrics)  │     │                  │     └─────────────┘
└─────────────┘     │                  │     ┌─────────────┐
                    │                  │────▶│ Prometheus  │
                    └──────────────────┘     │  (Metrics)  │
                                             └─────────────┘
                                                    │
                                             ┌──────▼──────┐
                                             │   Grafana   │
                                             │(Dashboards) │
                                             └─────────────┘
```

## Componentes

- **Go Application**: Servidor HTTP com endpoints instrumentados
- **OpenTelemetry Collector**: Recebe, processa e exporta telemetria
- **Jaeger**: Armazenamento e visualização de traces distribuídos
- **Prometheus**: Coleta e armazenamento de métricas
- **Grafana**: Dashboards e visualização (opcional)

## Quick Start

### Pré-requisitos

- Docker e Docker Compose instalados
- Go 1.21+ (apenas para desenvolvimento local)

### Executar com Docker Compose

```bash
# Iniciar todos os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f app

# Parar todos os serviços
docker-compose down
```

### Executar localmente (desenvolvimento)

```bash
# Instalar dependências
go mod tidy

# Iniciar infraestrutura
docker-compose up -d otel-collector jaeger prometheus grafana

# Executar a aplicação
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
go run ./cmd/server
```

## Endpoints da Aplicação

| Endpoint | Descrição |
|----------|-----------|
| `GET /` | Página inicial |
| `GET /health` | Health check |
| `GET /api/users` | Lista usuários (simula query no DB) |
| `GET /api/process` | Processa dados (demonstra spans aninhados) |

## UIs de Observabilidade

| Serviço | URL | Credenciais |
|---------|-----|-------------|
| **Aplicação** | http://localhost:8080 | - |
| **Jaeger UI** | http://localhost:16686 | - |
| **Prometheus** | http://localhost:9090 | - |
| **Grafana** | http://localhost:3000 | admin/admin |

## Testando a Instrumentação

### Gerar traces

```bash
# Requisição simples
curl http://localhost:8080/

# Buscar usuários (gera span de "database")
curl http://localhost:8080/api/users

# Processar dados (gera múltiplos spans aninhados)
curl http://localhost:8080/api/process

# Gerar carga para métricas
for i in {1..100}; do
  curl -s http://localhost:8080/api/users > /dev/null
  curl -s http://localhost:8080/api/process > /dev/null
  sleep 0.1
done
```

### Visualizar no Jaeger

1. Acesse http://localhost:16686
2. Selecione o serviço `gootel-service`
3. Clique em "Find Traces"
4. Explore os traces e spans

### Visualizar no Prometheus

1. Acesse http://localhost:9090
2. Consulte métricas como:
   - `http_requests_total`
   - `http_request_duration_seconds_bucket`

### Visualizar no Grafana

1. Acesse http://localhost:3000
2. Login: admin/admin
3. As datasources Prometheus e Jaeger já estão configuradas

## Estrutura do Projeto

```
gootel/
├── cmd/
│   └── server/
│       └── main.go          # Aplicação principal
├── internal/
│   └── otel/
│       └── otel.go          # Configuração OpenTelemetry
├── config/
│   ├── otel-collector-config.yaml
│   ├── prometheus.yaml
│   └── grafana/
│       └── provisioning/
├── docker-compose.yaml
├── Dockerfile
├── go.mod
└── README.md
```

## Métricas Coletadas

| Métrica | Tipo | Descrição |
|---------|------|-----------|
| `http_requests_total` | Counter | Total de requisições HTTP |
| `http_request_duration_seconds` | Histogram | Latência das requisições |

## Traces

Cada requisição HTTP gera um trace com spans para:
- Recebimento da requisição HTTP
- Operações de "banco de dados"
- Processamento de dados (validação, transformação, armazenamento)

## Variáveis de Ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `SERVICE_NAME` | gootel-service | Nome do serviço |
| `SERVICE_VERSION` | 1.0.0 | Versão do serviço |
| `ENVIRONMENT` | development | Ambiente |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | localhost:4317 | Endpoint do Collector |

## Desenvolvimento

### Adicionar novas métricas

```go
// No main.go ou em um arquivo de métricas
counter, _ := meter.Int64Counter("my_custom_counter",
    metric.WithDescription("Minha métrica customizada"),
)

// Usar
counter.Add(ctx, 1, metric.WithAttributes(
    attribute.String("key", "value"),
))
```

### Adicionar novos spans

```go
ctx, span := tracer.Start(ctx, "operation-name",
    trace.WithAttributes(
        attribute.String("key", "value"),
    ),
)
defer span.End()

// Seu código aqui
span.SetAttributes(attribute.Bool("success", true))
```

## Troubleshooting

### Traces não aparecem no Jaeger

1. Verifique se o collector está rodando: `docker-compose logs otel-collector`
2. Verifique a conectividade: `docker-compose exec app nc -zv otel-collector 4317`

### Métricas não aparecem no Prometheus

1. Verifique o exporter do collector: http://localhost:8889/metrics
2. Verifique os targets do Prometheus: http://localhost:9090/targets

## Licença

MIT
