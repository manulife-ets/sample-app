#!/bin/bash
set -o errexit
set -o xtrace

for filename in *.js; do

    k6 run --vus 1000 --duration 10s $filename

done

# k6 run --vus 1000 --duration 10s findOne_found.js

# k6 run --vus 1000 --duration 10s findOne_notfound.js

# k6 run --vus 1000 --duration 10s post-get-put-validate-script.js
