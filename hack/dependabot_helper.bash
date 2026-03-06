#!/usr/bin/env bash

set -euo pipefail

MERGEABILITY_MAX_RETRIES=60
MERGEABILITY_SLEEP=5
CHECKS_MAX_RETRIES=60
CHECKS_SLEEP=5
MERGE_MAX_RETRIES=30
MERGE_SLEEP=10
FAILURES=0

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"
}

# list open PRs from dependabot
prs=$(gh pr list --json number,author,title)

declare -A prMap

while read -r pr; do
    author=$(echo "$pr" | jq -r '.author.login')
    if [[ "$author" != "app/dependabot" ]]; then
        title=$(echo "$pr" | jq -r '.title')
        log "Skipping non-dependabot PR: $title (author: $author)"
        continue
    fi
    mod=$(echo "$pr" | jq -r '.title' | sed -e 's/.*[Bb]ump \([^ ]*\) .*/\1/')
    prMap[$mod]="$pr"
done < <(echo "${prs}" | jq -c '.[]')

merge() {
    local mod="$1"
    if [[ ! "${prMap[$mod]+x}" ]]; then
        return 0
    fi

    log "Merging $mod"
    local pr
    pr=$(echo "${prMap[$mod]}" | jq '.number')

    # Wait for PR to become mergeable
    local retries=0
    while ! gh pr view "$pr" --json 'mergeable' | grep -q 'MERGEABLE'; do
        retries=$((retries + 1))
        if [[ $retries -ge $MERGEABILITY_MAX_RETRIES ]]; then
            log "ERROR: Timed out waiting for PR #$pr to become mergeable"
            FAILURES=$((FAILURES + 1))
            return 1
        fi
        sleep "$MERGEABILITY_SLEEP"
    done

    gh pr review "$pr" --approve

    # Wait for checks to pass
    retries=0
    while true; do
        local checks_output
        checks_output=$(gh pr checks "$pr" 2>&1) || true

        if echo "$checks_output" | grep -qi 'fail'; then
            log "ERROR: Checks failed for PR #$pr ($mod)"
            FAILURES=$((FAILURES + 1))
            return 1
        fi

        if echo "$checks_output" | grep -qi 'pass'; then
            if ! echo "$checks_output" | grep -qi 'pending'; then
                log "All checks passed for PR #$pr"
                break
            fi
        fi

        retries=$((retries + 1))
        if [[ $retries -ge $CHECKS_MAX_RETRIES ]]; then
            log "ERROR: Timed out waiting for checks on PR #$pr"
            FAILURES=$((FAILURES + 1))
            return 1
        fi
        sleep "$CHECKS_SLEEP"
    done

    # Merge with retries
    retries=0
    while ! gh pr merge "$pr" --admin --merge; do
        retries=$((retries + 1))
        if [[ $retries -ge $MERGE_MAX_RETRIES ]]; then
            log "ERROR: Timed out trying to merge PR #$pr ($mod)"
            FAILURES=$((FAILURES + 1))
            return 1
        fi
        sleep "$MERGE_SLEEP"
    done

    log "Successfully merged PR #$pr ($mod)"
    unset "prMap[$mod]"
}

merge "k8s.io/apimachinery" || true
merge "k8s.io/api" || true
merge "k8s.io/client-go" || true

for k in "${!prMap[@]}"; do
    merge "$k" || true
done

if [[ $FAILURES -gt 0 ]]; then
    log "ERROR: $FAILURES PR(s) failed to merge"
    exit 1
fi

log "All dependabot PRs processed successfully"
