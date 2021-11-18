package ginney

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockStruct struct {
	data string
}

func TestRepository_NewSuccessResponse(t *testing.T) {
	t.Run("Happy", func(t *testing.T) {
		fakeData := mockStruct{"test"}

		newSuccessResponse := NewSuccessResponse(fakeData)
		assert.Equal(t, newSuccessResponse.Status, StatusSuccess)
		assert.Equal(t, newSuccessResponse.Message, MessageOk)
		assert.Equal(t, newSuccessResponse.Data, fakeData)

	})
}

func TestRepository_NewErrorResponse(t *testing.T) {
	t.Run("Happy", func(t *testing.T) {
		mockErrorMessage := "Error here!!"

		newSuccessResponse := NewErrorResponse(mockErrorMessage)
		assert.Equal(t, newSuccessResponse.Status, StatusFail)
		assert.Equal(t, newSuccessResponse.Message, mockErrorMessage)

	})
}
