package wallet

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_daysToNanoSeconds(t *testing.T) {
	duration := daysToNanoSeconds(28)
	fmt.Println(duration.Nanoseconds())
	assert.Equal(t, time.Duration(2419200000000000), duration)
}
