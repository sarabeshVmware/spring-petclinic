load("@ytt:data", "data")
load("@ytt:assert", "assert")

if data.values.ceip_policy_disclosed != True:
	assert.fail("The field ceip_policy_disclosed in values.yaml must be set to true in order to proceed with the installation.  For more information on VMware's Customer Experience Improvement Program, visit this link: https://www.vmware.com/solutions/trustvmware/ceip-products.html.")
end