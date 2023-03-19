package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Pasta",
		Price: 1.00,
		SKU:   "abc-def-ghi",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
