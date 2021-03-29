package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/fantashley/vaxalert/pkg/vaxspotter"
)

const VaxClientV0Path = "/api/v0"

type VaxClientV0 struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewVaxClientV0(apiURL string, httpClient *http.Client) (*VaxClientV0, error) {
	vaxSpotterURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %w", err)
	}
	vaxSpotterURL.Path = path.Join(vaxSpotterURL.Path, VaxClientV0Path)

	vaxHTTPClient := httpClient
	if vaxHTTPClient == nil {
		vaxHTTPClient = &http.Client{}
	}

	return &VaxClientV0{
		baseURL:    vaxSpotterURL,
		httpClient: vaxHTTPClient,
	}, nil
}

func (c *VaxClientV0) GetStates() ([]vaxspotter.State, error) {
	const statePath = "states.json"
	var states []vaxspotter.State

	stateURL := *c.baseURL
	stateURL.Path = path.Join(stateURL.Path, statePath)

	resp, err := c.httpClient.Get(stateURL.String())
	if err != nil {
		return states, fmt.Errorf("failed to get provider list: %w", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&states); err != nil {
		return states, fmt.Errorf("failed to decode provider list: %w", err)
	}

	return states, nil
}

func (c *VaxClientV0) GetLocations(stateCode string) (vaxspotter.Locations, error) {
	const locPath = "/states/%s.json"
	var locations vaxspotter.Locations

	locURL := *c.baseURL
	locURL.Path = path.Join(locURL.Path, fmt.Sprintf(locPath, stateCode))

	resp, err := c.httpClient.Get(locURL.String())
	if err != nil {
		return locations, fmt.Errorf("failed to get location list: %w", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return locations, fmt.Errorf("failed to decode location list: %w", err)
	}

	return locations, nil
}
