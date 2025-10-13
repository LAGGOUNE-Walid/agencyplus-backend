package pdfservice

import (
	"bytes"
	"log"
	"log/slog"
	"logispro/internal/db"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type InvoiceTemplateData struct {
	Title       string
	AgencyName  string
	Description string
	PaymentID   string
	Usernames   []db.User
	Amount      float64
}

type Generator struct {
	TemplatesDir string
	StorageDir   string
	Logger       *slog.Logger
}

func (g Generator) Generate(
	out string,
	template string,
	data any,
) {
	go func(templateName string, data any) {
		b, err := g.ParseTemplate(g.TemplatesDir+"/"+templateName, data)
		if err != nil {
			g.Logger.Error("failed to parse template", "error", err)
			return
		}
		t := time.Now().Unix()
		err = os.WriteFile(g.StorageDir+"/invoices/"+strconv.FormatInt(int64(t), 10)+".html", b, 0644)
		if err != nil {
			g.Logger.Error("failed to write bytes to html file", "error", err)
			return
		}
		f, err := os.Open(g.StorageDir + "/invoices/" + strconv.FormatInt(int64(t), 10) + ".html")
		if f != nil {
			defer f.Close()
		}
		if err != nil {
			g.Logger.Error("failed to open html file", "error", err)
			return
		}

		pdfg, err := wkhtmltopdf.NewPDFGenerator()
		if err != nil {
			g.Logger.Error("failed to init wkhtmltopdf", "error", err)
			return
		}
		pdfg.AddPage(wkhtmltopdf.NewPageReader(f))
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
		pdfg.Dpi.Set(300)
		err = pdfg.Create()
		if err != nil {
			g.Logger.Error("failed to pdfg.Create()", "error", err)
			return
		}
		err = pdfg.WriteFile(g.StorageDir + "/invoices/" + strconv.FormatInt(int64(t), 10) + ".pdf")
		if err != nil {
			log.Fatal(err)
		}
	}(template, data)
}

func (g Generator) ParseTemplate(templateName string, data any) ([]byte, error) {
	t, err := template.ParseFiles(templateName)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
