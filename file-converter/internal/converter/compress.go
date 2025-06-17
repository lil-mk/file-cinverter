// internal/converter/compress.go
package converter

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// CompressConverter gère la compression et décompression des fichiers
type CompressConverter struct {
	BaseConverter
	CompressionLevel int // Niveau de compression (1-9, par défaut 6)
}

// NewCompressConverter crée une nouvelle instance de CompressConverter
func NewCompressConverter() *CompressConverter {
	return &CompressConverter{
		CompressionLevel: gzip.DefaultCompression,
	}
}

// Formats supportés pour la compression
var CompressFormats = []SupportedFormat{
	{"GZIP", "gz", "application/gzip"},
	{"ZLIB", "zlib", "application/zlib"},
	{"ZIP", "zip", "application/zip"},
}

// Convert implémente l'interface Converter pour la compression
func (c *CompressConverter) Convert(input []byte, outputFormat string) ([]byte, error) {
	switch outputFormat {
	case "gz", "gzip":
		return c.compressGzip(input)
	case "gunzip":
		return c.decompressGzip(input)
	case "zlib":
		return c.compressZlib(input)
	case "unzlib":
		return c.decompressZlib(input)
	case "zip":
		return c.compressZip(input)
	case "unzip":
		return c.decompressZip(input)
	default:
		return nil, fmt.Errorf("format de compression non supporté: %s", outputFormat)
	}
}

// Compression GZIP
func (c *CompressConverter) compressGzip(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	gw, err := gzip.NewWriterLevel(&buf, c.CompressionLevel)
	if err != nil {
		return nil, fmt.Errorf("erreur d'initialisation gzip: %v", err)
	}

	if _, err := gw.Write(input); err != nil {
		return nil, fmt.Errorf("erreur de compression gzip: %v", err)
	}

	if err := gw.Close(); err != nil {
		return nil, fmt.Errorf("erreur de fermeture gzip: %v", err)
	}

	return buf.Bytes(), nil
}

// Décompression GZIP
func (c *CompressConverter) decompressGzip(input []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("erreur d'ouverture gzip: %v", err)
	}
	defer gr.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, gr); err != nil {
		return nil, fmt.Errorf("erreur de décompression gzip: %v", err)
	}

	return buf.Bytes(), nil
}

// Compression ZLIB
func (c *CompressConverter) compressZlib(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw, err := zlib.NewWriterLevel(&buf, c.CompressionLevel)
	if err != nil {
		return nil, fmt.Errorf("erreur d'initialisation zlib: %v", err)
	}

	if _, err := zw.Write(input); err != nil {
		return nil, fmt.Errorf("erreur de compression zlib: %v", err)
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("erreur de fermeture zlib: %v", err)
	}

	return buf.Bytes(), nil
}

// Décompression ZLIB
func (c *CompressConverter) decompressZlib(input []byte) ([]byte, error) {
	zr, err := zlib.NewReader(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("erreur d'ouverture zlib: %v", err)
	}
	defer zr.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, zr); err != nil {
		return nil, fmt.Errorf("erreur de décompression zlib: %v", err)
	}

	return buf.Bytes(), nil
}

// Compression ZIP
func (c *CompressConverter) compressZip(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	// Créer un fichier dans l'archive
	f, err := zw.Create("file")
	if err != nil {
		return nil, fmt.Errorf("erreur de création du fichier zip: %v", err)
	}

	// Écrire les données
	if _, err := f.Write(input); err != nil {
		return nil, fmt.Errorf("erreur d'écriture zip: %v", err)
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("erreur de fermeture zip: %v", err)
	}

	return buf.Bytes(), nil
}

// Décompression ZIP
func (c *CompressConverter) decompressZip(input []byte) ([]byte, error) {
	reader := bytes.NewReader(input)
	zr, err := zip.NewReader(reader, int64(len(input)))
	if err != nil {
		return nil, fmt.Errorf("erreur d'ouverture zip: %v", err)
	}

	// Lire le premier fichier de l'archive
	if len(zr.File) == 0 {
		return nil, fmt.Errorf("archive zip vide")
	}

	// Ouvrir le fichier
	f, err := zr.File[0].Open()
	if err != nil {
		return nil, fmt.Errorf("erreur d'ouverture du fichier dans zip: %v", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		return nil, fmt.Errorf("erreur de lecture du fichier zip: %v", err)
	}

	return buf.Bytes(), nil
}

// SetCompressionLevel définit le niveau de compression
func (c *CompressConverter) SetCompressionLevel(level int) error {
	if level < gzip.BestSpeed || level > gzip.BestCompression {
		return fmt.Errorf("niveau de compression invalide: %d (doit être entre %d et %d)",
			level, gzip.BestSpeed, gzip.BestCompression)
	}
	c.CompressionLevel = level
	return nil
}

// GetSupportedFormats retourne les formats supportés
func (c *CompressConverter) GetSupportedFormats() []SupportedFormat {
	return CompressFormats
}

// Méthodes utilitaires supplémentaires

// IsCompressedFile vérifie si un fichier est déjà compressé
func (c *CompressConverter) IsCompressedFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	compressedExts := []string{".gz", ".gzip", ".zip", ".zlib"}
	for _, compressedExt := range compressedExts {
		if ext == compressedExt {
			return true
		}
	}
	return false
}

// GetCompressionRatio calcule le ratio de compression
func (c *CompressConverter) GetCompressionRatio(original, compressed []byte) float64 {
	if len(original) == 0 {
		return 0
	}
	return float64(len(compressed)) / float64(len(original))
}
