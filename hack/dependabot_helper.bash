#!/usr/bin/env bash

set -euo pipefail

DIR=$(dirname "${BASH_SOURCE[0]}")/..

# list open PRs from dependabot touching KEB go modules
prs=$(gh pr list --json number,author,title)

# iterate over each PR, run go mod tidy under the KCP CLI dir, commit, push
git checkout main
git pull origin main
declare -A prMap

while read pr; do
    mod=$(echo "$pr" | jq '.title' | sed -e 's/.*Bump \([^ ]*\) .*/\1/')
    prMap[$mod]="$pr"
done < <(echo "${prs}" | jq -c '.[]')

function merge() {
    mod=$1
    if [[ "${prMap[$mod]+x}" ]]; then
        echo merging $mod
        pr=$(echo ${prMap[$mod]} | jq '.number')
        while ! gh pr view $pr --json 'mergeable' | grep -e 'MERGEABLE'; do
            sleep 5
        done
        gh pr review "${pr}" --approve
        while ! gh pr merge "$pr" --auto --merge; do
            sleep 10
        done
        gh pr merge "${pr}" --admin --rebase
        unset prMap[$mod]
    fi
    
}

merge "k8s.io/apimachinery"
merge "k8s.io/api"
merge "k8s.io/client-go"

for k in "${!prMap[@]}"; do
    merge "$k"
done
