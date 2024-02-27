package utils

import (
	"fmt"

	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

var L *lua.LState

func InitLua() error {
	L = lua.NewState()
	defer L.Close()

	// L.SetGlobal("require", luar.New(L, luaImport))
	module := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"setPluginType": func(L *lua.LState) int {
			return setPluginType(L)
		},
	})

	L.SetGlobal("sw", module)

	err := L.DoFile("test.lua")

	if err != nil {
		return err
	}

	return nil
}

func AddItemMetadata(L *lua.LState) int {
	data := L.ToString(1)
	fmt.Println(data)
	return 0
}

type PluginConfig struct {
	Name string
	Type string
}

func setPluginType(L *lua.LState) int {
	luaTable := L.ToTable(1)

	fmt.Println(luaTable)

	var pluginConfig *PluginConfig

	err := gluamapper.Map(luaTable, &pluginConfig)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Printf("blah\n %+v\n", pluginConfig)

	return 0
}
