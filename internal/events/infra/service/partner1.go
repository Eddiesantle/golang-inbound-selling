package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Partner1 representa um parceiro externo que processa reservas.
type Partner1 struct {
	BaseURL string // URL base para a API do parceiro.
}

// Partner1ReservationRequest estrutura de solicitação para reservar spots via Partner1.
type Partner1ReservationRequest struct {
	Spots      []string `json:"spots"`       // Lista de spots que estão sendo reservados.
	TicketKind string   `json:"ticket_kind"` // Tipo de ingresso (por exemplo, "meia", "inteira").
	Email      string   `json:"email"`       // E-mail do usuário que está fazendo a reserva.
}

// Partner1ReservationResponse estrutura de resposta do parceiro após a reserva.
type Partner1ReservationResponse struct {
	ID         string `json:"id"`          // ID da reserva gerada pelo parceiro.
	Email      string `json:"email"`       // E-mail do usuário para confirmação.
	Spot       string `json:"spot"`        // Nome do spot reservado.
	TicketKind string `json:"ticket_kind"` // Tipo de ingresso reservado.
	Status     string `json:"status"`      // Status da reserva (por exemplo, "confirmado").
	EventID    string `json:"event_id"`    // ID do evento associado à reserva.
}

// MakeReservation envia uma solicitação de reserva para o parceiro e retorna as respostas da reserva.
func (p *Partner1) MakeReservation(req *ReservationRequest) ([]ReservationResponse, error) {

	// Converte a solicitação de reserva genérica para o formato específico do parceiro.
	partnerReq := Partner1ReservationRequest{
		Spots:      req.Spots,
		TicketKind: req.TicketKind,
		Email:      req.Email,
	}

	// Serializa a solicitação do parceiro em JSON.
	body, err := json.Marshal(partnerReq)
	if err != nil {
		return nil, err
	}

	// Constrói a URL para a solicitação de reserva, incluindo o ID do evento.
	url := fmt.Sprintf("%s/event/%s/reserve", p.BaseURL, req.EventID)

	// Cria uma nova solicitação HTTP do tipo POST com o corpo da solicitação JSON.
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Cria um cliente HTTP e envia a solicitação.
	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	// Fecha o corpo da resposta quando a função terminar.
	defer httpResp.Body.Close()

	// Verifica se o código de status HTTP é 201 (Created).
	if httpResp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	// Decodifica a resposta JSON do parceiro em uma slice de Partner1ReservationResponse.
	var partnerResp []Partner1ReservationResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&partnerResp); err != nil {
		return nil, err
	}

	// Converte as respostas do parceiro para o formato genérico de ReservationResponse.
	responses := make([]ReservationResponse, len(partnerResp))
	for i, r := range partnerResp {
		responses[i] = ReservationResponse{
			ID:     r.ID,
			Spot:   r.Spot,
			Status: r.Status,
		}
	}

	// Retorna as respostas da reserva.
	return responses, nil
}
