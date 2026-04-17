package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	payload := []byte(`{"name":"Joao","age":16,"email":"invalid-email","amount":2.4,"status":"PENDENTE"}`)

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
