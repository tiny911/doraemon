package config

import (
	"os"
	"testing"
)

var (
	cfg = struct {
		TestKv    string            `yaml:"TestKv"`
		TestSlice []string          `yaml:"TestSlice"`
		TestMap   map[string]string `yaml:"TestMap"`
	}{}

	//下面的两个值从config.yaml中拷过来的
	prj  = "test_prj"
	test = "TestKvVal"
)

func TestConfigure_Parse(t *testing.T) {
	Parse(&cfg)

	v := os.Getenv("ENV_PRJ_NAME")
	if v != prj {
		t.Errorf("parse env failed! please check config.yaml.")
	}

	if cfg.TestKv != test {
		t.Errorf("parse cfg failed! please check config.yaml.")
	}
}
