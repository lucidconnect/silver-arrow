package merchant

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIdParser(t *testing.T) {
	uid := uuid.New()
	fmt.Println(uid)
	encodedId, err := EncodeUUIDToMerchantId(uid)
	assert.NoError(t, err)
	assert.NotEmpty(t, encodedId)
	fmt.Println(encodedId)
	// parsed := parseMerchantIdtoUUID(encodedId)
	// assert.Equal(t, uid, parsed)
}
