package main

import (
	// "encoding/json"
	"fmt"
	// "net/http"
	// "sync"
	// "time"

	"github.com/google/uuid"
	// "github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	fmt.Println("starting application server")

	fmt.Print(uuid.New().String())
}
