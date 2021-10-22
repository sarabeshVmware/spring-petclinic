load("@ytt:data", "data")
load("@ytt:assert", "assert")

if not "/" in data.values.rw_app_registry.server_repo:
  assert.fail("Expected data value 'rw_app_registry.server_repo' to include one '/'")
end
