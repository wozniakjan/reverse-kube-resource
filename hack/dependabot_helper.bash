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

# list open PRs (fetch all, not just default 30)
prs=$(gh pr list --limit 999 --json number,author,title)

declare -A prMap

while read -r pr; do
    author=$(echo "$pr" | jq -r '.author.login')
    # Accept common Dependabot author logins: dependabot, dependabot[bot], and app/dependabot
    if [[ ! "$author" =~ ^dependabot(\[bot\])?$ && "$author" != "app/dependabot" ]]; then
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
    pr=$(echo "${prMap[$mod]}" | jq -r '.number')

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

    if ! gh pr review "$pr" --approve; then
        log "ERROR: Failed to approve PR #$pr ($mod)"
        FAILURES=$((FAILURES + 1))
        return 1
    fi

    # Wait for checks to pass using structured JSON
    retries=0
    while true; do
        local checks_json
        if ! checks_json=$(gh pr view "$pr" --json statusCheckRollup 2>&1); then
            log "ERROR: Failed to fetch status checks for PR #$pr ($mod): $checks_json"
            FAILURES=$((FAILURES + 1))
            return 1
        fi

        local failed_count pending_count total_count
        failed_count=$(echo "$checks_json" | jq '[.statusCheckRollup[]? | select(.state == "FAILURE" or .state == "ERROR")] | length')
        pending_count=$(echo "$checks_json" | jq '[.statusCheckRollup[]? | select(.state == "PENDING" or .state == "IN_PROGRESS" or .state == "QUEUED")] | length')
        total_count=$(echo "$checks_json" | jq '.statusCheckRollup | length')

        if (( failed_count > 0 )); then
            log "ERROR: Checks failed for PR #$pr ($mod)"
            FAILURES=$((FAILURES + 1))
            return 1
        fi

        if (( total_count > 0 )) && (( pending_count == 0 )); then
            log "All checks passed for PR #$pr"
            break
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

# Merge pinned modules in order; unset from prMap on success to avoid double-processing
for mod in "k8s.io/apimachinery" "k8s.io/api" "k8s.io/client-go"; do
    if merge "$mod"; then
        unset "prMap[$mod]" 2>/dev/null || true
    fi
done

for k in "${!prMap[@]}"; do
    merge "$k" || true
done

if [[ $FAILURES -gt 0 ]]; then
    log "ERROR: $FAILURES PR(s) failed to merge"
    exit 1
fi

log "All dependabot PRs processed successfully"
