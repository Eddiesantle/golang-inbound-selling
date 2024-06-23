package repository

import (
	"database/sql"
	"errors"

	"github.com/Eddiesantle/golang-inbound-selling/internal/events/domain"
	_ "github.com/go-sql-driver/mysql"
)

// sqlc - O SQLC é uma ferramenta que facilita a geração de código Go a partir de consultas SQL.
// GORM - O GORM é um Object-Relational Mapping (ORM) para Go que facilita a interação com bancos de dados relacionais.

// mysqlEventRepository é uma implementação do repositório de eventos que usa o banco de dados MySQL.
type mysqlEventRepository struct {
	db *sql.DB // A conexão com o banco de dados.
}

// func NewMysqlEventRepository(db *sql.DB) (domain.EventRepository, error) {
// 	return &mysqlEventRepository{db: db}, nil
// }

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
	_, err := r.db.Exec(query, ticket.ID, ticket.EventID, ticket.Spot.ID, ticket.TicketType, ticket.Price)
	return err
}

// FindEventByID busca um evento no banco de dados pelo seu ID.
// Retorna um ponteiro para o objeto Event e um possível erro.
func (r *mysqlEventRepository) FindEventByID(eventID string) (*domain.Event, error) {
	query := `
		SELECT id, name, location, organization, tarion, date, image_url, capacity, price, partner_id
		FROM events
		WHERE id = ?
	`

	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var event *domain.Event
	// Faz a leitura do resultado da query para o objeto Event.
	err = rows.Scan(
		&event.ID,
		&event.Name,
		&event.Location,
		&event.Organization,
		&event.Rating,
		&event.Date,
		&event.ImageURL,
		&event.Capacity,
		&event.Price,
		&event.PartnerID,
	)

	if err != nil {
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
	var ticketID, ticketEventID, ticketSpotID, ticketType sql.NullString
	var ticketPrice sql.NullFloat64

	// Faz a leitura do resultado da query para os objetos Spot e Ticket.
	err := row.Scan(
		&spot.ID, &spot.EventID, &spot.Name, &spot.Status, &spot.TicketID,
		&ticketID, &ticketEventID, &ticketSpotID, &ticketType, &ticketPrice,
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
		ticket.TicketType = domain.TicketType(ticketType.String)
		ticket.Price = ticketPrice.Float64
		spot.TicketID = ticket.ID
	}

	return &spot, nil
}
