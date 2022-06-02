#!/usr/bin/env bash

local_dir="$(dirname "${BASH_SOURCE[0]}")"
cd "${local_dir}"

# prep
component_to_test=out/function-to-test.lib.yaml
trap 'rm -rf "${component_to_test}"' EXIT
mkdir -p out

cat <<EOF | ytt -f _profiles_test.star -f ../../../tap-pkg/config/_profiles.star -f-
#@data/values
---
profile: full
EOF

echo -e "-----"
echo -e "Start: Shared Values Test - App Accelerator"
echo -e "-----"
# can only load files a file with a `.lib.yaml` or `.star` extension.
ln -srf ../../../tap-pkg/config/app-accelerator.yaml $component_to_test
ytt --allow-symlink-destination ../../../tap-pkg -f _accelerator_test.star -f ../../../tap-pkg/config/_profiles.star -f $component_to_test -f ../../../tap-pkg/config/values.yml

echo -e "-----"
echo -e "Start: Shared Values Test - App Live View - Connector"
echo -e "-----"
# can only load files a file with a `.lib.yaml` or `.star` extension.
ln -srf ../../../tap-pkg/config/app-live-view-connector.yaml $component_to_test
ytt --allow-symlink-destination ../../../tap-pkg -f _appliveview_connector_test.star -f ../../../tap-pkg/config/_profiles.star -f $component_to_test -f ../../../tap-pkg/config/values.yml

echo -e "-----"
echo -e "Start: Shared Values Test - App Live View"
echo -e "-----"
# can only load files a file with a `.lib.yaml` or `.star` extension.
ln -srf ../../../tap-pkg/config/app-live-view.yaml $component_to_test
ytt --allow-symlink-destination ../../../tap-pkg -f _appliveview_test.star -f ../../../tap-pkg/config/_profiles.star -f $component_to_test -f ../../../tap-pkg/config/values.yml

echo -e "-----"
echo -e "Start: Shared Values Test - Cloud Native Runtimes"
echo -e "-----"
# can only load files a file with a `.lib.yaml` or `.star` extension.
ln -srf ../../../tap-pkg/config/cloud-native-runtimes.yaml $component_to_test
ytt --allow-symlink-destination ../../../tap-pkg -f _cnrs_test.star -f ../../../tap-pkg/config/_profiles.star -f $component_to_test -f ../../../tap-pkg/config/values.yml

echo -e "-----"
echo -e "Start: Shared Values Test - Learning Center"
echo -e "-----"
# can only load files a file with a `.lib.yaml` or `.star` extension.
ln -srf ../../../tap-pkg/config/learning-center.yaml $component_to_test
ytt --allow-symlink-destination ../../../tap-pkg -f _learningcenter_test.star -f ../../../tap-pkg/config/_profiles.star -f $component_to_test -f ../../../tap-pkg/config/values.yml

echo -e "-----"
echo -e "Start: Shared Values Test - TAP GUI"
echo -e "-----"
# can only load files a file with a `.lib.yaml` or `.star` extension.
ln -srf ../../../tap-pkg/config/tap-gui.yaml $component_to_test
ytt --allow-symlink-destination ../../../tap-pkg -f _tap_gui_test.star -f ../../../tap-pkg/config/_profiles.star -f $component_to_test -f ../../../tap-pkg/config/values.yml
