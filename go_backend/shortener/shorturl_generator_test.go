package shortener

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const timeStr = "1257894000000000000"

func TestShortLinkGenerator(t *testing.T) {
	initialLink_1 := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortLink_1 := generateShortLink(initialLink_1, timeStr)

	initialLink_2 := "https://www.eddywm.com/lets-build-a-url-shortener-in-go-with-redis-part-2-storage-layer/"
	shortLink_2 := generateShortLink(initialLink_2, timeStr)

	initialLink_3 := "https://spectrum.ieee.org/automaton/robotics/home-robots/hello-robots-stretch-mobile-manipulator"
	shortLink_3 := generateShortLink(initialLink_3, timeStr)


	assert.Equal(t, shortLink_1, "8UCsi2gD")
	assert.Equal(t, shortLink_2, "2fNGJSK6")
	assert.Equal(t, shortLink_3, "CGHVY6U7")
}