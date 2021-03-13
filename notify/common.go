package notify

import (
	"bytes"
	"log"
	"text/template"
)

func buildTemplate(currentIP, domain string, tplsrc string) string {
	t := template.New("notification template")
	if _, err := t.Parse(tplsrc); err != nil {
		log.Println("Failed to parse template")
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
		log.Println(err.Error())
		return ""
	}

	return tpl.String()
}
