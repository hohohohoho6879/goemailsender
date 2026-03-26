package template

import (
	"bytes"
	"embed"
	htmltemplate "html/template"
	"log"
)

//go:embed email.html
var templateFS embed.FS

var emailTmpl *htmltemplate.Template

func init() {
	var err error
	emailTmpl, err = htmltemplate.ParseFS(templateFS, "email.html")
	if err != nil {
		log.Fatalf("Failed to parse email template: %s", err)
	}
}

type EmailData struct {
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Code        string `json:"code"`
	MagicUrl    string `json:"magicUrl"`
	WebsiteName string `json:"websiteName"`
	WebsiteUrl  string `json:"websiteUrl"`
}

func Render(data EmailData) (string, error) {
	var buf bytes.Buffer
	if err := emailTmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
