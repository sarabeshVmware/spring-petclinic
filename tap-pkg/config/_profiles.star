load("@ytt:data", "data")
load("@ytt:assert", "assert")
load("@ytt:struct", "struct")

_full_profile = "full"
_workspace_profile = "workspace"
_build_profile = "build"
_run_profile = "run"
_tbd_profile = "tbd"

_all_profiles = [_full_profile, _workspace_profile, _build_profile, _run_profile, _tbd_profile]

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
  ingress_values_dict = struct.decode(ingress_values)

  if 'ingressEnabled' in ingress_values_dict and 'ingressEnabled' not in pkg_values_dict:
    pkg_values_dict['ingressEnabled'] = ingress_values['ingressEnabled']
  end

  if 'ingressDomain' in ingress_values_dict and 'ingressDomain' not in pkg_values_dict:
    pkg_values_dict['ingressDomain'] = ingress_values.ingressDomain
  end

  if 'tls' in ingress_values_dict:
    if 'tls' not in pkg_values_dict:
      pkg_values_dict['tls'] = {}
    end

    if 'namespace' in ingress_values_dict['tls'] and 'namespace' not in pkg_values_dict['tls']:
      pkg_values_dict['tls']['namespace'] = ingress_values.tls.namespace
    end

    if 'secretName' in ingress_values_dict['tls'] and 'secretName' not in pkg_values_dict['tls']:
      pkg_values_dict['tls']['secretName'] = ingress_values.tls.secretName
    end
  end


  return struct.encode(pkg_values_dict)
end

profiles = struct.make(
	full=_full_profile,
	workspace=_workspace_profile,
	build=_build_profile,
	run=_run_profile,
  tbd=_tbd_profile,

	is_any_enabled=_is_any_enabled,
	is_enabled=_is_enabled,
	is_pkg_enabled=_is_pkg_enabled,
    merge_ingress_values=_merge_ingress_values,
)
