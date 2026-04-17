# Especificacao - Biblioteca de Validacao sem Reflection (Collect-All)

## Objetivo

Definir uma abordagem de validacao para APIs em Go baseada em tipos primitivos e regras explicitas, sem uso de reflection/tag parsing, com foco em:

- Performance previsivel
- Legibilidade e manutencao
- Reuso de regras de validacao
- Retorno de multiplos erros por payload (padrao collect-all)

---

## Motivacao

No estado atual, validacoes com `if/else` espalhadas funcionam, mas tendem a:

- Duplicar regras entre handlers/use cases
- Misturar validacao com logica de negocio
- Dificultar padronizacao de erros para resposta HTTP

A proposta e centralizar regras em uma biblioteca pequena e explicita, sem custo de reflection.

---

## Principios de Design

- **Sem reflection**: somente funcoes e tipos explicitamente chamados
- **Composicao de regras**: funcoes pequenas, puras e reutilizaveis
- **Collect-all por padrao**: acumular erros por campo antes de retornar
- **Tipagem forte opcional**: wrappers de dominio (`Name`, `Age`, etc.) para garantir invariantes
- **Erros estruturados**: payload consistente para API e observabilidade

---

## Estrutura de Pacotes Sugerida

```text
fields-validator/
  validate/
    errors.go
    string_rules.go
    number_rules.go
    collect.go
  types/
    name.go
    age.go
  example/
    main.go
  specs/
    validation-collect-all.md
```

---

## Modelo de Erro (Estruturado)

```go
package validate

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

type Errors []FieldError

func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}
	return "validation failed"
}

func (e Errors) HasAny() bool {
	return len(e) > 0
}
```

### Convencao recomendada para `Code`

- `required`
- `min`
- `max`
- `range`
- `one_of`
- `invalid_format`

---

## Contrato das Regras

Cada regra recebe `field` e `value` e retorna erro estruturado (ou vazio).

```go
package validate

type StringRule func(field, value string) []FieldError
type IntRule func(field string, value int) []FieldError
```

Essa assinatura facilita composicao e elimina alocacoes desnecessarias de wrappers complexos.

---

## Regras Genericas - String

```go
package validate

import "strings"

func Required() StringRule {
	return func(field, value string) []FieldError {
		if strings.TrimSpace(value) == "" {
			return []FieldError{{
				Field: field, Code: "required", Message: "field is required", Value: value,
			}}
		}
		return nil
	}
}

func MinLen(min int) StringRule {
	return func(field, value string) []FieldError {
		if len(value) < min {
			return []FieldError{{
				Field: field, Code: "min", Message: "length is below minimum", Value: value,
			}}
		}
		return nil
	}
}

func MaxLen(max int) StringRule {
	return func(field, value string) []FieldError {
		if len(value) > max {
			return []FieldError{{
				Field: field, Code: "max", Message: "length exceeds maximum", Value: value,
			}}
		}
		return nil
	}
}
```

---

## Regras Genericas - Numero

```go
package validate

func Min(min int) IntRule {
	return func(field string, value int) []FieldError {
		if value < min {
			return []FieldError{{
				Field: field, Code: "min", Message: "value is below minimum", Value: value,
			}}
		}
		return nil
	}
}

func Max(max int) IntRule {
	return func(field string, value int) []FieldError {
		if value > max {
			return []FieldError{{
				Field: field, Code: "max", Message: "value exceeds maximum", Value: value,
			}}
		}
		return nil
	}
}

func Range(min, max int) IntRule {
	return func(field string, value int) []FieldError {
		if value < min || value > max {
			return []FieldError{{
				Field: field, Code: "range", Message: "value out of accepted range", Value: value,
			}}
		}
		return nil
	}
}
```

---

## Motor de Collect-All

```go
package validate

func ApplyStringRules(field, value string, rules ...StringRule) []FieldError {
	var errs []FieldError
	for _, rule := range rules {
		errs = append(errs, rule(field, value)...)
	}
	return errs
}

func ApplyIntRules(field string, value int, rules ...IntRule) []FieldError {
	var errs []FieldError
	for _, rule := range rules {
		errs = append(errs, rule(field, value)...)
	}
	return errs
}
```

Esse motor e simples de testar, barato em CPU e sem mecanismo dinamico oculto.

---

## Exemplo de DTO da API com Collect-All

```go
package input

import "fields-validator/validate"

type CreateUserInput struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (in CreateUserInput) Validate() validate.Errors {
	var errs validate.Errors

	errs = append(errs, validate.ApplyStringRules(
		"name",
		in.Name,
		validate.Required(),
		validate.MinLen(3),
		validate.MaxLen(80),
	)...)

	errs = append(errs, validate.ApplyIntRules(
		"age",
		in.Age,
		validate.Range(18, 120),
	)...)

	return errs
}
```

Uso no handler:

```go
func CreateUserHandler(in CreateUserInput) error {
	if errs := in.Validate(); errs.HasAny() {
		// retornar 422 com payload padronizado
		return errs
	}

	// fluxo normal
	return nil
}
```

---

## Exemplo de Tipos de Dominio (Opcional)

Essa etapa reduz dados invalidos circulando no dominio.

```go
package types

import "fields-validator/validate"

type Name string

func NewName(v string) (Name, validate.Errors) {
	errs := validate.ApplyStringRules(
		"name",
		v,
		validate.Required(),
		validate.MinLen(3),
		validate.MaxLen(80),
	)
	if len(errs) > 0 {
		return "", errs
	}
	return Name(v), nil
}
```

