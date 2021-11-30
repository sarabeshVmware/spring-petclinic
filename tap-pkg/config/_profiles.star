load("@ytt:data", "data")
load("@ytt:assert", "assert")
load("@ytt:struct", "struct")

_full_profile = "full"
_dev_profile = "dev"
_build_profile = "build"
_run_profile = "run"

_all_profiles = [_full_profile, _dev_profile, _build_profile, _run_profile]

if not data.values.profile in _all_profiles:
  assert.fail("Expected profile to be one of: {}".format(_all_profiles))
end

def _is_enabled(profile):
  return data.values.profile == profile
end

def _is_any_enabled(profiles):
  return data.values.profile in profiles
end

def _is_pkg_enabled(name):
  return (name not in data.values.excluded_packages) 
end

def _merge_ingress_values(pkg_values, ingress_values):
  pkg_values_dict = struct.decode(pkg_values)
  pkg_values_dict['ingressEnabled'] = ingress_values.ingressEnabled
  pkg_values_dict['ingressDomain'] = ingress_values.ingressDomain
  pkg_values_dict['tls'] = {}
  pkg_values_dict['tls']['namespace'] = ingress_values.tls.namespace
  pkg_values_dict['tls']['secretName'] = ingress_values.tls.secretName
  return struct.encode(pkg_values_dict)
end

profiles = struct.make(
	full=_full_profile,
	dev=_dev_profile,
	build=_build_profile,
	run=_run_profile,

	is_any_enabled=_is_any_enabled,
	is_enabled=_is_enabled,
	is_pkg_enabled=_is_pkg_enabled,
    merge_ingress_values=_merge_ingress_values,
)
