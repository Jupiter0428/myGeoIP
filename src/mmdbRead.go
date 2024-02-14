package src

import (
	"fmt"
	"github.com/oschwald/maxminddb-golang"
)

type GeoIP struct {
	DB *maxminddb.Reader
	DBType int
	BuildEpoch uint
}

type Location struct {
	Country struct {
		IsoCode string `maxminddb:"iso_code"`
		Names map[string]string `maxminddb:"names"`
	}	`maxminddb:"country"`
	City struct {
		Names map[string]string `maxminddb:"names"`
	}	`maxminddb:"city"`
}

type Asn struct {
	Number       uint   `maxminddb:"autonomous_system_number"`
	Organization string `maxminddb:"autonomous_system_organization"`
}

const (
	isASN = 1 << iota
	isCity
)

func NewGeoIP(dbPath string) (*GeoIP, error) {
	db ,err := maxminddb.Open(dbPath)
	if err != nil{
		return nil, err
	}
	dbType, err := getDBType(db)
	epoch := getDBEpoch(db)
	return &GeoIP{
		DB:db,
		DBType: dbType,
		BuildEpoch: epoch,
	}, nil
}

func getDBType(db *maxminddb.Reader) (int, error) {
	switch db.Metadata.DatabaseType {
	case "GeoLite2-City":
		return isCity, nil
	case "GeoLite2-ASN":
		return isASN, nil
	default:
		return 0, fmt.Errorf("dbReader does not support the %q database type",db.Metadata.DatabaseType)
	}
}

func getDBEpoch(db *maxminddb.Reader) uint{
	epoch := db.Metadata.BuildEpoch
	return epoch
}

func (g *GeoIP) Close() {
	g.DB.Close()
}

func (g *GeoIP) LookupLocation(ipaddr []byte) (*Location, error) {
	if g.DBType & isCity == 0 {
		return nil, fmt.Errorf("dbMethod does not support the %q database type",g.DB.Metadata.DatabaseType)
	}
	var location Location
	err := g.DB.Lookup(ipaddr, &location)
	if err != nil {
		return nil, err
	}
	return &location, err
}

func (g *GeoIP) LookupASN(ipaddr []byte) (*Asn, error) {
	if g.DBType & isASN == 0 {
		return nil, fmt.Errorf("dbMethod does not support the %q database type",g.DB.Metadata.DatabaseType)
	}
	var asn Asn
	err := g.DB.Lookup(ipaddr, &asn)
	if err != nil {
		return nil, err
	}
	return &asn, err
}