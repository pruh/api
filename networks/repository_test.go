package networks_test

import (
	"testing"

	. "github.com/pruh/api/networks"
	. "github.com/pruh/api/networks/tests"
	"github.com/stretchr/testify/assert"
)

func TestRepoGetControllerId(t *testing.T) {
	var mockCalled = false

	repo := Repository{
		OmadaApi: &MockOmadaApi{
			MockGetControllerId: func() (*ControllerIdResponse, error) {
				resp := &ControllerIdResponse{
					ErrorCode: 0,
					Msg:       "Success.",
					Result: ControllerIdResult{
						OmadacId: "someId",
					},
				}

				mockCalled = true

				return resp, nil
			},
		},
	}

	assert := assert.New(t)

	controllerId, _ := repo.GetControllerId()

	// todo verify mock method called
	omadacId := controllerId.Result.OmadacId

	assert.True(mockCalled, "mock is not called")
	assert.Equal("someId", omadacId, "controller id is not as expected")
}
