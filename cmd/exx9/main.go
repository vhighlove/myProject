package main

import (
	"exx9/masker"
	"fmt"
)

func customMasker(_ interface{}) interface{} {
	return "abracadabra"
}
func main() {
	jsonData := `{
		"name": "John",
		"age": 30,
		"salary": 1234.56,
		"isEmployee": true,
		"address": {
			"city": "New York",
			"street": "5th Avenue"
		},
		"creditCard": {
			"number": "1234-5678-9101-1121",
			"cvv": "123"
		},
		"transactions": [
			{"id": 1, "amount": 100},
			{"id": 2, "amount": 150}
		]
	}`

	// Маскування за ключами
	maskedJSONByKey, err := masker.MaskByKeys(jsonData, []string{"name", "transactions", "cvv"}, customMasker)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Masked by Keys:", string(maskedJSONByKey))

	// Маскування за шляхами
	maskedJSONByPath, err := masker.MaskByIndexes(jsonData, []int{1, 5})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Masked by Paths:", string(maskedJSONByPath))
}
