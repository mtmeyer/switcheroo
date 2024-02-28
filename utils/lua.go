package utils

import (
	"github.com/yuin/gopher-lua"
)

func MetadataInputToLuaTable(L *lua.LState, metadata []MetadataPluginData) *lua.LTable {
	metadataTable := L.NewTable()

	for i, item := range metadata {
		itemTable := L.NewTable()

		L.SetField(itemTable, "path", lua.LString(item.Path))
		L.SetField(itemTable, "name", lua.LString(item.Name))
		L.SetField(itemTable, "parentDirectory", lua.LString(item.ParentDirectory))

		L.SetTable(metadataTable, lua.LNumber(i+1), itemTable) // Lua arrays are 1-indexed
	}

	return metadataTable
}

func LuaTableToStringSlice(table *lua.LTable) []string {
	var slice []string

	table.ForEach(func(key, value lua.LValue) {
		slice = append(slice, value.String())
	})

	return slice
}
