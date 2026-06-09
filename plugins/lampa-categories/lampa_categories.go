package lampa_categories

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Config contains middleware settings.
type Config struct {
	ParameterName string `json:"parameterName,omitempty"`
	LogRequests   bool   `json:"logRequests,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		ParameterName: "categories",
		LogRequests:   true,
	}
}

type lampaCategories struct {
	next          http.Handler
	parameterName string
	logRequests   bool
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
		logRequests:   config.LogRequests,
	}, nil
}

func (m *lampaCategories) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rawQuery, report := splitCommaSeparatedParameter(req.URL.RawQuery, m.parameterName)
	req.URL.RawQuery = rawQuery
	req.RequestURI = requestURI(req.URL)

	req.Header.Set("X-Lampa-Categories-Plugin", "enabled")
	req.Header.Set("X-Lampa-Categories-Status", report.status)
	if report.original != "" {
		req.Header.Set("X-Lampa-Categories-Original", report.original)
	}
	if report.rewritten != "" {
		req.Header.Set("X-Lampa-Categories-Rewritten", report.rewritten)
	}
	if report.final != "" {
		req.Header.Set("X-Lampa-Categories-Final", report.final)
	}

	if m.logRequests && report.status != "missing" {
		fmt.Printf(
			"lampa-categories path=%q status=%s original=%q rewritten=%q final=%q\n",
			req.URL.Path,
			report.status,
			report.original,
			report.rewritten,
			report.final,
		)
	}

	m.next.ServeHTTP(rw, req)
}

type rewriteReport struct {
	status    string
	original  string
	rewritten string
	final     string
}

func splitCommaSeparatedParameter(rawQuery string, parameterName string) (string, rewriteReport) {
	report := rewriteReport{status: "missing"}
	if rawQuery == "" {
		return rawQuery, report
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
		report.original = decodedValue
		if err != nil || !strings.Contains(decodedValue, ",") {
			report.status = "passthrough"
			rewritten = append(rewritten, part)
			continue
		}

		encodedName := url.QueryEscape(parameterName)
		lastValue := lastNonEmptyValue(decodedValue)
		if lastValue == "" {
			report.status = "rewritten"
			continue
		}

		finalPart := encodedName + "=" + url.QueryEscape(lastValue)
		rewritten = append(rewritten, finalPart)
		report.status = "rewritten"
		report.rewritten = lastValue
		report.final = finalPart
	}

	return strings.Join(rewritten, "&"), report
}

func lastNonEmptyValue(value string) string {
	items := strings.Split(value, ",")
	for i := len(items) - 1; i >= 0; i-- {
		item := strings.TrimSpace(items[i])
		if item != "" {
			return item
		}
	}

	return ""
}

func requestURI(u *url.URL) string {
	if u.RawQuery == "" {
		return u.Path
	}

	return u.Path + "?" + u.RawQuery
}
