package main

import (
	"fmt"
	"log/slog"

	"github.com/OzkrOssa/freeradius-api/internal/adapter/config"
)

func main() {
	config, err := config.New()
	if err != nil {
		slog.Error("error reading config", "ERROR", err)
	}
}
