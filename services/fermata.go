package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Daniongithub/startfermate-api/models"
	"github.com/Daniongithub/startfermate-api/utils"

	"github.com/PuerkitoBio/goquery"
)

func GetFermata(param, param2, palina, det string) ([]any, error) {

	url := fmt.Sprintf(
		"https://infobus.startromagna.it/InfoFermata?param=%s&param2=%s&palina=%s",
		param, param2, palina,
	)

	client := &http.Client{
		Timeout: 8 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var result models.FermataResponse
	result.Bus = []models.Bus{}

	result.Fermata = strings.TrimSpace(
		doc.Find("h2.fw-bold.text-primary.title").Text(),
	)

	doc.Find(".container.mb-50 .bus-card").Each(func(i int, el *goquery.Selection) {

		isSopp := el.Find(".bus-status.sopp").Length() > 0

		orario := strings.TrimSpace(el.Find(".bus-times span").First().Text())
		stato := strings.TrimSpace(el.Find(".bus-status").Text())

		headerSpan := el.Find(".bus-header > span").First()

		// rimuove l'icona material-icons
		headerSpan.Find(".material-icons").Remove()

		linea := strings.TrimSpace(headerSpan.Text())
		linea = strings.TrimPrefix(linea, "Linea ")

		dest := strings.TrimSpace(el.Find(".bus-destination").Text())
		mezzo, _ := el.Find(".det a").Attr("data-vehicle")

		if len(mezzo) == 4 {
			mezzo = "3" + mezzo
		}

		if stato == "Non disp" {
			stato = "N.D."
		}

		if mezzo == "" {
			mezzo = "N.D."
		}

		if dest == "Fornace.Zarattini" {
			dest = "Fornace Zarattini"
		}

		var bus models.Bus

		utils.Normalize(&bus, linea, dest)

		bus.Orario = orario
		bus.Stato = stato
		bus.Mezzo = mezzo
		bus.Soppressa = isSopp

		if utils.ShouldShow(det, bus.Destinazione) {
			result.Bus = append(result.Bus, bus)
		}
	})

	var out []any

	// primo elemento: fermata (come prima)
	out = append(out, map[string]string{
		"fermata": result.Fermata,
	})

	// poi i bus
	for _, b := range result.Bus {
		out = append(out, b)
	}

	return out, nil
}
