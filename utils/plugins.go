package utils

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
)

var L *lua.LState

func InitLua() error {
	L = lua.NewState()
	defer L.Close()

	metadataData := []MetadataPluginData{
		{
			Name:            "some-name",
			Path:            "/some/path",
			ParentDirectory: "work",
		},
		{
			Name:            "some-other-name",
			Path:            "/some/other/path",
			ParentDirectory: "personal",
		},
	}

	metadata := GetMetadataFromPlugin("test", metadataData)

	fmt.Println(metadata)

	return nil
}

type MetadataPluginData struct {
	Path            string
	Name            string
	ParentDirectory string
}

func GetMetadataFromPlugin(pluginName string, data []MetadataPluginData) []string {
	L = lua.NewState()
	defer L.Close()

	metadata := CallMetadataPlugin("test", MetadataInputToLuaTable(L, data))

	var metadataStruct []string

	if tbl, ok := metadata.(*lua.LTable); ok {
		metadataStruct = LuaTableToStringSlice(tbl)
	} else {
		panic(errors.New("Return data type not valid"))
	}

	return metadataStruct
}

func CallMetadataPlugin(pluginName string, data *lua.LTable) lua.LValue {
	LoadPlugin(pluginName, "metadata")

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("GetPluginMetadata"),
		NRet:    1,
		Protect: true,
	}, data); err != nil {
		panic(err)
	}
	metadata := L.Get(-1) // returned value
	L.Pop(1)              // remove received value

	return metadata
}

type PluginConfig struct {
	Name string
	Type string
}

func LoadPlugin(pluginName string, pluginType string) {
	// Find and load plugin if exists
	pluginPath := path.Join(ConfigDirectory, "plugins", pluginName+".lua")
	fmt.Println(pluginPath)
	if _, err := os.Stat(pluginPath); err != nil {
		panic(errors.New("Plugin file doesn't exist"))
	}

	luaErr := L.DoFile(pluginPath)

	if luaErr != nil {
		panic(luaErr)
	}

	// Check if plugin is correct type
	var pluginConfig *PluginConfig
	err := gluamapper.Map(L.GetGlobal("Config").(*lua.LTable), &pluginConfig)

	if err != nil {
		panic(err)
	}

	if pluginConfig.Type != pluginType {
		return
	}
}
