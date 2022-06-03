#!/usr/bin/bash

cat <<EOF | ytt -f _profiles_test.star -f ../../../tap-pkg/config/_profiles.star -f-
#@data/values
---
profile: full
EOF
