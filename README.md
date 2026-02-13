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

## Themes

Switcheroo can load color themes from JSON files. Set the optional `theme` key in
`config.json` to the name (or path) of the theme you want:

```json
{
  "directory": "/path/to/repos",
  "output": "path",
  "theme": "dracula"
}
```

If `theme` is omitted or set to `"default"`, Switcheroo uses your terminal's
colors. When a theme name is provided, Switcheroo searches for a matching
`<name>.json` file in:

1. `./themes/`
2. `<switcheroo binary>/themes/`
3. `$XDG_CONFIG_HOME/switcheroo/themes/`
4. `~/.config/switcheroo/themes/`
5. `~/.switcheroo/themes/`

You can also provide an absolute/relative path (with or without `.json`). A
theme file lists color values for the UI; any fields you leave empty fall back
to the terminal defaults.

Example (`themes/dracula.json`):

```json
{
  "name": "Dracula",
  "accent": "#bd93f9",
  "muted": "#6272a4",
  "disabled": "#44475a",
  "success": "#50fa7b",
  "warning": "#ffb86c",
  "error": "#ff5555",
  "background": "#282a36",
  "foreground": "#f8f8f2",
  "border": "#6272a4",
  "current_line": "#8be9fd",
  "input_bg": "#1e1f29",
  "list_bg": "#21222c",
  "preview_bg": "#1d1e27",
  "selected_bg": "#44475a"
}
```
