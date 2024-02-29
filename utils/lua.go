package utils

import (
	"fmt"
	"os"
	"path"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func MetadataInputToLuaTable(L *lua.LState, metadata []Directory) *lua.LTable {
	metadataTable := L.NewTable()

	for i, item := range metadata {
		itemTable := L.NewTable()

		L.SetField(itemTable, "path", lua.LString(item.Path))
		L.SetField(itemTable, "name", lua.LString(item.Name))
		L.SetField(itemTable, "group", lua.LString(item.Group))

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

// Import allows a lua plugin to import a package
func LuaImport(pkg string) *lua.LTable {
	switch pkg {
	case "fmt":
		return importFmt()
	case "utils":
		return importLuaUtils()
	default:
		return nil
	}
}

func importFmt() *lua.LTable {
	pkg := L.NewTable()

	L.SetField(pkg, "Println", luar.New(L, fmt.Println))
	L.SetField(pkg, "Printf", luar.New(L, fmt.Printf))
	L.SetField(pkg, "Sprintln", luar.New(L, fmt.Sprintln))
	L.SetField(pkg, "Sprintf", luar.New(L, fmt.Sprintf))
	L.SetField(pkg, "Sprint", luar.New(L, fmt.Sprint))

	return pkg
}

func importLuaUtils() *lua.LTable {
	pkg := L.NewTable()

	L.SetField(pkg, "ReadDir", L.NewFunction(LuaReadDir))
	L.SetField(pkg, "Exists", L.NewFunction(LuaExists))
	L.SetField(pkg, "PathJoin", L.NewFunction(LuaPathJoin))

	return pkg
}

type DirectoryInfo struct {
	Name  string
	IsDir bool
}

func LuaReadDir(L *lua.LState) int {
	path := L.ToString(1)

	dirContents, err := os.ReadDir(path)

	if err != nil {
		panic(err)
	}

	directories := []DirectoryInfo{}

	for _, dir := range dirContents {
		directories = append(directories, DirectoryInfo{
			Name:  dir.Name(),
			IsDir: dir.IsDir(),
		})
	}

	dirTable := L.NewTable()

	for i, item := range directories {
		itemTable := L.NewTable()

		L.SetField(itemTable, "name", lua.LString(item.Name))
		L.SetField(itemTable, "isDir", lua.LBool(item.IsDir))

		L.SetTable(dirTable, lua.LNumber(i+1), itemTable) // Lua arrays are 1-indexed
	}

	L.Push(dirTable)

	return 1
}

func LuaExists(L *lua.LState) int {
	path := L.ToString(1)

	_, err := os.Stat(path)

	if err == nil {
		L.Push(lua.LBool(true))
		return 1
	}

	if os.IsNotExist(err) {
		L.Push(lua.LBool(false))
		return 1
	}

	L.Push(lua.LBool(false))
	return 1
}

func LuaPathJoin(L *lua.LState) int {
	var args []string
	currentIndex := 1

	for len(L.ToString(currentIndex)) > 0 {
		args = append(args, L.ToString(currentIndex))
		currentIndex += 1
	}

	L.Push(lua.LString(path.Join(args...)))

	return 1
}
