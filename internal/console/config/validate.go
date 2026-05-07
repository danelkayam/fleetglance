package config

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func ValidateFleet(fleet *Fleet) error {
	if fleet == nil {
		return errors.New("fleet config is required")
	}

	if fleet.Version != 1 {
		return fmt.Errorf("unsupported fleet config version: %d", fleet.Version)
	}

	if len(fleet.Ships) == 0 {
		return errors.New("fleet must contain at least one ship")
	}
	if len(fleet.Ships) > MaxShips {
		return fmt.Errorf("fleet supports at most %d ships", MaxShips)
	}

	shipNames := make([]string, 0, len(fleet.Ships))
	for name := range fleet.Ships {
		shipNames = append(shipNames, name)
	}
	sort.Strings(shipNames)

	for _, name := range shipNames {
		ship := fleet.Ships[name]
		if strings.TrimSpace(ship.URL) == "" {
			return fmt.Errorf("ship %q url is required", name)
		}

		parsedURL, err := url.Parse(ship.URL)
		if err != nil || parsedURL.Host == "" || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
			return fmt.Errorf("ship %q url must be absolute http/https URL", name)
		}
		if parsedURL.Path != "" && parsedURL.Path != "/" {
			return fmt.Errorf("ship %q url must be agent base URL without path", name)
		}
	}

	return nil
}
