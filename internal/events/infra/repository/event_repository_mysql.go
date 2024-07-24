package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Eddiesantle/golang-inbound-selling/internal/events/domain"
	_ "github.com/go-sql-driver/mysql"
)

// sqlc - O SQLC é uma ferramenta que facilita a geração de código Go a partir de consultas SQL.
// GORM - O GORM é um Object-Relational Mapping (ORM) para Go que facilita a interação com bancos de dados relacionais.

// mysqlEventRepository é uma implementação do repositório de eventos que usa o banco de dados MySQL.
type mysqlEventRepository struct {
	db *sql.DB // A conexão com o banco de dados.
}

func NewMysqlEventRepository(db *sql.DB) (domain.EventRepository, error) {
	return &mysqlEventRepository{db: db}, nil
}

// CreateSpot insere um novo spot (assento/lugar) no banco de dados.
// Recebe um ponteiro para um objeto Spot do domínio.
func (r *mysqlEventRepository) CreateSpot(spot *domain.Spot) error {

	query := `
		INSERT INTO spots (id, event_id, name, status, ticket_id)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, spot.ID, spot.EventID, spot.Name, spot.Status, spot.TicketID)
	return err
}

// ReserveSpot atualiza o status de um spot para reservado e associa um ticket a ele.
// Recebe os IDs do spot e do ticket.
func (r *mysqlEventRepository) ReserveSpot(spotID, ticketID string) error {
	query := `
		UPDATE spots
		SET status = ?, ticket_id = ?
		WHERE id= ? 
	`

	_, err := r.db.Exec(query, domain.SpotStatusSold, ticketID, spotID)
	return err
}

// CreateTicket insere um novo ticket no banco de dados.
// Recebe um ponteiro para um objeto Ticket do domínio.
func (r *mysqlEventRepository) CreateTicket(ticket *domain.Ticket) error {
	query := `
		INSERT INTO tickets (id, event_id, spot_id, ticket_type, price)
		VALUES(?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, ticket.ID, ticket.EventID, ticket.Spot.ID, ticket.TicketKind, ticket.Price)
	return err
}

