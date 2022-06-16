#!/usr/bin/env bash

local_dir="$(dirname "${BASH_SOURCE[0]}")"
cd "${local_dir}"

# prep
component_lib_file=out/function-to-test.lib.yaml
trap 'rm -rf "${component_lib_file}"' EXIT
mkdir -p out

cat <<EOF | ytt -f _profiles_test.star -f ../../../tap-pkg/config/_profiles.star -f-
#@data/values
---
profile: full
EOF

# components below were named according to their PackageInstall filename under ../../../tap-pkg/config
# i.e: ../../../tap-pkg/app-accelerator.yaml
COMPONENTS=(
  "app-accelerator"
  "app-live-view-connector"
  "app-live-view"
  "cloud-native-runtimes"
  "learning-center"
  "tap-gui"
  "metadata-store"
)

for c in "${COMPONENTS[@]}"; do
  component_to_test="$c"
  component_test_file="_$(echo "$component_to_test" | tr '-' '_')_test.star"
  echo -e "-----"
  echo -e "Start: Shared Values Test - $component_to_test"
  echo -e "-----"
  # can only load files a file with a `.lib.yaml` or `.star` extension.
  ln -sf "../../../../tap-pkg/config/${component_to_test}.yaml" $component_lib_file
  ytt --allow-symlink-destination ../../../tap-pkg -f "${component_test_file}" -f ../../../tap-pkg/config/_profiles.star -f $component_lib_file -f ../../../tap-pkg/config/values.yml
done
