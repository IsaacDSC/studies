package consistent

import (
	"fmt"
	"testing"
)

func TestHashRing_Basic(t *testing.T) {
	ring := New(100)

	// Adiciona nodes
	ring.AddNode("server-1")
	ring.AddNode("server-2")
	ring.AddNode("server-3")

	if ring.GetNodeCount() != 3 {
		t.Errorf("expected 3 nodes, got %d", ring.GetNodeCount())
	}

	// Testa consistência: mesma chave sempre vai para o mesmo node
	key := "user:123"
	node1 := ring.GetNode(key)
	node2 := ring.GetNode(key)

	if node1 != node2 {
		t.Errorf("inconsistent: same key mapped to different nodes: %s vs %s", node1, node2)
	}
}

func TestHashRing_RemoveNode(t *testing.T) {
	ring := New(100)

	ring.AddNode("server-1")
	ring.AddNode("server-2")
	ring.AddNode("server-3")

	// Guarda o mapeamento de algumas chaves
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	originalMapping := make(map[string]string)
	for _, k := range keys {
		originalMapping[k] = ring.GetNode(k)
	}

	// Remove um node
	ring.RemoveNode("server-2")

	if ring.GetNodeCount() != 2 {
		t.Errorf("expected 2 nodes after removal, got %d", ring.GetNodeCount())
	}

	// Verifica que chaves que NÃO estavam em server-2 mantêm o mesmo node
	for _, k := range keys {
		oldNode := originalMapping[k]
		newNode := ring.GetNode(k)

		// Se estava em server-2, deve ter mudado
		// Se não estava, deve continuar no mesmo
		if oldNode != "server-2" && newNode != oldNode {
			t.Errorf("key %s moved from %s to %s after removing unrelated node", k, oldNode, newNode)
		}
	}
}

func TestHashRing_Distribution(t *testing.T) {
	ring := New(150)

	ring.AddNode("server-1")
	ring.AddNode("server-2")
	ring.AddNode("server-3")

	// Conta a distribuição de 10000 chaves
	distribution := make(map[string]int)
	numKeys := 10000

	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("key-%d", i)
		node := ring.GetNode(key)
		distribution[node]++
	}

	// Com 3 nodes, esperamos ~3333 chaves por node
	// Aceita uma variação de 20%
	expected := numKeys / 3
	tolerance := float64(expected) * 0.20

	for node, count := range distribution {
		diff := float64(count - expected)
		if diff < 0 {
			diff = -diff
		}
		if diff > tolerance {
			t.Errorf("node %s has %d keys, expected ~%d (tolerance: %.0f)", node, count, expected, tolerance)
		}
	}

	t.Logf("Distribution: %v", distribution)
}

func TestHashRing_GetNNodes(t *testing.T) {
	ring := New(100)

	ring.AddNode("server-1")
	ring.AddNode("server-2")
	ring.AddNode("server-3")

	nodes := ring.GetNNodes("mykey", 2)

	if len(nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(nodes))
	}

	// Verifica que não há duplicatas
	seen := make(map[string]bool)
	for _, n := range nodes {
		if seen[n] {
			t.Errorf("duplicate node in result: %s", n)
		}
		seen[n] = true
	}
}

func TestHashRing_EmptyRing(t *testing.T) {
	ring := New(100)

	node := ring.GetNode("somekey")
	if node != "" {
		t.Errorf("expected empty string for empty ring, got %s", node)
	}

	nodes := ring.GetNNodes("somekey", 3)
	if nodes != nil {
		t.Errorf("expected nil for empty ring, got %v", nodes)
	}
}

func BenchmarkGetNode(b *testing.B) {
	ring := New(150)

	for i := 0; i < 100; i++ {
		ring.AddNode(fmt.Sprintf("server-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ring.GetNode(fmt.Sprintf("key-%d", i))
	}
}

