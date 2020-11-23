package test

import (
	"github.com/stretchr/testify/require"
	"regexp"
	"strings"
	"testing"
)

func TestReplace(t *testing.T) {
	otherReg, _ := regexp.Compile("[^a-z-.]")
	trimReg, _ := regexp.Compile("(^[.-]+|[.-]+$)")
	s := "github.com-maxiloEmmmm-go-web?"

	s = strings.ToLower(s)
	s = otherReg.ReplaceAllString(s, ".")
	s = trimReg.ReplaceAllString(s, "")

	require.Equal(t, s, "github.com-maxiloemmmm-go-web")
}
