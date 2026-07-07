package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
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

type cacheEntry struct {
	Data      []fermataRaw
	ExpiresAt time.Time
}

var (
	bacinoCache = make(map[string]cacheEntry)
	cacheMutex  sync.RWMutex
)

func GetBacino(selected string) ([]fermataRaw, error) {

	if selected != "ra" && selected != "rn" && selected != "fc" {
		return nil, errors.New("bacino non valido")
	}

	cacheMutex.RLock()
	entry, ok := bacinoCache[selected]
	cacheMutex.RUnlock()

	if ok && time.Now().Before(entry.ExpiresAt) {
		return entry.Data, nil
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

	cacheMutex.Lock()
	bacinoCache[selected] = cacheEntry{
		Data:      wrapper.D,
		ExpiresAt: time.Now().Add(12 * time.Hour),
	}
	cacheMutex.Unlock()

	return wrapper.D, nil
}
