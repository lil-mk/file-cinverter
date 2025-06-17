package converter

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

type TextConverter struct{}

// GetSupportedFormats retourne les formats supportés
func (t *TextConverter) GetSupportedFormats() []SupportedFormat {
	return []SupportedFormat{
		{Name: "JSON", Extension: "json", ContentType: "application/json"},
		{Name: "CSV", Extension: "csv", ContentType: "text/csv"},
		{Name: "XML", Extension: "xml", ContentType: "application/xml"},
		{Name: "Text", Extension: "txt", ContentType: "text/plain"},
	}
}

// Structure pour la conversion XML
type XMLData struct {
	XMLName xml.Name    `xml:"root"`
	Items   []XMLRecord `xml:"item"`
}

type XMLRecord struct {
	Fields []XMLField `xml:"field"`
}

type XMLField struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func (t *TextConverter) Convert(input []byte, outputFormat string) ([]byte, error) {
	// Valider le format de sortie
	if err := ValidateFormat(outputFormat, t.GetSupportedFormats()); err != nil {
		return nil, err
	}

	// Détecter le format d'entrée
	inputFormat := detectFormat(input)
	// Vérifier si le format d'entrée est supporté
	if err := ValidateFormat(inputFormat, t.GetSupportedFormats()); err != nil {
		return nil, fmt.Errorf("format d'entrée non reconnu: %s", inputFormat)
	}

	// Convertir les données en structure intermédiaire
	var data []map[string]interface{}

	switch inputFormat {
	case "json":
		if err := json.Unmarshal(input, &data); err != nil {
			return nil, fmt.Errorf("erreur lors du parsing JSON: %v", err)
		}
	case "csv":
		reader := csv.NewReader(bytes.NewReader(input))
		records, err := reader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("erreur lors du parsing CSV: %v", err)
		}

		if len(records) < 2 {
			return nil, fmt.Errorf("CSV invalide: besoin d'au moins un en-tête et une ligne de données")
		}

		headers := records[0]
		data = make([]map[string]interface{}, 0, len(records)-1)

		for _, record := range records[1:] {
			item := make(map[string]interface{})
			for i, value := range record {
				if i < len(headers) {
					item[headers[i]] = value
				}
			}
			data = append(data, item)
		}
	case "xml":
		var xmlData XMLData
		if err := xml.Unmarshal(input, &xmlData); err != nil {
			return nil, fmt.Errorf("erreur lors du parsing XML: %v", err)
		}

		data = make([]map[string]interface{}, len(xmlData.Items))
		for i, item := range xmlData.Items {
			record := make(map[string]interface{})
			for _, field := range item.Fields {
				record[field.Name] = field.Value
			}
			data[i] = record
		}
	case "txt":
		lines := strings.Split(string(input), "\n")
		data = make([]map[string]interface{}, len(lines))
		
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				data[i] = map[string]interface{}{
					"line": line,
				}
			}
		}
		
		cleanData := []map[string]interface{}{}
		for _, item := range data {
			if item != nil {
				cleanData = append(cleanData, item)
			}
		}
		data = cleanData
	}

	// Convertir vers le format de sortie
	switch outputFormat {
	case "json":
		return json.MarshalIndent(data, "", "  ")
	case "csv":
		if len(data) == 0 {
			return nil, fmt.Errorf("pas de données à convertir")
		}

		headers := make([]string, 0)
		for k := range data[0] {
			headers = append(headers, k)
		}

		buf := new(bytes.Buffer)
		writer := csv.NewWriter(buf)

		if err := writer.Write(headers); err != nil {
			return nil, fmt.Errorf("erreur lors de l'écriture des en-têtes CSV: %v", err)
		}

		for _, item := range data {
			record := make([]string, len(headers))
			for i, header := range headers {
				if val, ok := item[header]; ok {
					record[i] = fmt.Sprint(val)
				}
			}
			if err := writer.Write(record); err != nil {
				return nil, fmt.Errorf("erreur lors de l'écriture des données CSV: %v", err)
			}
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			return nil, fmt.Errorf("erreur lors de la finalisation du CSV: %v", err)
		}
		return buf.Bytes(), nil
	case "xml":
		xmlData := XMLData{
			Items: make([]XMLRecord, len(data)),
		}

		for i, item := range data {
			var fields []XMLField
			for key, value := range item {
				fields = append(fields, XMLField{
					Name:  key,
					Value: fmt.Sprint(value),
				})
			}
			xmlData.Items[i] = XMLRecord{Fields: fields}
		}

		buf := new(bytes.Buffer)
		buf.WriteString(xml.Header)

		encoder := xml.NewEncoder(buf)
		encoder.Indent("", "  ")
		if err := encoder.Encode(xmlData); err != nil {
			return nil, fmt.Errorf("erreur lors de l'encodage XML: %v", err)
		}

		return buf.Bytes(), nil
	case "txt":
		var builder strings.Builder
		for _, item := range data {
			for key, value := range item {
				builder.WriteString(fmt.Sprintf("%s: %v\n", key, value))
			}
			builder.WriteString("\n")
		}
		return []byte(builder.String()), nil
	}

	return nil, nil
}

func detectFormat(input []byte) string {
	trimmed := bytes.TrimSpace(input)
	if len(trimmed) == 0 {
		return ""
	}

	if trimmed[0] == '{' || trimmed[0] == '[' {
		return "json"
	}

	if trimmed[0] == '<' {
		return "xml"
	}

	firstLine := bytes.SplitN(trimmed, []byte{'\n'}, 2)[0]
	if bytes.Count(firstLine, []byte{','}) > 0 && !bytes.ContainsAny(firstLine, "{}[]<>") {
		return "csv"
	}

	return "txt"
}