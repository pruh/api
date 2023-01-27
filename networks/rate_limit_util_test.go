package networks_test

import (
	"testing"

	. "github.com/pruh/api/networks"
	"github.com/stretchr/testify/assert"
)

func TestToUploadRateLimit(t *testing.T) {
	testsData := []struct {
		description string
		data        *Data
		rateLimit   *int
		expectError bool
	}{
		{
			description: "enabled flag is missing",
			data:        &Data{},
			expectError: true,
		},
		{
			description: "enabled flag is missing",
			data: &Data{
				RateLimit: &RateLimit{
					UpLimitEnable: NewBool(false),
				}},
			rateLimit:   NewInt(DISABLED),
			expectError: false,
		},
		{
			description: "enabled is true; rate limit missing",
			data: &Data{
				RateLimit: &RateLimit{
					UpLimitEnable: NewBool(true),
				}},
			expectError: true,
		},
		{
			description: "enabled is true; type missing",
			data: &Data{
				RateLimit: &RateLimit{
					UpLimitEnable: NewBool(true),
					UpLimit:       NewInt(1),
				}},
			expectError: true,
		},
		{
			description: "enabled is true; rate missing",
			data: &Data{
				RateLimit: &RateLimit{
					UpLimitEnable: NewBool(true),
					UpLimitType:   NewInt(1),
				}},
			expectError: true,
		},
		{
			description: "rate limit in kbps",
			data: &Data{
				RateLimit: &RateLimit{
					UpLimitEnable: NewBool(true),
					UpLimit:       NewInt(1000),
					UpLimitType:   NewInt(0),
				}},
			rateLimit:   NewInt(1000),
			expectError: false,
		},
		{
			description: "rate limit in mbps",
			data: &Data{
				RateLimit: &RateLimit{
					UpLimitEnable: NewBool(true),
					UpLimit:       NewInt(10),
					UpLimitType:   NewInt(1),
				}},
			rateLimit:   NewInt(10 * 1024),
			expectError: false,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		r, e := testData.data.ToUploadRateLimit()

		if testData.expectError {
			assert.NotNil(e, "should return error")
		} else {
			assert.Nil(e, "should not return error")
			assert.Equal(*testData.rateLimit, *r, "upload rate limit is incorrect")
		}
	}

}

func TestToDownloadRateLimit(t *testing.T) {
	testsData := []struct {
		description string
		data        *Data
		rateLimit   *int
		expectError bool
	}{
		{
			description: "enabled flag is missing",
			data:        &Data{},
			expectError: true,
		},
		{
			description: "enabled flag is missing",
			data: &Data{
				RateLimit: &RateLimit{
					DownLimitEnable: NewBool(false),
				}},
			rateLimit:   NewInt(DISABLED),
			expectError: false,
		},
		{
			description: "enabled is true; rate limit missing",
			data: &Data{
				RateLimit: &RateLimit{
					DownLimitEnable: NewBool(true),
				}},
			expectError: true,
		},
		{
			description: "enabled is true; type missing",
			data: &Data{
				RateLimit: &RateLimit{
					DownLimitEnable: NewBool(true),
					DownLimit:       NewInt(1),
				}},
			expectError: true,
		},
		{
			description: "enabled is true; rate missing",
			data: &Data{
				RateLimit: &RateLimit{
					DownLimitEnable: NewBool(true),
					DownLimitType:   NewInt(1),
				}},
			expectError: true,
		},
		{
			description: "rate limit in kbps",
			data: &Data{
				RateLimit: &RateLimit{
					DownLimitEnable: NewBool(true),
					DownLimit:       NewInt(1000),
					DownLimitType:   NewInt(0),
				}},
			rateLimit:   NewInt(1000),
			expectError: false,
		},
		{
			description: "rate limit in mbps",
			data: &Data{
				RateLimit: &RateLimit{
					DownLimitEnable: NewBool(true),
					DownLimit:       NewInt(10),
					DownLimitType:   NewInt(1),
				}},
			rateLimit:   NewInt(10 * 1024),
			expectError: false,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		r, e := testData.data.ToDownloadRateLimit()

		if testData.expectError {
			assert.NotNil(e, "should return error")
		} else {
			assert.Nil(e, "should not return error")
			assert.Equal(*testData.rateLimit, *r, "download rate limit is incorrect")
		}
	}

}

func TestIsSpeedLimitEqual(t *testing.T) {
	testsData := []struct {
		description       string
		requestSpeedLimit *int
		limitEnable       *bool
		limit             *int
		limitType         *int
		result            bool
	}{
		{
			description: "no speed limit in request",
			result:      true,
		},
		{
			description:       "request no limit, current limit on",
			requestSpeedLimit: NewInt(DISABLED),
			limitEnable:       NewBool(true),
			result:            false,
		},
		{
			description:       "request no limit, current limit off",
			requestSpeedLimit: NewInt(DISABLED),
			limitEnable:       NewBool(false),
			result:            true,
		},
		{
			description:       "request 1024 limit, current limit off",
			requestSpeedLimit: NewInt(1024),
			limitEnable:       NewBool(false),
			result:            false,
		},
		{
			description:       "request 1024 limit, current limit 2048 kbps",
			requestSpeedLimit: NewInt(1024),
			limitEnable:       NewBool(true),
			limitType:         NewInt(0),
			limit:             NewInt(2048),
			result:            false,
		},
		{
			description:       "request 1024 limit, current limit 1024 kbps",
			requestSpeedLimit: NewInt(1024),
			limitEnable:       NewBool(true),
			limitType:         NewInt(0),
			limit:             NewInt(1024),
			result:            true,
		},
		{
			description:       "request 1024 limit, current limit 2 mbps",
			requestSpeedLimit: NewInt(1024),
			limitEnable:       NewBool(true),
			limitType:         NewInt(1),
			limit:             NewInt(2),
			result:            false,
		},
		{
			description:       "request 1024 limit, current limit 1 mbps",
			requestSpeedLimit: NewInt(1024),
			limitEnable:       NewBool(true),
			limitType:         NewInt(1),
			limit:             NewInt(1),
			result:            true,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		r := IsSpeedLimitEqual(testData.limitEnable, testData.limit, testData.limitType, testData.requestSpeedLimit)

		assert.Equal(testData.result, r, "speed limit equal result missing")
	}
}
