package pokeapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	// ErrNotFound occurs when the requested pokemon doesn't exist (404).
	ErrNotFound = errors.New("Pokemon not found")
	// ErrClientSide occurs when there's any type of client side error (4xx).
	ErrClientSide = errors.New("Client side error")
	// ErrServerSide occurs when there's any type of servier side error (5xx).
	ErrServerSide = errors.New("Server side error")
)

// Client represents an HTTP client for pokeapi.co.
type Client struct {
	client *resty.Client
}

// NewClient creates new pokeapi HTTP client.
func NewClient(apiURL string, timeoutSeconds int) *Client {
	client := resty.New()
	client.SetHostURL(apiURL)
	client.SetTimeout(time.Duration(timeoutSeconds) * time.Second)
	return &Client{
		client: client,
	}
}

// GetPokemonByName gets a pokemon based on its name.
func (c *Client) GetPokemonByName(name string) (*Pokemon, error) {
	resp, err := c.client.R().SetHeader("Accept", "application/json").
		SetPathParams(map[string]string{
			"name": name,
		}).Get("/pokemon/{name}")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode() >= http.StatusBadRequest &&
		resp.StatusCode() < http.StatusInternalServerError {
		return nil, ErrClientSide
	}

	if resp.StatusCode() >= http.StatusInternalServerError {
		return nil, ErrServerSide
	}

	var pokemon Pokemon

	err = json.Unmarshal(resp.Body(), &pokemon)
	if err != nil {
		return nil, err
	}

	return &pokemon, nil
}
