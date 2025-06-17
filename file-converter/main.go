// === main.go ===
package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "file-converter/internal/api"
)

func main() {
    r := mux.NewRouter()
    
    // Routes API
    api.RegisterRoutes(r)
    
    // Démarrer le serveur
    log.Println("Serveur démarré sur le port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

