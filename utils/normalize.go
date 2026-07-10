package utils

import "github.com/Daniongithub/startfermate-api/models"

var variants = map[string]string{
	"4|Mirabilandia":                     "4B",
	"4|Classe Cantoniera":                "4C",
	"4|ClasseCantoniera":                 "4C",
	"4|Lido di Dante":                    "4D",
	"4|Fosso Ghiaia":                     "4F",
	"4|Fosso Ghiaia-Standiana":           "4F",
	"4|Savio di Cervia":                  "4F",
	"4|Classe via Liburna":               "4R",
	"4|Classe Romea Vecchia":             "4R",
	"4|R.Vecchia":                        "4R",
	"4|Romea Vecchia via del Centurione": "4C",
	"4|Classe":                           "4R",

	"1|Borgo Nuovo":        "1B",
	"3|via Sant'Alberto":   "3/",
	"1|Antica Milizia":     "1/",
	"1|Via Antica Milizia": "1/",
}

var navetto = map[string]struct {
	linea string
	dest  string
}{
	"Navetto Mare Marina": {
		linea: "65",
		dest:  "Navetto Mare Marina di Ravenna",
	},
	"Navetto Mare Punta": {
		linea: "66",
		dest:  "Navetto Mare Punta Marina",
	},
	"Navetto Mare Punta Marina-Campeggi": {
		linea: "67",
		dest:  "Navetto Mare Punta M.-Campeggi",
	},
}

var ignoredDest = map[string]bool{
	"RIENTRO DEPOSITO": true,
	"PRESA SERVIZIO":   true,
	"FUORI LINEA":      true,
}

var truncatedLines = map[string]bool{
	"1":  true,
	"1B": true,
	"3":  true,
	"4":  true,
	"4B": true,
	"4C": true,
	"4D": true,
	"4F": true,
	"4R": true,
	"5":  true,
	"8":  true,
	"18": true,
	"70": true,
	"80": true,
}

var truncatedDestinations = map[string]bool{
	"Stazione FS":        true,
	"STAZIONE FS":        true,
	"Ravenna FS":         true,
	"Ravenna Radio Taxi": true,
}

func Normalize(bus *models.Bus, linea, dest string) {

	// NAVETTI
	if v, ok := navetto[linea]; ok {
		bus.Linea = v.linea
		bus.Destinazione = v.dest
		return
	}

	key := linea + "|" + dest
	if v, ok := variants[key]; ok {
		bus.Linea = v
	} else {
		bus.Linea = linea
	}

	bus.Destinazione = dest

	if truncatedLines[bus.Linea] && truncatedDestinations[bus.Destinazione] {
		bus.Linea += "/"
	}
}

func ShouldShow(det string, dest string) bool {
	if det == "true" {
		return true
	}
	return !ignoredDest[dest]
}
