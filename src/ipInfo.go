package src

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type IPInfo struct {
	IP string `json:"ip"`
	Country string `json:"country"`
	City string `json:"city"`
	Asn string `json:"asn"`
	Organization string `json:"organization"`
	Language string `json:"language"`
}

func NewIPInfo() *IPInfo {
	return &IPInfo{
		IP: "",
		Country:      "unknown",
		City:         "unknown",
		Asn:          "unknown",
		Organization: "unknown",
		Language: "en",
	}
}

func (info *IPInfo) SetLanguage(language string) {
	info.Language = language
}

func (info *IPInfo) SetIP(ip string) {
	info.IP = ip
}

func (info *IPInfo) SetCountry(location *Location) {
	country := location.Country.Names[info.Language]
	if country == "" {
		return
	} else {
		info.Country = country
	}
}

func (info *IPInfo) SetCity(location *Location) {
	city := location.City.Names[info.Language]
	if city == "" {
		return
	} else {
		info.City = city
	}
}

func (info *IPInfo) SetAsn(asn *Asn) {
	asnNumber := asn.Number
	if asnNumber == 0 {
		return
	} else {
		info.Asn = "ASN" + strconv.Itoa(int(asnNumber))
	}
}

func (info *IPInfo) SetOrganization(asn *Asn) {
	org := asn.Organization
	if org == "" {
		return
	} else {
		info.Organization = org
	}
}

func (info *IPInfo) ToJson() ([]byte, error) {
	jsonData, err := json.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal IPInfo to JSON: %q",err)
	}
	return jsonData, nil
}

func GetInfo(ip string, location *Location,asn *Asn, lang string) *IPInfo {
	info := NewIPInfo()
	info.SetIP(ip)
	info.SetLanguage(lang)
	info.SetCountry(location)
	info.SetCity(location)
	info.SetAsn(asn)
	info.SetOrganization(asn)
	return info
}