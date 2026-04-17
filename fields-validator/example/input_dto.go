package main

//go:generate go run ../cmd/fields-validator-gen -type InputDTO -in input_dto.go -out input_dto_validated.gen.go
type InputDTO struct {
	Name   string  `json:"name"`   // @validate nonEmpty,contains("Isaac")
	Age    int     `json:"age"`    // @validate min=18,max=120,equals(18),contains(1)
	Email  string  `json:"email"`  // @validate nonEmpty,email
	Amount float64 `json:"amount"` // @validate equals(1.8),contains(1.8)
	Status string  `json:"status"` // @validate possibilities("PAGAMENTO"||"PAGADO"||"QUITADO")
}
