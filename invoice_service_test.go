package main

import "testing"

func TestInvoiceHTMLTotals(t *testing.T) {
	cases := []struct {
		name  string
		order Order
		total int
		want  string
	}{
		{"two lines", Order{ID: "o-1", Seller: "Shop", Buyer: "Lee", Items: []Item{{Name: "Pen", Quantity: 2, UnitCents: 150}, {Name: "Book", Quantity: 1, UnitCents: 500}}}, 800, "Total: $8.00"},
		{"reject empty", Order{ID: "o-2", Seller: "Shop", Buyer: "Lee"}, 0, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, total, err := invoiceHTML(tc.order)
			if tc.want == "" {
				if err == nil {
					t.Fatal("expected validation error")
				}
				return
			}
			if err != nil || total != tc.total {
				t.Fatalf("total=%d err=%v", total, err)
			}
			if !contains(got, tc.want) {
				t.Fatalf("missing %q", tc.want)
			}
		})
	}
}
func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
