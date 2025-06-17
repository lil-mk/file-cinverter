package api

import (
	"file-converter/internal/converter"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/convert/{format}", convertHandler).Methods("POST")
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	format := vars["format"]

	// Lire le fichier
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du fichier", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du contenu", http.StatusInternalServerError)
		return
	}

	// Sélectionner le convertisseur
	var conv converter.Converter
	switch format {
	case "json", "csv", "xml", "txt":
		conv = &converter.TextConverter{}
	case "jpeg", "png", "gif":
		conv = &converter.ImageConverter{}
	case "gzip", "gunzip":
		conv = &converter.CompressConverter{}
	default:
		http.Error(w, "Format non supporté", http.StatusBadRequest)
		return
	}

	// Convertir
	result, err := conv.Convert(content, format)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Définir le bon Content-Type en fonction du format
	contentType := getContentType(format)
	w.Header().Set("Content-Type", contentType)
	w.Write(result)
}

// getContentType retourne le Content-Type approprié pour chaque format
func getContentType(format string) string {
	switch format {
	case "json":
		return "application/json"
	case "csv":
		return "text/csv"
	case "xml":
		return "application/xml"
	case "txt":
		return "text/plain"
	case "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "gzip":
		return "application/gzip"
	default:
		return "application/octet-stream"
	}
}