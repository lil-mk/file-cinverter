// internal/converter/converter.go
package converter

import (
    "fmt"
    "mime/multipart"
    "os"
    "io"
    "path/filepath"
)

// SupportedFormat représente un format de fichier supporté
type SupportedFormat struct {
    Name        string
    Extension   string
    ContentType string
}

// Interface principale pour la conversion
type Converter interface {
    Convert(input []byte, outputFormat string) ([]byte, error)
    GetSupportedFormats() []SupportedFormat
}

// BaseConverter implémente les fonctionnalités communes
type BaseConverter struct{}

// ValidateFormat vérifie si un format est valide
func ValidateFormat(format string, supportedFormats []SupportedFormat) error {
    for _, f := range supportedFormats {
        if f.Extension == format {
            return nil
        }
    }
    return fmt.Errorf("format non supporté: %s", format)
}

// ConvertFile convertit un fichier
func ConvertFile(conv Converter, inputPath string, outputFormat string) error {
    input, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("erreur de lecture: %v", err)
    }

    output, err := conv.Convert(input, outputFormat)
    if err != nil {
        return fmt.Errorf("erreur de conversion: %v", err)
    }

    outputPath := fmt.Sprintf("%s.%s", 
        inputPath[:len(inputPath)-len(filepath.Ext(inputPath))],
        outputFormat)

    return os.WriteFile(outputPath, output, 0644)
}

// ConvertMultipartFile convertit un fichier multipart
func ConvertMultipartFile(conv Converter, file *multipart.FileHeader, format string) ([]byte, error) {
    src, err := file.Open()
    if err != nil {
        return nil, fmt.Errorf("erreur d'ouverture: %v", err)
    }
    defer src.Close()

    content, err := io.ReadAll(src)
    if err != nil {
        return nil, fmt.Errorf("erreur de lecture: %v", err)
    }

    return conv.Convert(content, format)
}