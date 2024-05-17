# Switcheroo

Switcheroo is a terminal based project/directory switcher written in Go.

## Setup

Create a switcheroo config directory in `~/.config/switcheroo`

### config.json

Create a `config.json` file in the Switcheroo config directory with the following data:

```json
{
  "directories": {
    "personal": "/Users/SomeUser/git/personal",
    "work": "/Users/SomeUser/git/work"
  },
  "output": "path"
}
```

### Plugins

For any plugins you want to use, copy them from the `/plugins` directory in the repo
to the `plugins` directory in the Switcheroo config folder.

The two plugins currently in this repo are:

- group.lua - Display which group a particular project is part of (e.g. is this a work
or personal project)

- gitWorktrees.lua - Display if a given project is using git worktrees, and if yes, 
what worktrees are currently cloned

## Flags

`--directory`: Override what directory will be passed into the fuzzy finder

`--configFile`: Path to a config file not in the default directory

`--skipPlugins`: Run Switcheroo without running any plugins

`--output`: What type of output to return. Either 'path', 'index' or 'name' 
