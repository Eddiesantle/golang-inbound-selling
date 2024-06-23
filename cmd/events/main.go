package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Eddiesantle/golang-inbound-selling/internal/events/infra/repository"
	"github.com/Eddiesantle/golang-inbound-selling/internal/events/infra/service"
	"github.com/Eddiesantle/golang-inbound-selling/internal/events/usecase"

	httpHandler "github.com/Eddiesantle/golang-inbound-selling/internal/events/infra/http"
)

func main() {
	db, err := sql.Open("mysql", "test_user:test_password@tcp(localgost:3306)/test_go")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventRepo, err := repository.NewMysqlEventRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	partnerBaseURLs := map[int]string{
		1: "http://localhost:9080/api1",
		2: "http://localhost:9080/api2",
	}

	listEventsUseCase := usecase.NewListEventsUseCase(eventRepo)
	getEventUseCase := usecase.NewGetEventUseCase(eventRepo)
	//createEventUseCase := usecase.NewCreateEventUseCase(eventRepo)
	partnerFactory := service.NewPartnerFactory(partnerBaseURLs)
	buyTicketsUseCase := usecase.NewBuyTicketsUseCase(eventRepo, partnerFactory)
	//createSpotsUseCase := usecase.NewCreateSpotsUseCase(eventRepo)
	listSpotsUseCase := usecase.NewListSpotsUseCase(eventRepo)

	eventsHandler := httpHandler.NewEventsHandler(
		listEventsUseCase,
		listSpotsUseCase,
		getEventUseCase,
		buyTicketsUseCase,
	)

	r := http.NewServeMux()
	// r.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	r.HandleFunc("/events", eventsHandler.ListEvents)
	r.HandleFunc("/events/{eventID}", eventsHandler.GetEvent)
	r.HandleFunc("/events/{eventID}/spots", eventsHandler.ListSpots)
	r.HandleFunc("POST /checkout", eventsHandler.BuyTickets)

	http.ListenAndServe(":8080", r)

	// server := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: r,
	// }

	// // Canal para escutar sinais do sistema operacional
	// idleConnsClosed := make(chan struct{})
	// go func() {
	// 	sigint := make(chan os.Signal, 1)
	// 	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
	// 	<-sigint

	// 	// Recebido sinal de interrupção, iniciando o graceful shutdown
	// 	log.Println("Recebido sinal de interrupção, iniciando o graceful shutdown...")

	// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 	defer cancel()

	// 	if err := server.Shutdown(ctx); err != nil {
	// 		log.Printf("Erro no graceful shutdown: %v\n", err)
	// 	}
	// 	close(idleConnsClosed)
	// }()

	// // Iniciando o servidor HTTP
	// log.Println("Servidor HTTP rodando na porta 8080")
	// if err := server.ListenAndServe(); err != http.ErrServerClosed {
	// 	log.Fatalf("Erro ao iniciar o servidor HTTP: %v\n", err)
	// }

	// <-idleConnsClosed
	// log.Println("Servidor HTTP finalizado")
}
