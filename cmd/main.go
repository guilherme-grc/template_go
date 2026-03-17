package main

import "template.go/config"

func main() {
	cfg := config.Load()

	// Agora você usa os dados
	println("Iniciando servidor na porta: " + cfg.AppPort)
	// r.Run(":" + cfg.AppPort)
}
