# Consistent Hashing em Go

Uma implementação simples e eficiente de **Consistent Hashing** em Go, útil para distribuição de carga em sistemas distribuídos.

## O que é Consistent Hashing?

Consistent Hashing é uma técnica que permite distribuir dados entre servidores de forma que, quando um servidor é adicionado ou removido, apenas uma fração mínima das chaves precisa ser remapeada.

### Benefícios

- **Mínima redistribuição**: Ao adicionar/remover nodes, apenas ~1/n das chaves são movidas
- **Escalabilidade horizontal**: Fácil adicionar ou remover servidores
- **Balanceamento de carga**: Réplicas virtuais garantem distribuição uniforme

## Uso

```go
package main

import (
    "fmt"
    "github.com/isaacdsc/consistent-sharding/consistent"
)

func main() {
    // Cria um ring com 150 réplicas virtuais por node
    ring := consistent.New(150)

    // Adiciona servidores
    ring.AddNode("server-1:6379")
    ring.AddNode("server-2:6379")
    ring.AddNode("server-3:6379")

    // Encontra o servidor para uma chave
    server := ring.GetNode("user:123")
    fmt.Printf("user:123 → %s\n", server)

    // Obtém múltiplos servidores para replicação
    servers := ring.GetNNodes("user:123", 2)
    fmt.Printf("Réplicas: %v\n", servers)

    // Remove um servidor (simula falha)
    ring.RemoveNode("server-2:6379")
}
```

## API

### `New(replicas int) *HashRing`
Cria um novo hash ring. `replicas` define o número de réplicas virtuais por node (recomendado: 100-200).

### `AddNode(node string)`
Adiciona um node ao ring.

### `RemoveNode(node string)`
Remove um node do ring.

### `GetNode(key string) string`
Retorna o node responsável pela chave.

### `GetNNodes(key string, n int) []string`
Retorna os N nodes mais próximos (útil para replicação).

### `GetNodes() []string`
Lista todos os nodes.

### `GetNodeCount() int`
Retorna o número de nodes.

## Executando

```bash
# Rodar o exemplo
go run main.go

# Rodar os testes
go test ./consistent -v

# Benchmark
go test ./consistent -bench=.
```

## Como funciona

1. **Hash Ring**: Os nodes são mapeados em um anel circular usando seus hashes
2. **Réplicas Virtuais**: Cada node físico tem múltiplas posições no anel para melhor distribuição
3. **Busca de Chave**: Para encontrar o node de uma chave, calculamos seu hash e encontramos o próximo node no sentido horário

```
        Node A          Node B
           ↓               ↓
    ┌──────●───────────────●──────┐
    │                             │
    ●                             │
    ↑                             │
  Node C                          │
    │                             ●
    │                             ↑
    │                           Node A'
    │                          (réplica)
    └─────────────●───────────────┘
                  ↑
                Node B'
               (réplica)
```

## Casos de Uso

- **Cache distribuído** (Redis Cluster, Memcached)
- **Sharding de banco de dados**
- **Load balancing** com afinidade de sessão
- **CDN** para roteamento de conteúdo





Adicionar um novo nó não teriamos problemas, somente teria problema caso seja removido um nó, sendo necessário realizar uma migração de dados.