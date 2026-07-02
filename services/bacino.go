package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type bacinoRequest struct {
	Bacino string `json:"bacino"`
}

type fermataRaw struct {
	Nome     string `json:"nome"`
	Palina   string `json:"palina"`
	TargetID string `json:"targetID"`
}

func GetBacino(selected string) ([]fermataRaw, error) {

	if selected != "ra" && selected != "rn" && selected != "fc" {
		return nil, errors.New("bacino non valido")
	}

	body, _ := json.Marshal(bacinoRequest{
		Bacino: selected,
	})

	client := &http.Client{
		Timeout: 8 * time.Second,
	}

	resp, err := client.Post(
		"https://infobus.startromagna.it/FermateService.asmx/GetFermateByBacino",
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wrapper struct {
		D []fermataRaw `json:"d"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	return wrapper.D, nil
}