Esse padrao permite separar claramente:

- validacao de borda (DTO/API)
- construcao segura de objetos de dominio

---

## Comparacao Rapida de Abordagens

### 1) `if/else` manual em cada handler

- Pro: simples no curto prazo
- Contra: duplicacao e falta de padrao

### 2) Lib com tags e reflection

- Pro: menos codigo inicial
- Contra: custo de reflection + menor controle fino

### 3) Regras explicitas + collect-all (proposta)

- Pro: sem reflection, previsivel, testavel, padronizado
- Contra: exige definicao inicial de regras e convencoes

---

## Abordagem Visual no Go: Comentario + Codegen

Para manter uma experiencia declarativa dentro do proprio arquivo Go (sem reflection no runtime), usar anotacoes em comentario e gerar codigo estatico.

### Exemplo de DTO anotado

```go
package input

//go:generate go run ./cmd/fields-validator-gen -type InputDTO -in input_dto.go -out input_dto_validated.gen.go
type InputDTO struct {
	Name  string `json:"name"`  // @validate nonEmpty
	Age   int    `json:"age"`   // @validate min=18,max=120
	Email string `json:"email"` // @validate email
}
```

### O que o generator deve produzir

- Tipos validados para os campos (opcional, conforme configuracao):
  - `type InputDTOName string`
  - `type InputDTOAge int`
  - `type InputDTOEmail string`
- Funcoes de validacao por campo, sem reflection:
  - `validateInputDTOName(field string, value string) []validate.FieldError`
  - `validateInputDTOAge(field string, value int) []validate.FieldError`
- Metodo collect-all no DTO:
  - `func (in InputDTO) ValidateCollect() validate.Errors`
- Opcional: metodo de decode com collect-all:
  - `func DecodeInputDTO(data []byte) (InputDTO, validate.Errors)`

### Exemplo de arquivo gerado (resumo)

```go
package input

import "fields-validator/validate"

func (in InputDTO) ValidateCollect() validate.Errors {
	var errs validate.Errors

	errs = append(errs, validate.ApplyStringRules(
		"name",
		in.Name,
		validate.Required(),
		validate.MinLen(1),
	)...)

	errs = append(errs, validate.ApplyIntRules(
		"age",
		in.Age,
		validate.Min(18),
		validate.Max(120),
	)...)

	errs = append(errs, validate.ApplyStringRules(
		"email",
		in.Email,
		validate.Required(),
		validate.Email(),
	)...)

	return errs
}
```

### Pipeline de geracao recomendado

1. Desenvolvedor adiciona anotacoes `// @validate ...` no DTO.
2. `go generate ./...` executa `fields-validator-gen`.
3. Generator parseia AST/comentarios e gera `*.gen.go`.
4. Build/test executam normalmente sem reflection custom.

### Regras de parsing da anotacao

- Prefixo obrigatorio: `@validate`
- Separador de regras: virgula
- Formatos aceitos:
  - `nonEmpty`
  - `email`
  - `min=18`
  - `max=120`
  - `oneOf=admin|user|viewer`
- Erro de geracao se:
  - regra desconhecida
  - parametro invalido
  - combinacao incoerente (ex.: `min=10,max=2`)

### Vantagens desta abordagem

- Mantem o "visual de schema" proximo do DTO
- Nao depende de reflection/tag parser no hot path
- Codigo final permanece idiomatico e debuggavel em Go
- Facilita padronizacao em times grandes

### Limitacoes conhecidas

- Sintaxe exata `String[min(...)]` nao e valida em Go
- Exige etapa de geracao (`go generate`) no fluxo do projeto
- Mudanca de regra requer regerar arquivos `*.gen.go`

---

## Plano de Adocao Incremental

1. Definir `FieldError`, `Errors` e codigos padrao.
2. Implementar 5 regras base:
   - string: `Required`, `MinLen`, `MaxLen`
   - numero: `Min`, `Max` ou `Range`
3. Criar MVP do generator com `@validate nonEmpty|min|max|email`.
4. Adotar em 1 endpoint critico com collect-all gerado.
5. Padronizar resposta HTTP de validacao (`422`).
6. Gradualmente migrar demais endpoints.

---

## Testes Recomendados

- Testes de unidade por regra (table-driven tests)
- Testes do agregador collect-all garantindo:
  - acumulo correto de erros
  - ordem consistente (se relevante para contrato)
- Benchmark simples comparando:
  - validacao manual atual
  - regras compositas propostas

Exemplo de benchmark:

```go
func BenchmarkValidateCreateUserInput(b *testing.B) {
	in := CreateUserInput{Name: "ab", Age: 10}
	for i := 0; i < b.N; i++ {
		_ = in.Validate()
	}
}
```

---

## Riscos e Mitigacoes

- **Risco**: proliferacao de regras muito especificas
  - **Mitigacao**: manter regras genericas em `validate` e regras de negocio em `types`
- **Risco**: mensagens de erro inconsistentes
  - **Mitigacao**: catalogo de `Code` e testes de contrato de resposta
- **Risco**: regressao em migracao gradual
  - **Mitigacao**: migrar endpoint por endpoint com testes de integracao

---

## Conclusao

O padrao collect-all com regras explicitas e tipagem forte opcional atende bem o objetivo de evitar reflection, manter performance previsivel e melhorar padronizacao de validacao na API. A estrategia incremental reduz risco e permite ganho rapido em consistencia sem refactor grande de uma vez.
