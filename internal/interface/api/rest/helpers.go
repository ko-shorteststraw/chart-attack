package rest

import (
	"bytes"
	"html/template"
	"math"
	"strconv"
)

func renderTemplate(tmpl *template.Template, name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// parseIntSignal converts a signal value (which may be a JSON number, string, or nil) to *int.
func parseIntSignal(v any) *int {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case float64:
		i := int(math.Round(val))
		return &i
	case string:
		if val == "" {
			return nil
		}
		i, err := strconv.Atoi(val)
		if err != nil {
			return nil
		}
		return &i
	}
	return nil
}

// parseFloatSignal converts a signal value to *float64.
func parseFloatSignal(v any) *float64 {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case float64:
		return &val
	case string:
		if val == "" {
			return nil
		}
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil
		}
		return &f
	}
	return nil
}
