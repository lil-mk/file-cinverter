
// === pkg/utils/utils.go ===
package utils

import (
    "path/filepath"
    "strings"
)

// GetFileExtension retourne l'extension d'un fichier sans le point
func GetFileExtension(filename string) string {
    return strings.TrimPrefix(filepath.Ext(filename), ".")
}