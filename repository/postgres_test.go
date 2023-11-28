package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock implementation of your DB struct
type MockDB struct {
	mock.Mock
	*gorm.DB
}

func TestFindPaymentByReference(t *testing.T) {
	// Create a new instance of the mock DB
	mockDB := new(MockDB)

	// Create an instance of the struct that you expect to receive from the database
	expectedPayment := &models.Payment{
		// Initialize with your expected values
	}

	// Set expectations for the FindPaymentByReference method
	mockDB.On("FindPaymentByReference", mock.Anything).Return(expectedPayment, nil)

	// Call your method with the mockDB
	db := &DB{Db: mockDB.DB} // Assuming DB is the struct containing the FindPaymentByReference method
	result, err := db.FindPaymentByReference(uuid.New())

	// Assert that the expectations were met
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedPayment, result)

	// Assert that the expected methods were called
	mockDB.AssertExpectations(t)
}