// FindEventByID returns an event by its ID, including associated spots and tickets.
func (r *mysqlEventRepository) FindEventByID(eventID string) (*domain.Event, error) {
	query := `
		SELECT 
			e.id, e.name, e.location, e.organization, e.rating, e.date, e.image_url, e.capacity, e.price, e.partner_id,
			s.id, s.event_id, s.name, s.status, s.ticket_id,
			t.id, t.event_id, t.spot_id, t.ticket_kind, t.price
		FROM events e
		LEFT JOIN spots s ON e.id = s.event_id
		LEFT JOIN tickets t ON s.id = t.spot_id
		WHERE e.id = ?
	`
	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var event *domain.Event
	for rows.Next() {
		var eventIDStr, eventName, eventLocation, eventOrganization, eventRating, eventImageURL, spotID, spotEventID, spotName, spotStatus, spotTicketID, ticketID, ticketEventID, ticketSpotID, ticketKind sql.NullString
		var eventDate sql.NullString
		var eventCapacity int
		var eventPrice, ticketPrice sql.NullFloat64
		var partnerID sql.NullInt32

		err := rows.Scan(
			&eventIDStr, &eventName, &eventLocation, &eventOrganization, &eventRating, &eventDate, &eventImageURL, &eventCapacity, &eventPrice, &partnerID,
			&spotID, &spotEventID, &spotName, &spotStatus, &spotTicketID,
			&ticketID, &ticketEventID, &ticketSpotID, &ticketKind, &ticketPrice,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrEventNotFound
			}
			return nil, err
		}

		if !eventIDStr.Valid || !eventName.Valid || !eventLocation.Valid || !eventOrganization.Valid || !eventRating.Valid || !eventDate.Valid || !eventImageURL.Valid || !eventPrice.Valid || !partnerID.Valid {
			continue
		}

		if event == nil {
			eventDateParsed, err := time.Parse("2006-01-02 15:04:05", eventDate.String)
			if err != nil {
				return nil, err
			}
			event = &domain.Event{
				ID:           eventIDStr.String,
				Name:         eventName.String,
				Location:     eventLocation.String,
				Organization: eventOrganization.String,
				Rating:       domain.Rating(eventRating.String),
				Date:         eventDateParsed,
				ImageURL:     eventImageURL.String,
				Capacity:     eventCapacity,
				Price:        eventPrice.Float64,
				PartnerID:    int(partnerID.Int32),
				Spots:        []domain.Spot{},
				Tickets:      []domain.Ticket{},
			}
		}

		if spotID.Valid {
			spot := domain.Spot{
				ID:       spotID.String,
				EventID:  spotEventID.String,
				Name:     spotName.String,
				Status:   domain.SpotStatus(spotStatus.String),
				TicketID: spotTicketID.String,
			}
			event.Spots = append(event.Spots, spot)

			if ticketID.Valid {
				ticket := domain.Ticket{
					ID:         ticketID.String,
					EventID:    ticketEventID.String,
					Spot:       &spot,
					TicketKind: domain.TicketKind(ticketKind.String),
					Price:      ticketPrice.Float64,
				}
				event.Tickets = append(event.Tickets, ticket)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if event == nil {
		return nil, domain.ErrEventNotFound
	}

	return event, nil
}

// FindSpotsByEventID busca os spots de um evento no banco de dados pelo ID do evento.
// Retorna um slice de ponteiros para objetos Spot e um possível erro.
func (r *mysqlEventRepository) FindSpotsByEventID(eventID string) ([]*domain.Spot, error) {
	query := `
		SELECT id, event_id, name, status, ticket_id
		FROM spots
		WHERE event_id = ?
	`

	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spots []*domain.Spot
	// Itera sobre os resultados da query e popula o slice de spots.
	for rows.Next() {
		var spot domain.Spot
		if err := rows.Scan(
			&spot.ID,
			&spot.EventID,
			&spot.Name,
			&spot.Status,
			&spot.TicketID,
		); err != nil {
			return nil, err
		}
		spots = append(spots, &spot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return spots, nil

}

// FindSpotByName busca um spot específico pelo nome e ID do evento no banco de dados.
// Retorna um ponteiro para o objeto Spot e um possível erro.
func (r *mysqlEventRepository) FindSpotByName(eventID, name string) (*domain.Spot, error) {
	query := `
	SELECT
		s.id, s.event_id, s.name, s.status, s.ticket_id,
		t.id, t.event_id, t.spot_id, t.ticket_type, t.price
	FROM spots s
	LEFT JOIN tickets t ON s.id = t.spot_id
	WHERE s.event_id = ? AND s.name = ?
	`
	// Executa a query de busca com o ID do evento e o nome do spot.
	row := r.db.QueryRow(query, eventID, name)

	var spot domain.Spot
	var ticket domain.Ticket
	// Variáveis para armazenar os valores retornados da query.
	var ticketID, ticketEventID, ticketSpotID, ticketKind sql.NullString
	var ticketPrice sql.NullFloat64

	// Faz a leitura do resultado da query para os objetos Spot e Ticket.
	err := row.Scan(
		&spot.ID, &spot.EventID, &spot.Name, &spot.Status, &spot.TicketID,
		&ticketID, &ticketEventID, &ticketSpotID, &ticketKind, &ticketPrice,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSpotNotFound
		}
		return nil, err

	}

	// Verifica se o ticket associado ao spot é válido.
	if ticketID.Valid {
		ticket.ID = ticketID.String
		ticket.EventID = ticketEventID.String
		ticket.Spot = &spot
		ticket.TicketKind = domain.TicketKind(ticketKind.String)
		ticket.Price = ticketPrice.Float64
		spot.TicketID = ticket.ID
	}

	return &spot, nil
}

// ListEvents retorna uma lista de todos os eventos.
func (r *mysqlEventRepository) ListEvents() ([]domain.Event, error) {
	query := `
		SELECT 
			e.id, e.name, e.location, e.organization, e.rating, e.date, e.image_url, e.capacity, e.price, e.partner_id,
			s.id, s.event_id, s.name, s.status, s.ticket_id,
			t.id, t.event_id, t.spot_id, t.ticket_kind, t.price
		FROM events e
		LEFT JOIN spots s ON e.id = s.event_id
		LEFT JOIN tickets t ON s.id = t.spot_id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventMap := make(map[string]*domain.Event)
	spotMap := make(map[string]*domain.Spot)
	for rows.Next() {
		var eventID, eventName, eventLocation, eventOrganization, eventRating, eventImageURL, spotID, spotEventID, spotName, spotStatus, spotTicketID, ticketID, ticketEventID, ticketSpotID, ticketKind sql.NullString
		var eventDate sql.NullString
		var eventCapacity int
		var eventPrice, ticketPrice sql.NullFloat64
		var partnerID sql.NullInt32

		err := rows.Scan(
			&eventID, &eventName, &eventLocation, &eventOrganization, &eventRating, &eventDate, &eventImageURL, &eventCapacity, &eventPrice, &partnerID,
			&spotID, &spotEventID, &spotName, &spotStatus, &spotTicketID,
			&ticketID, &ticketEventID, &ticketSpotID, &ticketKind, &ticketPrice,
		)
		if err != nil {
			return nil, err
		}

		if !eventID.Valid || !eventName.Valid || !eventLocation.Valid || !eventOrganization.Valid || !eventRating.Valid || !eventDate.Valid || !eventImageURL.Valid || !eventPrice.Valid || !partnerID.Valid {
			continue
		}

		event, exists := eventMap[eventID.String]
		if !exists {
			eventDateParsed, err := time.Parse("2006-01-02 15:04:05", eventDate.String)
			if err != nil {
				return nil, err
			}
			event = &domain.Event{
				ID:           eventID.String,
				Name:         eventName.String,
				Location:     eventLocation.String,
				Organization: eventOrganization.String,
				Rating:       domain.Rating(eventRating.String),
				Date:         eventDateParsed,
				ImageURL:     eventImageURL.String,
				Capacity:     eventCapacity,
				Price:        eventPrice.Float64,
				PartnerID:    int(partnerID.Int32),
				Spots:        []domain.Spot{},
				Tickets:      []domain.Ticket{},
			}
			eventMap[eventID.String] = event
		}

		if spotID.Valid {
			spot, spotExists := spotMap[spotID.String]
			if !spotExists {
				spot = &domain.Spot{
					ID:       spotID.String,
					EventID:  spotEventID.String,
					Name:     spotName.String,
					Status:   domain.SpotStatus(spotStatus.String),
					TicketID: spotTicketID.String,
				}
				event.Spots = append(event.Spots, *spot)
				spotMap[spotID.String] = spot
			}

			if ticketID.Valid {
				ticket := domain.Ticket{
					ID:         ticketID.String,
					EventID:    ticketEventID.String,
					Spot:       spot,
					TicketKind: domain.TicketKind(ticketKind.String),
					Price:      ticketPrice.Float64,
				}
				event.Tickets = append(event.Tickets, ticket)
			}
		}
	}

	events := make([]domain.Event, 0, len(eventMap))
	for _, event := range eventMap {
		events = append(events, *event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
