package main

import (
	"fmt"

	"github.com/isaacdsc/consistent-sharding/consistent"
)

func main() {
	fmt.Println("=== Consistent Hashing Demo ===\n")

	// Cria um hash ring com 150 réplicas virtuais por node
	ring := consistent.New(150)

	// Simula servidores de cache
	servers := []string{
		"cache-server-1:6379",
		"cache-server-2:6379",
		"cache-server-3:6379",
	}

	// Adiciona os servidores ao ring
	fmt.Println("📡 Adicionando servidores ao ring:")
	for _, server := range servers {
		ring.AddNode(server)
		fmt.Printf("   + %s\n", server)
	}
	fmt.Println()

	// Demonstra o roteamento de chaves
	keys := []string{
		"user:1001",
		"user:1002",
		"user:1003",
		"session:abc123",
		"session:xyz789",
		"product:notebook",
		"product:mouse",
		"cart:user1001",
	}

	fmt.Println("🔑 Roteamento de chaves:")
	keyToServer := make(map[string]string)
	for _, key := range keys {
		server := ring.GetNode(key)
		keyToServer[key] = server
		fmt.Printf("   %s → %s\n", key, server)
	}
	fmt.Println()

	// Demonstra consistência: mesma chave sempre vai para o mesmo servidor
	fmt.Println("✅ Verificando consistência:")
	testKey := "user:1001"
	for i := 0; i < 5; i++ {
		server := ring.GetNode(testKey)
		fmt.Printf("   Tentativa %d: %s → %s\n", i+1, testKey, server)
	}
	fmt.Println()

	// Demonstra GetNNodes para replicação
	fmt.Println("📋 Replicação - obtendo 2 servidores para cada chave:")
	ring.AddNode("cache-server-4:6379")
	ring.AddNode("cache-server-5:6379")

	fmt.Println("✅ Verificando consistência:")
	testKey = "user:1001"
	for i := 0; i < 5; i++ {
		server := ring.GetNode(testKey)
		fmt.Printf("   Tentativa %d: %s → %s\n", i+1, testKey, server)
	}
	fmt.Println()

	// Estatísticas de distribuição
	fmt.Println("📊 Testando distribuição com 10.000 chaves:")
	distribution := make(map[string]int)
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("random-key-%d", i)
		server := ring.GetNode(key)
		distribution[server]++
	}

	for server, count := range distribution {
		bar := ""
		for j := 0; j < count/100; j++ {
			bar += "█"
		}
		fmt.Printf("   %s: %d %s\n", server, count, bar)
	}
}
