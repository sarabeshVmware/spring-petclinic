load("@ytt:data", "data")
load("@ytt:assert", "assert")
load("@ytt:struct", "struct")

_full_profile = "full"
_dev_light_profile = "dev-light"
_all_profiles = [_full_profile, _dev_light_profile]

if not data.values.profile in _all_profiles:
  assert.fail("Expected profile to be one of: {}".format(_all_profiles))
end

def _is_enabled(profile):
  return data.values.profile == profile
end

def _is_any_enabled(profiles):
  return data.values.profile in profiles
end

profiles = struct.make(
	full=_full_profile,
	dev_light=_dev_light_profile,

	is_any_enabled=_is_any_enabled,
	is_enabled=_is_enabled,
)
