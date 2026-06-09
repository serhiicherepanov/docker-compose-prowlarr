package lampa_categories

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// Config contains middleware settings.
type Config struct {
	ParameterName string `json:"parameterName,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		ParameterName: "categories",
	}
}

type lampaCategories struct {
	next          http.Handler
	parameterName string
}

// New creates the middleware instance.
func New(_ context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	parameterName := strings.TrimSpace(config.ParameterName)
	if parameterName == "" {
		parameterName = "categories"
	}

	return &lampaCategories{
		next:          next,
		parameterName: parameterName,
	}, nil
}

func (m *lampaCategories) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.URL.RawQuery = splitCommaSeparatedParameter(req.URL.RawQuery, m.parameterName)
	m.next.ServeHTTP(rw, req)
}

func splitCommaSeparatedParameter(rawQuery string, parameterName string) string {
	if rawQuery == "" {
		return rawQuery
	}

	parts := strings.Split(rawQuery, "&")
	rewritten := make([]string, 0, len(parts))

	for _, part := range parts {
		name, value, hasValue := strings.Cut(part, "=")
		decodedName, err := url.QueryUnescape(name)
		if err != nil || decodedName != parameterName || !hasValue {
			rewritten = append(rewritten, part)
			continue
		}

		decodedValue, err := url.QueryUnescape(value)
		if err != nil || !strings.Contains(decodedValue, ",") {
			rewritten = append(rewritten, part)
			continue
		}

		values := strings.Split(decodedValue, ",")
		encodedName := url.QueryEscape(parameterName)
		for _, item := range values {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}

			rewritten = append(rewritten, encodedName+"="+url.QueryEscape(item))
		}
	}

	return strings.Join(rewritten, "&")
}
