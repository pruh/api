package networks

import (
	"errors"
	"math"
)

const DISABLED = -1

// Returns upload rate limit in KBPS from omada SSID data response
// Returns -1 if no limit is set
func (d *Data) ToUploadRateLimit() (*int, error) {
	if d == nil || d.RateLimit == nil || d.RateLimit.UpLimitEnable == nil {
		return nil, errors.New("upload rate limit enabled is not set for data")
	}

	if !*d.RateLimit.UpLimitEnable {
		return NewInt(DISABLED), nil
	}

	if d.RateLimit.UpLimit == nil || d.RateLimit.UpLimitType == nil {
		return nil, errors.New("upload rate limit is not set for data")
	}

	return NewInt(*d.RateLimit.UpLimit * int(powInt(1024, *d.RateLimit.UpLimitType))), nil
}

// Returns download rate limit in KBPS from omada SSID data response
// Returns -1 if no limit is set
func (d *Data) ToDownloadRateLimit() (*int, error) {
	if d == nil || d.RateLimit == nil || d.RateLimit.DownLimitEnable == nil {
		return nil, errors.New("download rate limit enabled is not set for data")
	}

	if !*d.RateLimit.DownLimitEnable {
		return NewInt(DISABLED), nil
	}

	if d.RateLimit.DownLimit == nil || d.RateLimit.DownLimitType == nil {
		return nil, errors.New("download rate limit is not set for data")
	}

	return NewInt(*d.RateLimit.DownLimit * powInt(1024, *d.RateLimit.DownLimitType)), nil
}

func IsSpeedLimitEqual(
	speedLimitEnable *bool,
	speedLimit *int,
	speedLimitType *int,
	requestSpeedLimit *int) bool {
	if requestSpeedLimit == nil {
		// no speed limit in request
		return true
	}

	if *requestSpeedLimit < 1 {
		// request to set no speed limit
		// speed is equal if speed limit not set
		return !*speedLimitEnable
	}

	// request to set speed limit

	if !*speedLimitEnable {
		// speed is NOT equal if speed limit is NOT enabled
		return false
	}

	// speed is equal if speed limit speed is the same
	speedLimitKbps := *speedLimit * powInt(1024, *speedLimitType)
	return speedLimitKbps == *requestSpeedLimit
}

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
