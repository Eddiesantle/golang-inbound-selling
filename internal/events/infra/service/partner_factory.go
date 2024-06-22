package service

import "fmt"

// PartnerFactory é uma interface que define o método para criar um parceiro.
type PartnerFactory interface {
	// CreatePartner cria e retorna uma instância de um parceiro com base no partnerID fornecido.
	CreatePartner(partnerID int) (Partner, error)
}

// DefaultPartnerFactory é a implementação padrão da interface PartnerFactory.
type DefaultPartnerFactory struct {
	partnerBaseURLs map[int]string // Map de IDs de parceiros para suas URLs base.
}

// NewPartnerFactory cria uma nova instância de DefaultPartnerFactory.
// Recebe um mapa de IDs de parceiros para URLs base e retorna uma interface PartnerFactory.
func NewPartnerFactory(partnerBaseURLs map[int]string) PartnerFactory {
	// Retorna a fábrica de parceiros com o mapa de URLs base inicializado.
	return &DefaultPartnerFactory{partnerBaseURLs: partnerBaseURLs}
}

// CreatePartner cria um parceiro com base no partnerID fornecido.
// Retorna um erro se o parceiro não for encontrado.
func (f *DefaultPartnerFactory) CreatePartner(partnerID int) (Partner, error) {
	// Busca a URL base do parceiro a partir do partnerID.
	BaseURL, ok := f.partnerBaseURLs[partnerID]
	if !ok {
		return nil, fmt.Errorf("partner with ID %d not found", partnerID)
	}

	// Cria e retorna a instância do parceiro apropriado com base no partnerID.
	// Adicione casos adicionais conforme novos parceiros forem introduzidos.
	switch partnerID {
	case 1:
		return &Partner1{BaseURL: BaseURL}, nil
	case 2:
		return &Partner1{BaseURL: BaseURL}, nil
	default:
		return nil, fmt.Errorf("partner with ID %d not found", partnerID)
	}
}
