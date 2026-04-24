package main

//go:generate go run ../cmd/fields-validator-gen -type InputDTO -in input_dto.go -out input_dto_validated.gen.go
type InputDTO struct {
	Name     string  `json:"name"`     // @validate nonEmpty,contains("Isaac")
	Age      int     `json:"age"`      // @validate min=18,max=120,equals(18),contains(1)
	Email    string  `json:"email"`    // @validate nonEmpty,email
	Amount   float64 `json:"amount"`   // @validate equals(1.8),contains(1.8)
	Amount2  float32 `json:"amount2"`  // @validate equals(1.8),contains(1.8)
	Status   string  `json:"status"`   // @validate possibilities("PAGAMENTO"||"PAGADO"||"QUITADO")
	Username *string `json:"username"` // @validate notNil,nonEmpty,minLen(3),maxLen(20)
	Age2     int32   `json:"age2"`     // @validate min=18,max=120,equals(18),contains(1)
	Age3     int64   `json:"age3"`     // @validate min=18,max=120,equals(18),contains(1)
	Date     string  `json:"date"`     // @validate date(YYYY-MM-DD)
	Time     string  `json:"time"`     // @validate time(HH:MM:SS)
	DateTime string  `json:"dateTime"` // @validate dateTime(YYYY-MM-DD HH:MM:SS)
}
