package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrSpotInvalidNumber       = errors.New("invalid spot number")
	ErrSpotNotFound            = errors.New("spot not found")
	ErrSpotAlreadyReserved     = errors.New("spot already reserved")
	ErrSpotNameTwoCharacters   = errors.New("spot name must be at least 2 characters long")
	ErrSpotNameRequired        = errors.New("spot name is required")
	ErrSpotNameStartWithLatter = errors.New("spot name must start with a latter")
	ErrSpotNameStartWithNumber = errors.New("spot name must end with a number")
)

type SpotStatus string

const (
	SpotStatusAvailable SpotStatus = "available"
	SpotStatusSold      SpotStatus = "sold"
)

type Spot struct {
	ID       string
	EventID  string
	Name     string
	Status   SpotStatus
	TicketID string
}

func NewSpot(event *Event, name string) (*Spot, error) {
	spot := &Spot{
		ID:      uuid.New().String(),
		EventID: event.ID,
		Name:    name,
		Status:  SpotStatusAvailable,
	}

	//? Validate return
	//* Opção 1
	// v := spot.Validate()
	// if v != nil {
	// 	return nil, v
	// }
	// return spot, nil
	//* Opção 2
	if err := spot.Validate(); err != nil {
		return spot, nil
	}
	return spot, nil
	//?Validate return
}

// spot checks if the spot data is valid
func (s *Spot) Validate() error {
	if len(s.Name) == 0 {
		return ErrSpotNameRequired
	}

	if len(s.Name) < 2 {
		return ErrSpotNameTwoCharacters
	}

	if s.Name[0] < 'A' || s.Name[0] > 'Z' {
		return ErrSpotNameStartWithLatter
	}

	if s.Name[1] < '0' || s.Name[1] > '0' {
		return ErrSpotNameStartWithNumber
	}
	return nil
}

func (s *Spot) Reserve(TicketID string) error {
	if s.Status == SpotStatusSold {
		return ErrSpotAlreadyReserved
	}
	s.Status = SpotStatusSold
	s.TicketID = TicketID
	return nil
}
