package main

import (
	"net/http"
	"os"

	purchaseRepository "haf.systems/purchases/repositories/file/purchaseRepository"
	exchangeRepository "haf.systems/purchases/repositories/treasury/exchangeRateRepository"

	services "haf.systems/purchases/services"

	controllers "haf.systems/purchases/controllers"

	server "haf.systems/purchases/server"
)

const PURCHASES_FILENAME string = "purchases.json"

func main() {
	port := os.Getenv("PORT")

	// ctx := context.Background()
	httpClient := &http.Client{}
	defer httpClient.CloseIdleConnections()

	purchaseRepository := purchaseRepository.NewPurchaseFileRepository(PURCHASES_FILENAME)
	defer purchaseRepository.Close()
	exchangeRepository := exchangeRepository.NewExchangeRateTreasuryRepository(httpClient)

	services := services.NewServices(purchaseRepository, exchangeRepository)
	controllers := controllers.NewControllers(services)

	server := server.NewServer(controllers, port)
	server.CorsSetup()

	server.Run()
}
