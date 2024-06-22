package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Partner2 representa um parceiro externo que processa reservas.
type Partner2 struct {
	BaseURL string // URL base para a API do parceiro.
}

// Partner2ReservationRequest estrutura de solicitação para reservar spots via Partner2.
type Partner2ReservationRequest struct {
	Lugares      []string `json:"lugares"`       // Lista de spots que estão sendo reservados.
	TipoIngresso string   `json:"tipo_ingresso"` // Tipo de ingresso (por exemplo, "meia", "inteira").
	Email        string   `json:"email"`         // E-mail do usuário que está fazendo a reserva.
}

// Partner2ReservationResponse estrutura de resposta do parceiro após a reserva.
type Partner2ReservationResponse struct {
	ID           string `json:"id"`            // ID da reserva gerada pelo parceiro.
	Email        string `json:"email"`         // E-mail do usuário para confirmação.
	Lugar        string `json:"lugar"`         // Nome do spot reservado.
	TipoIngresso string `json:"tipo_ingresso"` // Tipo de ingresso reservado.
	Estado       string `json:"estado"`        // Status da reserva (por exemplo, "confirmado").
	EventID      string `json:"evento_id"`     // ID do evento associado à reserva.
}

// MakeReservation envia uma solicitação de reserva para o parceiro e retorna as respostas da reserva.
func (p *Partner2) MakeReservation(req *ReservationRequest) ([]ReservationResponse, error) {

	// Converte a solicitação de reserva genérica para o formato específico do parceiro.
	partnerReq := Partner2ReservationRequest{
		Lugares:      req.Spots,
		TipoIngresso: req.TicketType,
		Email:        req.Email,
	}

	// Serializa a solicitação do parceiro em JSON.
	body, err := json.Marshal(partnerReq)
	if err != nil {
		return nil, err
	}

	// Constrói a URL para a solicitação de reserva, incluindo o ID do evento.
	url := fmt.Sprintf("%s/eventos/%s/reservar", p.BaseURL, req.EventID)

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

	// Decodifica a resposta JSON do parceiro em uma slice de Partner2ReservationResponse.
	var partnerResp []Partner2ReservationResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&partnerResp); err != nil {
		return nil, err
	}

	// Converte as respostas do parceiro para o formato genérico de ReservationResponse.
	responses := make([]ReservationResponse, len(partnerResp))
	for i, r := range partnerResp {
		responses[i] = ReservationResponse{
			ID:     r.ID,
			Spot:   r.Lugar,
			Status: r.Estado,
		}
	}

	// Retorna as respostas da reserva.
	return responses, nil
}
