package models

type FermataResponse struct {
	Fermata string `json:"fermata"`
	Bus     []Bus  `json:"bus"`
}
