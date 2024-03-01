# Plugins

Plugins allow for extended behaviour with `lua`. Currently the only extensible part
of Switcheroo is the metadata creation.

## Using Plugins

Plugins need to be located in the `/plugins` directory in the config directory:
e.g. `~/.config/switcheroo/plugins`

There are a few example plugins in the `/plugins` directory in this repo which 
you can copy in to the correct directory as outlined above.

All plugins within the config plugins directory will be automatically loaded.

## Creating plugins

Each plugin requires a global `Config` object to outline a few basic parameters:

```lua
Config = {
  name = "Project Group",
  description = "Which group does a given project belong to?",
  type = "metadata",
}
```

Currently the only `type` that will work is `metadata`

### Metadata plugins

Metadata plugins need to export a `GetPluginMetadata` function.  This function will
 be called when Switcheroo runs.  

`GetPluginMetadata` takes in a single lua table property which is all of the directories 
 that will be passed to the fuzzy finder.  This table has the following structure:

```lua
data = {
    {
        path = "string"
        name = "string"
        group = "string"
    },
    ...
}

```

The `GetMetadataPlugin` function must return a table (array) of strings which will be 
concatenated with the output of the other metadata plugins later

**Example:**

```lua
function GetPluginMetadata(data)
  local metadata = {}

  for _, v in ipairs(data) do
    table.insert(metadata, "Context: " .. v.group)
  end

  return metadata
end
```
