// === internal/converter/image.go ===
package converter

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
)

type ImageConverter struct{}

// GetSupportedFormats implements Converter.
func (i *ImageConverter) GetSupportedFormats() []SupportedFormat {
	panic("unimplemented")
}

func (i *ImageConverter) Convert(input []byte, outputFormat string) ([]byte, error) {
	// Décoder l'image d'entrée
	img, format, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("erreur de décodage de l'image: %v", err)
	}
	fmt.Printf("Format d'entrée détecté : %s\n", format)

	buf := new(bytes.Buffer)
	switch outputFormat {
	case "jpeg":
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 85})
	case "png":
		err = png.Encode(buf, img)
	case "gif":
		err = gif.Encode(buf, img, &gif.Options{NumColors: 256})
	default:
		return nil, fmt.Errorf("format d'image non supporté: %s", outputFormat)
	}

	if err != nil {
		return nil, fmt.Errorf("erreur d'encodage de l'image: %v", err)
	}

	return buf.Bytes(), nil
}
