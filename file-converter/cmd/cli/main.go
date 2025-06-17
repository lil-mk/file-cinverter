package main

import (
	"file-converter/internal/converter"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFormat string
	outputDir  string
)

var rootCmd = &cobra.Command{
	Use:   "converter",
	Short: "Convertisseur de fichiers",
	Long: `Un convertisseur de fichiers qui supporte plusieurs formats:
JSON, CSV, XML, TXT et plus encore.`,
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convertit un fichier vers un autre format",
	Long: `Convertit un fichier d'entrée vers le format spécifié.
Exemple: converter convert -i input.json -f csv -o result`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Vérifier que le fichier d'entrée existe
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			return fmt.Errorf("le fichier %s n'existe pas", inputFile)
		}

		// Créer le dossier de sortie s'il n'existe pas
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("impossible de créer le dossier de sortie: %v", err)
		}

		// Lire le fichier d'entrée
		input, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("erreur lors de la lecture du fichier: %v", err)
		}

		// Sélectionner le convertisseur approprié
		var conv converter.Converter
		switch outputFormat {
		case "json", "csv", "xml", "txt":
			conv = &converter.TextConverter{}
		case "jpeg", "png", "gif":
			conv = &converter.ImageConverter{}
		case "gzip", "gunzip":
			conv = &converter.CompressConverter{}
		default:
			return fmt.Errorf("format non supporté: %s", outputFormat)
		}

		// Convertir le fichier
		result, err := conv.Convert(input, outputFormat)
		if err != nil {
			return fmt.Errorf("erreur lors de la conversion: %v", err)
		}

		// Générer le nom du fichier de sortie
		baseName := filepath.Base(inputFile)
		baseNameWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]
		outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.%s", baseNameWithoutExt, outputFormat))

		// Sauvegarder le résultat
		if err := os.WriteFile(outputFile, result, 0644); err != nil {
			return fmt.Errorf("erreur lors de la sauvegarde du fichier: %v", err)
		}

		fmt.Printf("Conversion réussie ! Fichier sauvegardé : %s\n", outputFile)
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Liste les formats supportés",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Formats supportés :")
		fmt.Println("Texte :")
		fmt.Println("  - json")
		fmt.Println("  - csv")
		fmt.Println("  - xml")
		fmt.Println("  - txt")
		fmt.Println("\nImage :")
		fmt.Println("  - jpeg")
		fmt.Println("  - png")
		fmt.Println("  - gif")
		fmt.Println("\nCompression :")
		fmt.Println("  - gzip")
		fmt.Println("  - gunzip")
	},
}

func init() {
	// Ajouter les commandes au rootCmd
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(listCmd)

	// Ajouter les flags à la commande convert
	convertCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Fichier d'entrée à convertir")
	convertCmd.Flags().StringVarP(&outputFormat, "format", "f", "", "Format de sortie")
	convertCmd.Flags().StringVarP(&outputDir, "output", "o", "result", "Dossier de sortie")

	// Marquer les flags requis
	convertCmd.MarkFlagRequired("input")
	convertCmd.MarkFlagRequired("format")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}