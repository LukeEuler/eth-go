package eg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_luckyLength(t *testing.T) {
	l := luckyPreLen("9999999976ee517dc82b1efbeaa96d0ca400944c")
	assert.Equal(t, 8, l)

	l = luckyPreLen("76ee517dc82b1efbeaa96d0ca400944c1111")
	assert.Equal(t, 1, l)
}
