package wallet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccessKey(t *testing.T) {
	pubKey, privKey, err := CreateAccessKey()
	assert.NoError(t, err)
	fmt.Println(pubKey)
	fmt.Println(privKey)
	t.Fail()
}
