package main

func setupAPIRoutes(service *Service) {
	// TODO: This will change when API is designed
	service.Mux.Handle("/api/distances", service.Distance)
}
