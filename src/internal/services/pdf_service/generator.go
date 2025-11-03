package pdfservice

import (
	"encoding/json"
	"log/slog"
	"logispro/internal/db"
	"logispro/internal/utils"
	"os"
	"strconv"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	amqp "github.com/rabbitmq/amqp091-go"
)

type InvoiceTemplateData struct {
	Title       string
	AgencyName  string
	Description string
	PaymentID   string
	Users       []db.User
	Amount      float64
}

type Generator struct {
	TemplatesDir string
	StorageDir   string
	Logger       *slog.Logger
	RabbitMqConn *amqp.Connection
}

type GeneratedInvoice struct {
	Path  string
	Users []int64
	Id    string
}

func (g Generator) Generate(
	out string,
	template string,
	data InvoiceTemplateData,
) {
	go func(templateName string, data InvoiceTemplateData) {

		paymentId := data.PaymentID
		usersSlice := utils.ExtractField(data.Users, func(u db.User) int64 {
			return u.ID
		})
		b, err := utils.ParseTemplate(g.TemplatesDir+"/"+templateName, data)
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
		err = pdfg.WriteFile(g.StorageDir + "/invoices/" + out)
		if err != nil {
			g.Logger.Error("failed to pdfg.WriteFile()", "error", err)
			return
		}

		invoice := GeneratedInvoice{Path: g.StorageDir + "/invoices/" + out, Users: usersSlice, Id: paymentId}
		// push to rabitmq
		ch, err := g.RabbitMqConn.Channel()
		if err != nil {
			g.Logger.Error("failed to get rabbitmq channel", "error", err)
			return
		}
		defer ch.Close()
		rmq := &utils.RabbitMQ{Conn: g.RabbitMqConn, Channel: ch}
		dataMq, err := json.Marshal(invoice)
		if err != nil {
			g.Logger.Error("failed to get marshal payload", "error", err)
			return
		}
		err = rmq.Publish("created_invoices", dataMq, amqp.Table{"x-retry": 1})
		if err != nil {
			// retry ? maybe
			g.Logger.Error("failed to publish to rabbitmq", "error", err)
			return
		}
	}(template, data)
}
