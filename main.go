package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type Location struct {
	Name                string  `json:"name" yaml:"name"`
	DisplayName         string  `json:"displayName" yaml:"displayName"`
	ShortName           *string `yaml:"shortName"`
	RegionalDisplayName string  `json:"regionalDisplayName" yaml:"regionalDisplayName"`
	PairedRegionName    *string `json:"pairedRegionName" yaml:"pairedRegionName"`
}

type LocationData struct {
	Locations []Location `yaml:"locations"`
}

type LocationFetcher interface {
	FetchLocations(context.Context) ([]Location, error)
}

type FileHandler interface {
	LoadExistingShortNames(string) (map[string]*string, error)
	WriteLocations(string, LocationData) error
}

type CLIFetcher struct{}

func (f *CLIFetcher) FetchLocations(context.Context) ([]Location, error) {
	cmd := exec.Command("az", "account", "list-locations", "--query",
		"[].{regionalDisplayName: regionalDisplayName, name: name, displayName: displayName, pairedRegionName: metadata.pairedRegion[0].name}")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("azure CLI failed: %w", err)
	}

	var locations []Location
	if err := json.Unmarshal(output, &locations); err != nil {
		return nil, fmt.Errorf("failed to parse CLI output: %w", err)
	}

	return locations, nil
}

type YAMLHandler struct{}

func (h *YAMLHandler) LoadExistingShortNames(filename string) (map[string]*string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return make(map[string]*string), nil
	}

	var existing LocationData
	if err := yaml.Unmarshal(data, &existing); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	shortNames := make(map[string]*string)
	for _, loc := range existing.Locations {
		if loc.ShortName != nil && *loc.ShortName != "" {
			shortNames[loc.Name] = loc.ShortName
		}
	}

	return shortNames, nil
}

func (h *YAMLHandler) WriteLocations(filename string, data LocationData) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := os.WriteFile(filename, yamlData, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	return nil
}

type LocationUpdater struct {
	fetcher     LocationFetcher
	fileHandler FileHandler
	outputFile  string
}

func (u *LocationUpdater) UpdateLocations(ctx context.Context) error {
	locations, err := u.fetcher.FetchLocations(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %w", err)
	}

	existingShortNames, err := u.fileHandler.LoadExistingShortNames(u.outputFile)
	if err != nil {
		return fmt.Errorf("failed to load existing short names: %w", err)
	}

	for i := range locations {
		if existing, ok := existingShortNames[locations[i].Name]; ok {
			locations[i].ShortName = existing
		}
	}

	if err := u.fileHandler.WriteLocations(u.outputFile, LocationData{Locations: locations}); err != nil {
		return fmt.Errorf("failed to write locations: %w", err)
	}

	fmt.Printf("Successfully updated %s with %d locations\n", u.outputFile, len(locations))
	return nil
}

func main() {
	updater := &LocationUpdater{
		fetcher:     &CLIFetcher{},
		fileHandler: &YAMLHandler{},
		outputFile:  "locations.yaml",
	}

	if err := updater.UpdateLocations(context.Background()); err != nil {
		log.Fatal(err)
	}
}
