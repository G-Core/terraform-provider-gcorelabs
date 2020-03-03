package testing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractToken(t *testing.T) {
	result := getGetResult(t)

	token, err := result.ExtractTokens()
	require.NoError(t, err)
	require.Equal(t, &expectedToken, token)
}
