// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// File represents a bootstrap registry file (i.e. {asn,dns,ipv4,ipv6}.json).
type File struct {
	// Fields from the JSON document.
	Description string
	Publication string
	Version     string

	// Map of service entries to RDAP base URLs.
	//
	// e.g. in ipv6.json, the following mapping:
	// "2c00::/12" => https://rdap.afrinic.net/rdap/,
	//                http://rdap.afrinic.net/rdap/.
	Entries map[string][]*url.URL

	// The file's JSON document.
	JSON []byte
}

// NewFile constructs a File from a bootstrap registry file.
func NewFile(jsonDocument []byte, serviceOverride map[string]string) (*File, error) {
	var doc struct {
		Description string
		Publication string
		Version     string

		Services [][][]string
	}

	err := json.Unmarshal(jsonDocument, &doc)
	if err != nil {
		return nil, err
	}

	f := &File{}
	f.Description = doc.Description
	f.Publication = doc.Publication
	f.Version = doc.Version
	f.JSON = jsonDocument

	f.Entries = make(map[string][]*url.URL)

	for _, s := range doc.Services {
		if len(s) != 2 {
			return nil, errors.New("Malformed bootstrap (bad services array)")
		}

		entries := s[0]
		rawURLs := s[1]

		var urls []*url.URL

		for _, rawURL := range rawURLs {
			url, err := url.Parse(rawURL)

			// Ignore unparsable URLs.
			if err != nil {
				continue
			}

			urls = append(urls, url)
		}

		if len(urls) > 0 {
			for _, entry := range entries {
				if override, ok := hasOverride(entry, serviceOverride); ok {
					oURL, err := url.Parse(override)
					if err != nil {
						return nil, errors.New(fmt.Sprintf("Service override cannot be parsed as a URL: %s:%s", entry, override))
					}

					urls = []*url.URL{oURL}
				}
				f.Entries[entry] = urls
			}
		}
	}

	return f, nil
}

func hasOverride(entry string, overrides map[string]string) (string, bool) {
	override, ok := overrides[strings.ToLower(entry)]
	return override, ok
}
