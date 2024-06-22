package domain

import "errors"

type TicketType string

const (
	TicketTypeHalf TicketType = "half"
	TicketTypeFull TicketType = "full"
)

var ErrTicketPriceZero = errors.New("Ticket price must be greater than zero")

type Ticket struct {
	ID         string
	EventID    string
	Spot       *Spot
	TicketType TicketType
	Price      float64
}

func IsValidTicketType(ticketType TicketType) bool {
	return ticketType == TicketTypeHalf || ticketType == TicketTypeFull
}

func (t *Ticket) CalculatePrice() {
	if t.TicketType == TicketTypeHalf {
		t.Price /= 2
	}
}

func (t *Ticket) Validate() error {
	if t.Price <= 0 {
		return ErrTicketPriceZero
	}
	return nil
}
