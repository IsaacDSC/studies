package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	corrects()
	fmt.Println("--------------------------------")
	wrongs()
}
func corrects() {
	// Valores alinhados a input_dto.go / ValidateCollect():
	// name contém "Isaac"; age == 18; email válido; amount == 1.8;
	// status ∈ {PAGAMENTO,PAGADO,QUITADO}; username presente, 3–20 caracteres.
	payload := []byte(`{
		"name": "Isaac Silva",
		"age": 18,
		"email": "isaac@example.com",
		"amount": 1.8,
		"amount2": 1.8,
		"status": "PAGAMENTO",
		"username": "isaac_user",
		"age2": 18,
		"age3": 18,
		"date": "2026-04-24",
		"time": "10:30:59",
		"dateTime": "2026-04-24 10:30:59"
	}`)

	var in InputDTO
	if err := json.Unmarshal(payload, &in); err != nil {
		fmt.Println("decode error:", err)
		return
	}

	if errs := in.ValidateCollect(); errs.HasAny() {
		fmt.Println("validation errors:", errs)
		for _, e := range errs {
			fmt.Printf("- field=%s code=%s message=%s value=%v\n", e.Field, e.Code, e.Message, e.Value)
		}
		return
	}

	fmt.Println("payload is valid")
	fmt.Printf("%+v\n", in)

}

func wrongs() {
	payload := []byte(`{
		"name": "Joao",
		"age": 16,
		"email": "invalid-email",
		"amount": 2.4,
		"amount2": 2.4,
		"status": "PENDENTE",
		"username": "",
		"age2": 16,
		"age3": 16,
		"date": "24/04/2026",
		"time": "10:30",
		"dateTime": "2026-04-24T10:30:59"
	}`)

	var in InputDTO
	if err := json.Unmarshal(payload, &in); err != nil {
		fmt.Println("decode error:", err)
		return
	}

	if errs := in.ValidateCollect(); errs.HasAny() {
		fmt.Println("validation errors:", errs)
		for _, e := range errs {
			fmt.Printf("- field=%s code=%s message=%s value=%v\n", e.Field, e.Code, e.Message, e.Value)
		}
		return
	}

	fmt.Println("payload is valid")
}
