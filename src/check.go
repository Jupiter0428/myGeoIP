package src

import (
	"fmt"
	"log"
	"net"
)

func CheckLang(language string, languages []string, defaultLang string) string {
	for _, lang := range languages {
		if lang == language {
			return language
		}
	}
	return defaultLang
}

func CheckIP(ipaddr string) ([]byte, error) {
	parsedIP := net.ParseIP(ipaddr)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid IP address: %q", ipaddr)
	}
	return parsedIP, nil
}

func CheckGeoIPVersion(mmDB **GeoIP, config *Config) {
	var path string
	if (*mmDB).DBType & isCity == 0 {
		path = config.City.Path
	} else if (*mmDB).DBType & isASN == 0 {
		path = config.Asn.Path
	}

	newVersion, err := func() (uint, error){
		geoIP, err := NewGeoIP(path)
		return geoIP.BuildEpoch, err
	} ()

	if err != nil {
		log.Println("Failed to get database version:", err)
	}

	if newVersion != (*mmDB).BuildEpoch {
		(*mmDB).Close()
		*mmDB, err = NewGeoIP(path)
		if err != nil {
			log.Println("Failed to update object mmdb:", err)
		}
		log.Println("GeoIP versions updated successfully.")
	}
}