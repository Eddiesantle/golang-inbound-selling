package domain
import "Time"

type Rating string

const {
	RatingLivre Rating = "L"
	Rating10 Rating = "L10"
	Rating12 Rating = "L12"
	Rating14 Rating = "L14"
	Rating16 Rating = "L16"
	Rating18 Rating = "L18"

}

type SpotStatus string

const {
	SpotStatusAvailable SpotStatus = "available"
	SpotStatusSold SpotStatus = "sold"
}

type Spot struct{
	ID string
	EventID string
	Name string
	Status SpotStatus
	TicketID string
}




type Event struct {
	ID string
	Name string
	Location string
	Organization string
	Rating Rating
	Date time.Time
	ImageURL string
	Capacity int
	Price float64
	PartnerID int
	Spots []Spots
	// Tickets
}