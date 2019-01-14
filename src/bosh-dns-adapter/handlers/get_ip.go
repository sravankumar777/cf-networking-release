package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"code.cloudfoundry.org/lager"
	"golang.org/x/net/dns/dnsmessage"
)

//go:generate counterfeiter -o fakes/sdc_client.go --fake-name SDCClient . sdcClient
type sdcClient interface {
	IPs(hostname string) ([]string, error)
}

type GetIP struct {
	SDCClient sdcClient
	Logger    lager.Logger
}

func (g GetIP) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	name := getQueryParam(req, "name", "")
	dnsType := getQueryParam(req, "type", "1")

	requestLogger := g.Logger.Session("serve-request")

	if dnsType != "1" {
		g.writeResponse(w, dnsmessage.RCodeSuccess, name, dnsType, nil)
		requestLogger.Debug("unsupported record type", lager.Data{
			"ips":          "",
			"service-name": name,
		})
		return
	}

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		g.writeResponse(w, dnsmessage.RCodeServerFailure, name, dnsType, nil)
		requestLogger.Debug("name parameter empty", lager.Data{
			"ips":          "",
			"service-name": "",
		})
		return
	}

	ips, err := g.SDCClient.IPs(name)
	if err != nil {
		g.writeErrorResponse(w, err)
		return
	}

	g.writeResponse(w, dnsmessage.RCodeSuccess, name, dnsType, ips)
}
func (g GetIP) writeErrorResponse(resp http.ResponseWriter, err error) {
	resp.WriteHeader(http.StatusInternalServerError)
	_, err = resp.Write([]byte(err.Error()))
	if err != nil {
		g.Logger.Error("Error writing to http response body", err)
	}
}

func getQueryParam(req *http.Request, key, defaultValue string) string {
	queryValue := req.URL.Query().Get(key)
	if queryValue == "" {
		return defaultValue
	}

	return queryValue
}

type Answer struct {
	Name   string `json:"name"`
	RRType uint16 `json:"type"`
	TTL    uint32 `json:"TTL"`
	Data   string `json:"data"`
}

func (g GetIP) writeResponse(resp http.ResponseWriter, dnsResponseStatus dnsmessage.RCode, requestedInfraName string, dnsType string, ips []string) {
	responseBody, err := buildResponseBody(dnsResponseStatus, requestedInfraName, dnsType, ips)
	if err != nil {
		g.Logger.Error("Error building response", err)
		return
	}

	_, err = resp.Write([]byte(responseBody))
	if err != nil {
		g.Logger.Error("Error writing to http response body", err)
	}

	g.Logger.Debug("HTTPServer access")
}

func buildResponseBody(dnsResponseStatus dnsmessage.RCode, requestedInfraName string, dnsType string, ips []string) (string, error) {
	answers := make([]Answer, len(ips), len(ips))
	for i, ip := range ips {
		answers[i] = Answer{
			Name:   requestedInfraName,
			RRType: uint16(dnsmessage.TypeA),
			Data:   ip,
			TTL:    0,
		}
	}

	bytes, err := json.Marshal(answers)
	if err != nil {
		return "", err // not tested
	}

	template := `{
		"Status": %d,
		"TC": false,
		"RD": false,
		"RA": false,
		"AD": false,
		"CD": false,
		"Question":
		[
			{
				"name": "%s",
				"type": %s
			}
		],
		"Answer": %s,
		"Additional": [ ],
		"edns_client_subnet": "0.0.0.0/0"
	}`

	return fmt.Sprintf(template, dnsResponseStatus, requestedInfraName, dnsType, string(bytes)), nil
}
