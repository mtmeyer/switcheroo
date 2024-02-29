package utils

import (
	"errors"
	"os"
	"path"

	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

var L *lua.LState

func GetMetadataFromPlugin(pluginFileName string, data []Directory) []string {
	L = lua.NewState()
	defer L.Close()

	L.SetGlobal("require", luar.New(L, LuaImport))

	metadata := CallMetadataPlugin(pluginFileName, MetadataInputToLuaTable(L, data))

	if metadata == nil {
		return nil
	}

	var metadataStruct []string

	if tbl, ok := metadata.(*lua.LTable); ok {
		metadataStruct = LuaTableToStringSlice(tbl)
	} else {
		panic(errors.New("Return data type not valid"))
	}

	return metadataStruct
}

func CallMetadataPlugin(pluginFileName string, data *lua.LTable) lua.LValue {
	metadataPlugin := LoadPlugin(pluginFileName, "metadata")

	if !metadataPlugin {
		return nil
	}

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
	Name        string
	Type        string
	Description string
}

func LoadPlugin(pluginFileName string, pluginType string) bool {
	// Find and load plugin if exists
	pluginPath := path.Join(ConfigDirectory, "plugins", pluginFileName)

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
		return false
	}

	return true
}
