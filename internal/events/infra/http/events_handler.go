package http

import (
	"encoding/json"
	"net/http"

	"github.com/Eddiesantle/golang-inbound-selling/internal/events/usecase"
)

type EventsHandler struct {
	listEventsUseCase  *usecase.ListEventsUseCase
	listSpotsUseCase   *usecase.ListSpotsUseCase
	getEventUseCase    *usecase.GetEventUseCase
	createEventUseCase *usecase.CreateEventUseCase
	buyTicketsUseCase  *usecase.BuyTicketsUseCase
	createSpotsUseCase *usecase.CreateSpotsUseCase
}

func NewEventsHandler(
	listEventsUseCase *usecase.ListEventsUseCase,
	listSpotsUseCase *usecase.ListSpotsUseCase,
	getEventUseCase *usecase.GetEventUseCase,
	buyTicketsUseCase *usecase.BuyTicketsUseCase,
	createEventUseCase *usecase.CreateEventUseCase,
	createSpotsUseCase *usecase.CreateSpotsUseCase,
) *EventsHandler {
	return &EventsHandler{
		listEventsUseCase:  listEventsUseCase,
		listSpotsUseCase:   listSpotsUseCase,
		getEventUseCase:    getEventUseCase,
		buyTicketsUseCase:  buyTicketsUseCase,
		createEventUseCase: createEventUseCase,
		createSpotsUseCase: createSpotsUseCase,
	}
}

// ListEvents handles the request to list all events.
// @Summary List all events
// @Description Get all events with their details
// @Tags Events
// @Accept json
// @Produce json
// @Success 200 {object} usecase.ListEventsOutputDTO
// @Failure 500 {object} string
// @Router /events [get]
func (h *EventsHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	output, err := h.listEventsUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// GetEvent handles the request to get details of a specific event.
// @Summary Get event details
// @Description Get details of an event by ID
// @Tags Events
// @Accept json
// @Produce json
// @Param eventID path string true "Event ID"
// @Success 200 {object} usecase.GetEventOutputDTO
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /events/{eventID} [get]
func (h *EventsHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	input := usecase.GetEventInputDTO{ID: eventID}

	output, err := h.getEventUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// ListSpots lists spots for an event.
// @Summary List spots for an event
// @Description List all spots for a specific event
// @Tags Events
// @Accept json
// @Produce json
// @Param eventID path string true "Event ID"
// @Success 200 {object} usecase.ListSpotsOutputDTO
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /events/{eventID}/spots [get]
func (h *EventsHandler) ListSpots(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	input := usecase.ListSpotsInputDTO{EventID: eventID}

	output, err := h.listSpotsUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// BuyTickets handles the request to buy tickets for an event.
// @Summary Buy tickets for an event
// @Description Buy tickets for a specific event
// @Tags Events
// @Accept json
// @Produce json
// @Param input body usecase.BuyTicketsInputDTO true "Input data"
// @Success 200 {object} usecase.BuyTicketsOutputDTO
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /checkout [post]
func (h *EventsHandler) BuyTickets(w http.ResponseWriter, r *http.Request) {
	var input usecase.BuyTicketsInputDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := h.buyTicketsUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// CreateEvent handles the request create event.
// @Summary create event
// @Description create event
// @Tags Events
// @Accept json
// @Produce json
// @Param input body usecase.CreateEventInputDTO true "Input data"
// @Success 200 {object} usecase.CreateEventInputDTO
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /event [post]
func (h *EventsHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateEventInputDTO

	// Adicione logs para verificar o método HTTP e o corpo da requisição
	//fmt.Println("Method:", r.Method)

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.createEventUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *EventsHandler) CreateSpots(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	var input usecase.CreateSpotsInputDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input.EventID = eventID

	output, err := h.createSpotsUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
