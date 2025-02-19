// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/openrdap/rdap/test"
)

func TestNetRegistryLookupsDNSNested(t *testing.T) {
	test.Start(test.BootstrapComplex)
	defer test.Finish()

	var bytes []byte = test.Get("https://rdap.example.org/dns.json")

	var d *DNSRegistry
	d, err := NewDNSRegistry(bytes, nil)

	if err != nil {
		t.Fatal(err)
	}

	tests := []registryTest{
		{
			"",
			false,
			"",
			[]string{"https://example.root", "http://example.root"},
		},
		{
			"example.com",
			false,
			"com",
			[]string{"https://example.com", "http://example.com"},
		},
		{
			"sub.example.com",
			false,
			"sub.example.com",
			[]string{"https://example.com/sub", "http://example.com/sub"},
		},
		{
			"sub.sub.example.com",
			false,
			"sub.example.com",
			[]string{"https://example.com/sub", "http://example.com/sub"},
		},
		{
			"example.xyz",
			false,
			"",
			[]string{"https://example.root", "http://example.root"},
		},
	}

	runRegistryTests(t, tests, d)
}

func TestNetRegistryLookupsDNS(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	var bytes []byte = test.Get("https://data.iana.org/rdap/dns.json")

	var d *DNSRegistry
	d, err := NewDNSRegistry(bytes, nil)

	if err != nil {
		t.Fatal(err)
	}

	tests := []registryTest{
		{
			"",
			false,
			"",
			[]string{},
		},
		{
			"www.EXAMPLE.BR",
			false,
			"br",
			[]string{"https://rdap.registro.br/"},
		},
		{
			"example.xyz",
			false,
			"",
			[]string{},
		},
	}

	runRegistryTests(t, tests, d)
}

func TestNewDNSRegistry(t *testing.T) {
	type args struct {
		json            []byte
		serviceOverride map[string]*url.URL
	}
	tests := []struct {
		name    string
		args    args
		want    *DNSRegistry
		wantErr bool
	}{
		{
			name: "replaces server with rdap server override",
			args: args{
				json: test.LoadFile("../testdata/bootstrap/testingdns.json"),
				serviceOverride: map[string]*url.URL{
					"ar": {
						Scheme: "https",
						Host:   "rdap.testingserveroverride.ar",
						Path:   "/",
					},
				},
			},
			want: &DNSRegistry{
				dns: map[string][]*url.URL{
					"ar": {
						{
							Scheme: "https",
							Host:   "rdap.testingserveroverride.ar",
							Path:   "/",
						},
					},
				},

				file: &File{
					Description: "RDAP bootstrap file for Domain Name System registrations",
					Publication: "2017-03-15T21:26:24Z",
					Version:     "1.0",
					Entries: map[string][]*url.URL{
						"ar": {
							{
								Scheme: "https",
								Host:   "rdap.testingserveroverride.ar",
								Path:   "/",
							},
						},
					},
					JSON: test.LoadFile("../testdata/bootstrap/testingdns.json"),
				},
			},
			wantErr: false,
		},
		{
			name: "replaces server with rdap server override",
			args: args{
				json: test.LoadFile("../testdata/bootstrap/testingdnsmultipletld.json"),
				serviceOverride: map[string]*url.URL{
					"ai": {
						Scheme: "https",
						Host:   "rdap.nic.ai",
						Path:   "/",
					},
				},
			},
			want: &DNSRegistry{
				dns: map[string][]*url.URL{
					"ar": {
						{
							Scheme: "https",
							Host:   "rdap.nic.ar",
							Path:   "/",
						},
					},
					"ai": {
						{
							Scheme: "https",
							Host:   "rdap.nic.ai",
							Path:   "/",
						},
					},
					"com": {
						{
							Scheme: "https",
							Host:   "rdap.nic.com",
							Path:   "/",
						},
					},
					"test": {
						{
							Scheme: "https",
							Host:   "rdap.nic.com",
							Path:   "/",
						},
					},
				},

				file: &File{
					Description: "RDAP bootstrap file for Domain Name System registrations",
					Publication: "2017-03-15T21:26:24Z",
					Version:     "1.0",
					Entries: map[string][]*url.URL{
						"ar": {
							{
								Scheme: "https",
								Host:   "rdap.nic.ar",
								Path:   "/",
							},
						},
						"ai": {
							{
								Scheme: "https",
								Host:   "rdap.nic.ai",
								Path:   "/",
							},
						},
						"com": {
							{
								Scheme: "https",
								Host:   "rdap.nic.com",
								Path:   "/",
							},
						},
						"test": {
							{
								Scheme: "https",
								Host:   "rdap.nic.com",
								Path:   "/",
							},
						},
					},
					JSON: test.LoadFile("../testdata/bootstrap/testingdnsmultipletld.json"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Start(test.Bootstrap)
			defer test.Finish()

			got, err := NewDNSRegistry(tt.args.json, tt.args.serviceOverride)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDNSRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDNSRegistry() deepequal got = %v, want %v", got, tt.want)
			}
		})
	}
}

func createTestURLSlice(u string) ([]*url.URL, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	urlSlice := []*url.URL{parsedURL}

	return urlSlice, nil
}
