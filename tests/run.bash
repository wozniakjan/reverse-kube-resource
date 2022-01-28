#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname $0)/.."

declare -A test_args=()
test_args["basic"]="" 
test_args["unstructured"]="-unstructured" 
test_args["only_meta"]="-only-meta" 
test_args["boilerplate"]="-go-header-file=./tests/boilerplate.txt"
#test_args["kubermatic"]="-kubermatic-interfaces"

for f in examples/*; do
    t=$(basename "$f")
    for ta in "${!test_args[@]}"; do
        a=${test_args[$ta]}
        echo "go run ./main.go -package=examples -src="$f" $a > ./tests/output/${t}_${ta}.go"
        go run ./main.go -package=examples -src="$f" $a > ./tests/output/${t}_${ta}.go
        (
            cd tests
            go build ./output/${t}_${ta}.go
        )
    done
done
