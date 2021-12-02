#!/usr/bin/bash

cat <<EOF | ytt -f _profiles_test.star -f ../_profiles.star -f-
#@data/values
---
profile: full
EOF
