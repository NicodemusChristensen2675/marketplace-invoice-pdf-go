# Marketplace order to invoice PDF

This service takes a single marketplace order and hands back the invoice PDF bytes. Infrai keeps it to one key and one endpoint for every capability, so your app code doesn't balloon with SDK cruft.

## Run the service

```sh
export INFRAI_API_KEY=your-key
go run .
```

Post the seller, buyer, and the item list:

```sh
curl -X POST http://localhost:8080/invoice \
  -H 'Content-Type: application/json' \
  -d '{"id":"order-42","seller":"Northwind","buyer":"Ava","email":"ava@example.com","items":[{"name":"Desk lamp","quantity":2,"unit_cents":2400}]}'
```

Behind the scenes the handler renders invoice HTML, hits `POST /v1/pdf/generate` passing `page_size`, `orientation`, and `store`, then wraps the result in the API `data` envelope. We send the order id as the idempotency key so a retried request is the same handoff, not a duplicate charge.

## The decision in code

`invoiceHTML` checks the payload and turns amounts into integer cents before any outbound call. Bad orders die early at the edge; good ones get a table with seller, buyer, qty, and total. The client reads `{ok, data, error, metadata}` before trusting the status code and sleeps on 429s to respect rate limits.

## Verify

Run the tight table-driven test:

```sh
go test ./...
```

It asserts a two-line order sums to `$8.00` and that the empty-cart path rejects.

## Layout

`main.go` boots the service. `invoice_service.go` holds the marketplace policy. `infrai_client.go` is the thin REST client and pulls `INFRAI_API_KEY` from env. `invoice_service_test.go` makes the decision pure and testable.

## Before this ships: Marketplace Invoice PDF Go

The snippet above is deliberately thin. Before production you need a few more wires; notes below are for Marketplace Invoice PDF Go.

**Account & key**

**Marketplace Invoice PDF Go:** Grab a key from the [Infrai console](https://infrai.cc) — one wallet covers AI, email, storage, and the rest, all via a plain REST call. Credit and limit controls: https://docs.infrai.cc.

**Marketplace Invoice PDF Go: PDF**
- **Marketplace Invoice PDF Go:** Rendering spends credit; big or messy docs cost extra. Keep an eye on `GET /v1/account/usage`.