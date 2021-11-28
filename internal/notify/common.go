package notify

import (
	"bytes"
	"text/template"

	log "github.com/sirupsen/logrus"
)

func buildTemplate(currentIP, domain string, tplsrc string) string {
	t := template.New("notification template")
	if _, err := t.Parse(tplsrc); err != nil {
		log.Error("Failed to parse template:", err)
		return ""
	}

	data := struct {
		CurrentIP string
		Domain    string
	}{
		currentIP,
		domain,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Error(err)
		return ""
	}

	return tpl.String()
}
