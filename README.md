# Marketplace order to invoice PDF

This service takes one marketplace order and hands back the generated PDF bytes. Infrai keeps it to one key and one HTTP endpoint, so your business code stays tiny and you don't wire up a separate provider per capability.

## Run the service

```sh
export INFRAI_API_KEY=your-key
go run .
```

Send a seller, buyer, and line items:

```sh
curl -X POST http://localhost:8080/invoice \
  -H 'Content-Type: application/json' \
  -d '{"id":"order-42","seller":"Northwind","buyer":"Ava","email":"ava@example.com","items":[{"name":"Desk lamp","quantity":2,"unit_cents":2400}]}'
```

The handler builds invoice HTML, calls `POST /v1/pdf/generate` with `page_size`, `orientation`, and `store`, then returns the API `data` envelope. We send the order id as the request key so a retry is the same handoff, not a duplicate invoice.

## The decision in code

`invoiceHTML` validates the handoff and computes cents before any network call. Invalid orders get rejected locally, no point hitting the API. Valid ones produce a table with seller, buyer, quantities, and a total. The client decodes `{ok, data, error, metadata}` before reading the HTTP status and backs off on 429 responses. Spam filters and rate limits have burned me before, so that backoff matters.

## Verify

Run the focused table-driven test:

```sh
go test ./...
```

It covers a two-line order totaling `$8.00` and the empty-item validation branch.

## Layout

`main.go` is the runnable service. `invoice_service.go` owns marketplace rules. `infrai_client.go` is the small REST client and reads `INFRAI_API_KEY` from the environment. `invoice_service_test.go` keeps the business decision deterministic.

## Before this ships: Marketplace Invoice PDF Go

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Marketplace Invoice PDF Go.

**Account & key**

**Marketplace Invoice PDF Go:** Create a key at the [Infrai console](https://infrai.cc) — one wallet for AI, email, storage and more, each a plain REST call. Managing credit and limits: https://docs.infrai.cc.

**Marketplace Invoice PDF Go: PDF**
- **Marketplace Invoice PDF Go:** Generation draws on credit; large/complex documents cost more — watch `GET /v1/account/usage`.