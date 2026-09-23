// SPDX-License-Identifier: MIT
// ADAPTED from Microsoft BCApps SalesQuotetoOrder.Codeunit.al at
// 2eae56d704a1fd035d104f333602aea7091b7749. See docs/provenance/BC_SALES_DERIVATION.md.
package bcsales

import "errors"

type DocumentType string

const (
	Quote DocumentType = "Quote"
	Order DocumentType = "Order"
)

var ErrQuoteType = errors.New("source document must be a quote")

// Header/Line represent only fields already present in the portable owner's
// quote/order contract. IDs are supplied by the caller's existing ID generator.
// No posting, VAT, prepayment, reservation, order state or consent is selected.
type Header struct {
	DocumentType DocumentType
	ID           string
	QuoteID      string
	CustomerID   string
	Currency     string
}

type Line struct {
	DocumentType   DocumentType
	DocumentID     string
	VariantID      string
	Quantity       int64
	UnitPriceMinor int64
}

func TransferQuoteToOrder(source Header, lines []Line, orderID string) (Header, []Line, error) {
	// OnRun:37 TestField(Document Type, Quote).
	if source.DocumentType != Quote {
		return Header{}, nil, ErrQuoteType
	}
	// CreateSalesHeader:126-132 copies the header, retargets type/number and
	// preserves the source quote number. Number-series allocation stays caller glue.
	order := source
	order.DocumentType = Order
	order.ID = orderID
	order.QuoteID = source.ID
	result := make([]Line, 0, len(lines))
	for _, sourceLine := range lines {
		// TransferQuoteToOrderLines:331-333 scopes the source records; :340-342
		// copies a matching line and retargets only its document type/number.
		if sourceLine.DocumentType != source.DocumentType || sourceLine.DocumentID != source.ID {
			continue
		}
		line := sourceLine
		line.DocumentType = Order
		line.DocumentID = order.ID
		result = append(result, line)
	}
	return order, result, nil
}
