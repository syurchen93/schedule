package util

import (
	"context"

	"googlemaps.github.io/maps"
)

var googleMapsClient *maps.Client
var ctx context.Context

func InitGoogleMapsClient(initCtx context.Context, apiKey string) {
	var err error
	googleMapsClient, err = maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}
	ctx = initCtx
}

func GetTimezoneByCityName(cityName string) (string, error) {
	lat, lng, err := getLatLongByCityName(cityName)
	if err != nil {
		return "", err
	}

	return GetTimezoneByLatLong(lat, lng)
}

func GetTimezoneByLatLong(lat, lng float64) (string, error) {
	timezone, err := googleMapsClient.Timezone(ctx, &maps.TimezoneRequest{
		Location: &maps.LatLng{
			Lat: lat,
			Lng: lng,
		},
	})
	if err != nil {
		return "", err
	}

	return timezone.TimeZoneID, nil
}

func getLatLongByCityName(cityName string) (float64, float64, error) {
	geocodeRequest := &maps.GeocodingRequest{
		Address: cityName,
	}

	geocodeResponse, err := googleMapsClient.Geocode(ctx, geocodeRequest)
	if err != nil {
		return 0, 0, err
	}

	if len(geocodeResponse) == 0 {
		return 0, 0, nil
	}

	return geocodeResponse[0].Geometry.Location.Lat, geocodeResponse[0].Geometry.Location.Lng, nil
}
