# Venda de tickets de eventos

## Exploração de entidades
- Event
- Ticket(Ingresso emitido pelo Usuário)
- Spot (Lugar / Cadeira)

## Dominio
- Event
-- Validate
-- AddSpot (Adicionar Spots ao evento)
- Spot
-- Validate%
-- Reserve (reserva lugar - Processo de compra)
-- Domain Service: GenerateSpots
- Ticket
-- Calculate Price
-- Validate
- Definição aa acessps externos (Repository)

## Infra / Repository - Acesso ao banco de dados
- Repository
-- ListEvents
-- FindEventByIs
-- CreateSpot
-- CreateTicket
-- ReserveSpot
-- FindSpotsByEventID
-- FindSpotByName

## Usecases
-- ListEvents
-- GetEvents
-- ListSpots
-- BuyTickets

## main
Entrypoint

