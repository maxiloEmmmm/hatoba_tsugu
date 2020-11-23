package test

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"hatoba_tsugu/api/types"
	"testing"
)

func TestLaunch(t *testing.T) {
	response, err := resty.New().R().SetBody(types.Launch{
		Env:   "dev",
		Git:   "https://github.com/maxiloEmmmm/go-web",
		Image: "192.168.1.246:5000/nginx:latest",
	}).Post("http://localhost:8000/launch")
	require.Nil(t, err)
	t.Log(response)
	require.Equal(t, response.StatusCode(), 200)
}
