package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gorules/zen-go"
)

func main() {
	engine := zen.NewEngine(zen.EngineConfig{})
	graph, _ := os.ReadFile("../jdm_graph.json")

	decision, err := engine.CreateDecision(graph)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	fmt.Println("=== Testando Regras ===")

	// Teste 1: Menor de idade
	testar(decision, "Menor (15 anos) + R$1000 VIP", map[string]any{
		"idade":       15,
		"valor":       1000,
		"tipoCliente": "vip",
	})

	// Teste 2: Maior de idade VIP
	testar(decision, "Maior (25 anos) + R$1000 VIP", map[string]any{
		"idade":       25,
		"valor":       1000,
		"tipoCliente": "vip",
	})

	// Teste 3: Maior de idade normal
	testar(decision, "Maior (30 anos) + R$600 Normal", map[string]any{
		"idade":       30,
		"valor":       600,
		"tipoCliente": "normal",
	})
}

func testar(decision zen.Decision, nome string, input map[string]any) {
	resultado, _ := decision.Evaluate(input)

	data, _ := json.Marshal(resultado)
	var resp map[string]any
	json.Unmarshal(data, &resp)

	result := resp["result"].(map[string]any)

	fmt.Printf("%s\n", nome)
	fmt.Printf("  Idade: %v → %v\n", result["statusIdade"], result["maiorIdade"])
	fmt.Printf("  Desconto: %v%% | Economia: R$%.2f | Final: R$%.2f\n\n",
		result["descontoPercent"],
		result["valorDesconto"],
		result["valorFinal"],
	)
}
