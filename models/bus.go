package models

type Bus struct {
	Linea        string `json:"linea"`
	Destinazione string `json:"destinazione"`
	Orario       string `json:"orario"`
	Stato        string `json:"stato"`
	Mezzo        string `json:"mezzo"`
	Soppressa    bool   `json:"soppressa"`
}
