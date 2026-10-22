package main

import (
	"context"
	"fmt"
	"html"
	"strings"
)

type Order struct {
	ID     string `json:"id"`
	Seller string `json:"seller"`
	Buyer  string `json:"buyer"`
	Email  string `json:"email"`
	Items  []Item `json:"items"`
}
type Item struct {
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	UnitCents int    `json:"unit_cents"`
}

func invoiceHTML(order Order) (string, int, error) {
	if order.ID == "" || order.Seller == "" || order.Buyer == "" || len(order.Items) == 0 {
		return "", 0, fmt.Errorf("order needs id, seller, buyer, and items")
	}
	total := 0
	var rows strings.Builder
	for _, item := range order.Items {
		if item.Quantity <= 0 || item.UnitCents < 0 {
			return "", 0, fmt.Errorf("invalid item %q", item.Name)
		}
		line := item.Quantity * item.UnitCents
		total += line
		fmt.Fprintf(&rows, "<tr><td>%s</td><td>%d</td><td>$%.2f</td><td>$%.2f</td></tr>", html.EscapeString(item.Name), item.Quantity, float64(item.UnitCents)/100, float64(line)/100)
	}
	return fmt.Sprintf("<!doctype html><html><body><h1>Invoice %s</h1><p>Seller: %s</p><p>Buyer: %s (%s)</p><table><tr><th>Item</th><th>Qty</th><th>Unit</th><th>Total</th></tr>%s</table><p>Total: $%.2f</p></body></html>", html.EscapeString(order.ID), html.EscapeString(order.Seller), html.EscapeString(order.Buyer), html.EscapeString(order.Email), rows.String(), float64(total)/100), total, nil
}

func GenerateInvoice(ctx context.Context, client *InfraiClient, order Order) (interface{}, error) {
	markup, _, err := invoiceHTML(order)
	if err != nil {
		return nil, err
	}
	return client.GeneratePDF(ctx, generateRequest{HTML: markup, PageSize: "A4", Orientation: "portrait", Store: true}, order.ID)
}
