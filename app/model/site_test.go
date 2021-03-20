package model_test

import (
	"testing"

	"github.com/josephspurrier/polarbearblog/app/model"
	"github.com/stretchr/testify/assert"
)

func TestSiteURL(t *testing.T) {
	s := new(model.Site)
	s.Scheme = "http"
	s.URL = "localhost"
	assert.Equal(t, "http://localhost", s.SiteURL())
}
