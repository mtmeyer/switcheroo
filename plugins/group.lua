Config = {
  name = "Project Group",
  description = "Which group does a given project belong to?",
  type = "metadata",
}

-- data = {
--  {
--    path: string
--    name: string
--    group: string
--  },
--  ...
-- }
function GetPluginMetadata(data)
  local metadata = {}

  for _, v in ipairs(data) do
    table.insert(metadata, "Context: " .. FirstToUpper(v.group))
  end

  return metadata
end

function FirstToUpper(str)
  return (str:gsub("^%l", string.upper))
end
