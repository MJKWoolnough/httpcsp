package csp

import (
	"encoding/json"
	"net/http"
)

type report struct {
	Report `json:"csp-report"`
}

type Report struct {
	BlockedURI         string `json:"blocked-uri"`
	Disposition        string `json:"disposition"`
	DocumentURI        string `json:"document-uri"`
	EffectiveDirective string `json:"effective-directive"`
	OriginalPolicy     string `json:"original-policy"`
	Referrer           string `json:"referrer"`
	ScriptSample       string `json:"script-sample"`
	StatusCode         int    `json:"status-code"`
	ViolatedDirective  string `json:"violated-directive"`
}

type Reporter interface {
	Report(Report)
}

type ReporterFunc func(Report)

func (rf ReporterFunc) Report(r Report) {
	rf(r)
}

type reporter struct {
	Reporter
}

func NewReporter(r Reporter) http.Handler {
	return reporter{r}
}

func (r reporter) ServeHTTP(w http.ResponseWriter, hr *http.Request) {
	var rp report
	err := json.NewDecoder(hr.Body).Decode(*rp)
	hr.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	go r.Report(rp)
}
