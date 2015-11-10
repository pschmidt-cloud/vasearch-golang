package test

import (
	"vasearch/model"
	"fmt"
	"github.com/BurntSushi/toml"
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"runtime"
)

/*
func TestMain(m *testing.M) {

}
*/

func TestJson(t *testing.T) {
	m := model.Sample{"Liver Test", 1024, "hg18"}

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	results := "{\"results\":" + string(b) + "}"
	fmt.Println("results= ", results)
	//fmt.Println(string(b))

}

func TestStub(t *testing.T) {
	fmt.Println("testStub")
	assert.Equal(t, 100, 100, "101 != 100")
	assert.True(t, true, "This is good. Canary test passing")
}

func TestMisc(t *testing.T) {
	data := "{\"description\" : \"This a sample json response\"}"
	fmt.Println(data)
	fmt.Printf("env vars: %s %s\n\n", runtime.GOOS, runtime.GOARCH)
}

func TestConfigFile(t *testing.T) {
	var app model.AppLoader
	app.Test = 5

	var conf model.Config
	if _, err := toml.DecodeFile("../application.properties", &conf); err != nil {
		panic(err)
	}
	app.Config = conf

	fmt.Printf("Test=%d\n", app.Test)
	fmt.Printf("ctx=%v\n", app.Config)
	fmt.Printf("ctx.host=%s\n", app.Config.Host)

	assert.Equal(t, app.Config.Host, "localhost")
	assert.Equal(t, app.Config.Port, 8080)
	assert.Equal(t, app.Config.Cluster, "elasticsearch-lumacpschmidt")
}
