local utils = require("utils")

Config = {
  name = "Git worktrees",
  description = "Is the current directory using git worktrees?",
  type = "metadata"
}

function GetPluginMetadata(data)
  local metadata = {}

  for _, dir in ipairs(data) do
    local currentDir = utils.PathJoin(dir.path, "worktrees")
    if utils.Exists(currentDir) then
      local dirContents = utils.ReadDir(currentDir)
      local worktrees = {}

      for _, item in ipairs(dirContents) do
        if item.isDir then
          table.insert(worktrees, item.name)
        end
      end

      table.insert(metadata, "Worktrees: " .. table.concat(worktrees, ", "))
    else
      table.insert(metadata, "Worktrees: ❌")
    end
  end

  return metadata
end
