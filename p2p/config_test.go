package p2p

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile(t *testing.T) {
	config, err := LoadConfigFromFile("../config", "bootstrap")
	assert.Nil(t, err)
	assert.Equal(t, 65532, config.port)
	assert.Equal(t, "route_table", config.routeFilePath)

}
