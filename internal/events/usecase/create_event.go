package usecase

import (
	"time"

	"github.com/Eddiesantle/golang-inbound-selling/internal/events/domain"
)

type CreateEventInputDTO struct {
	Name         string    `json:"name"`
	Location     string    `json:"location"`
	Organization string    `json:"organization"`
	Rating       string    `json:"rating"`
	Date         time.Time `json:"date"`
	Capacity     int       `json:"capacity"`
	ImageURL     string    `json:"image_url"`
	Price        float64   `json:"price"`
	PartnerID    int       `json:"partner_id"`
}

type CreateEventOutputDTO struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Location     string    `json:"location"`
	Organization string    `json:"organization"`
	Rating       string    `json:"rating"`
	Date         time.Time `json:"date"`
	ImageURL     string    `json:"image_url"`
	Capacity     int       `json:"capacity"`
	Price        float64   `json:"price"`
	PartnerID    int       `json:"partner_id"`
}

type CreateEventUseCase struct {
	repo domain.EventRepository
}

func NewCreateEventUseCase(repo domain.EventRepository) *CreateEventUseCase {
	return &CreateEventUseCase{repo: repo}
}

func (uc *CreateEventUseCase) Execute(input CreateEventInputDTO) (CreateEventOutputDTO, error) {

	event, err := domain.NewEvent(
		input.Name,
		input.Location,
		input.Organization,
		domain.Rating(input.Rating),
		input.Date,
		input.Capacity,
		input.Price,
		input.ImageURL,
		input.PartnerID,
	)

	if err != nil {
		return CreateEventOutputDTO{}, err
	}

	err = uc.repo.CreateEvent(event)
	if err != nil {
		return CreateEventOutputDTO{}, err
	}

	output := CreateEventOutputDTO{
		ID:           event.ID,
		Name:         event.Name,
		Location:     event.Location,
		Organization: event.Organization,
		Rating:       string(event.Rating),
		Date:         event.Date,
		Capacity:     event.Capacity,
		ImageURL:     event.ImageURL,
		Price:        event.Price,
		PartnerID:    event.PartnerID,
	}

	return output, nil
}
