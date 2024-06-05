package main

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"text/template"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/vitorconti/pdf-threading/seed"
)

func main() {
	var wg sync.WaitGroup
	// Parse the HTML template
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatal(err)
	}

	rs, err := seed.RetriveAirlineDelays()
	if err != nil {
		log.Fatal(err)
	}

	for rs.Next() {
		var airline seed.AirlineDelay
		err = rs.Scan(&airline.Id, &airline.Airport_code, &airline.Airport_name, &airline.Time_label, &airline.Delay_late, &airline.Delay_security)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go func(airline seed.AirlineDelay) {
			defer wg.Done() // Ensure Done is called even if there's an error

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, airline); err != nil {
				log.Println(err)
				return
			}

			pdfg, err := wkhtmltopdf.NewPDFGenerator()
			if err != nil {
				log.Println(err)
				return
			}

			pdfg.AddPage(wkhtmltopdf.NewPageReader(&buf))

			if err := pdfg.Create(); err != nil {
				log.Println(err)
				return
			}

			if err := pdfg.WriteFile(fmt.Sprintf("report%v.pdf", airline.Id)); err != nil {
				log.Println(err)
				return
			}

			log.Println("PDF generated successfully for airline ID:", airline.Id)
		}(airline)
	}

	wg.Wait() // Wait for all goroutines to finish
}
