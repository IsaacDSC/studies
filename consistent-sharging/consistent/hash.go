package consistent

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

// HashRing representa o anel de hash consistente
type HashRing struct {
	mu       sync.RWMutex
	ring     map[uint32]string // hash -> node
	sorted   []uint32          // hashes ordenados para busca binária
	nodes    map[string]bool   // nodes reais
	replicas int               // número de réplicas virtuais por node
}

// New cria um novo HashRing com o número especificado de réplicas virtuais
func New(replicas int) *HashRing {
	if replicas <= 0 {
		replicas = 100 // padrão recomendado para boa distribuição
	}
	return &HashRing{
		ring:     make(map[uint32]string),
		sorted:   make([]uint32, 0),
		nodes:    make(map[string]bool),
		replicas: replicas,
	}
}

// hashKey gera um hash uint32 para uma string
func hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// AddNode adiciona um node ao anel com suas réplicas virtuais
func (h *HashRing) AddNode(node string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.nodes[node] {
		return // node já existe
	}

	h.nodes[node] = true

	// Adiciona réplicas virtuais
	for i := 0; i < h.replicas; i++ {
		virtualKey := fmt.Sprintf("%s#%d", node, i)
		hash := hashKey(virtualKey)
		h.ring[hash] = node
		h.sorted = append(h.sorted, hash)
	}

	// Mantém a lista ordenada
	sort.Slice(h.sorted, func(i, j int) bool {
		return h.sorted[i] < h.sorted[j]
	})
}

// RemoveNode remove um node do anel
func (h *HashRing) RemoveNode(node string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.nodes[node] {
		return // node não existe
	}

	delete(h.nodes, node)

	// Remove todas as réplicas virtuais
	for i := 0; i < h.replicas; i++ {
		virtualKey := fmt.Sprintf("%s#%d", node, i)
		hash := hashKey(virtualKey)
		delete(h.ring, hash)
	}

	// Reconstrói a lista ordenada
	h.sorted = make([]uint32, 0, len(h.ring))
	for hash := range h.ring {
		h.sorted = append(h.sorted, hash)
	}
	sort.Slice(h.sorted, func(i, j int) bool {
		return h.sorted[i] < h.sorted[j]
	})
}

// GetNode retorna o node responsável pela chave dada
func (h *HashRing) GetNode(key string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.sorted) == 0 {
		return ""
	}

	hash := hashKey(key)

	// Busca binária para encontrar o primeiro hash >= key hash
	idx := sort.Search(len(h.sorted), func(i int) bool {
		return h.sorted[i] >= hash
	})

	// Se passou do fim, volta para o início (anel circular)
	if idx >= len(h.sorted) {
		idx = 0
	}

	return h.ring[h.sorted[idx]]
}

// GetNodes retorna uma lista de todos os nodes
func (h *HashRing) GetNodes() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	nodes := make([]string, 0, len(h.nodes))
	for node := range h.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetNodeCount retorna o número de nodes no anel
func (h *HashRing) GetNodeCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.nodes)
}

// GetNNodes retorna os N nodes mais próximos para uma chave (útil para replicação)
func (h *HashRing) GetNNodes(key string, n int) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.sorted) == 0 || n <= 0 {
		return nil
	}

	if n > len(h.nodes) {
		n = len(h.nodes)
	}

	hash := hashKey(key)
	idx := sort.Search(len(h.sorted), func(i int) bool {
		return h.sorted[i] >= hash
	})

	if idx >= len(h.sorted) {
		idx = 0
	}

	result := make([]string, 0, n)
	seen := make(map[string]bool)

	for len(result) < n {
		node := h.ring[h.sorted[idx]]
		if !seen[node] {
			seen[node] = true
			result = append(result, node)
		}
		idx = (idx + 1) % len(h.sorted)
	}

	return result
}

