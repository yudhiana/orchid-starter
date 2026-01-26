package common

import (
	"fmt"
	"orchid-starter/constants"
	"time"
)

func GetIndonesianTimeZone(timeZoneCode string) (*time.Location, error) {
	switch timeZoneCode {
	case constants.WIB:
		return time.LoadLocation(constants.TimeZoneAsiaJakarta)
	case constants.WITA:
		return time.LoadLocation(constants.TimeZoneAsiaMakassar)
	case constants.WIT:
		return time.LoadLocation(constants.TimeZoneAsiaJayapura)
	default:
		return nil, fmt.Errorf("unknown Indonesian time zone: %s", timeZoneCode)
	}
}
