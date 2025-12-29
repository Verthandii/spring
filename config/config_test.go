package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Verthandii/spring/config/source"
	"github.com/Verthandii/spring/encoding"
	_ "github.com/Verthandii/spring/encoding/form"
	_ "github.com/Verthandii/spring/encoding/json"
	_ "github.com/Verthandii/spring/encoding/xml"
	_ "github.com/Verthandii/spring/encoding/yaml"
)

var testCfg = testConfig{
	App:     "test",
	Version: 1,
}

type testConfig struct {
	App     string `json:"app" yaml:"app"`
	Version int    `json:"version" yaml:"version"`
}

func TestInitYaml(t *testing.T) {
	SetDefaultCodec(encoding.GetCodec("yaml"))
	cfg := testConfig{}
	err := Init(source.NewFile("testdata/config.yaml"), &cfg)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testCfg, cfg)
}

func TestInitJson(t *testing.T) {
	cfg := testConfig{}
	err := Init(source.NewFile("testdata/config.json"), &cfg)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testCfg, cfg)
}
