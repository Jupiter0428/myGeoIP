package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"myGeoIP/src"
	"net/http"
)


func main() {
	config, err := src.GetConfig("./conf/dbConfig.yaml")
	if err != nil {
		log.Fatal(err)
	}

	geoipCity, err := src.NewGeoIP(config.City.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer geoipCity.Close()

	geoipASN, err := src.NewGeoIP(config.Asn.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer geoipASN.Close()

	scheduled := cron.New()
	scheduled.AddFunc("@daily", func() {src.CheckGeoIPVersion(&geoipCity, config)})
	scheduled.AddFunc("@daily", func() {src.CheckGeoIPVersion(&geoipASN, config)})

	r := gin.Default()

	r.POST("/geoip", func(c *gin.Context) {
		var requestData struct {
			IPAddr string `json:"ipaddr"`
			Lang   string `json:"lang"`
		}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		ip, err := src.CheckIP(requestData.IPAddr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IP address"})
			return
		}

		lang := src.CheckLang(requestData.Lang, config.City.SupportedLanguages, config.City.DefaultLanguage)

		location, err := geoipCity.LookupLocation(ip)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to lookup location: %s", err)})
			return
		}

		asn, err:= geoipASN.LookupASN(ip)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to lookup ASN: %s", err)})
			return
		}

		ipInfo := src.GetInfo(requestData.IPAddr, location, asn, lang)

		if ipInfo != nil {
			c.JSON(http.StatusOK, gin.H{
				"ip":      ipInfo.IP,
				"country": ipInfo.Country,
				"city":    ipInfo.City,
				"asn":     ipInfo.Asn,
				"org":     ipInfo.Organization,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get IP information"})
		}
	})

	r.Run(":5000")
}