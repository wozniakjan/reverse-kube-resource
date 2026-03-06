#!/usr/bin/env bash
set -euo pipefail
FAILURES=0
log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"; }
err() { log "ERROR: $*"; FAILURES=$((FAILURES + 1)); }

prs=$(gh pr list --limit 999 --json number,author,title) || { log "ERROR: Failed to fetch PRs"; exit 1; }
declare -A prMap
while read -r pr; do
    author=$(echo "$pr" | jq -r '.author.login')
    [[ "$author" =~ ^dependabot(\[bot\])?$ || "$author" == "app/dependabot" ]] || continue
    title=$(echo "$pr" | jq -r '.title')
    mod=$(echo "$title" | sed 's/.*[Bb]ump \([^ ]*\) .*/\1/')
    [[ "$mod" != *" "* && -n "$mod" ]] || continue
    dir=$(echo "$title" | sed -n 's/.* in \(\/[^ ]*\) from .*/\1/p')
    prMap["${dir:-/}:${mod}"]="$pr"
done < <(echo "$prs" | jq -c '.[]')

merge() {
    local key="$1"; [[ "${prMap["$key"]+x}" ]] || return 0
    local pr; pr=$(echo "${prMap["$key"]}" | jq -r '.number')
    [[ "$pr" =~ ^[0-9]+$ ]] || { err "invalid PR number '$pr' for $key"; return 1; }
    log "Processing PR #$pr ($key)"
    local i state=""
    for ((i = 0; i < 60; i++)); do # wait for mergeable (5min)
        state=$(gh pr view "$pr" --json mergeable --jq .mergeable) || state="API_ERROR"
        [[ "$state" == "MERGEABLE" ]] && break
        [[ "$state" == "UNKNOWN" || "$state" == "API_ERROR" ]] || break # terminal state, stop polling
        sleep 5
    done
    [[ "$state" == "MERGEABLE" ]] || { err "PR #$pr not mergeable ($state)"; return 1; }
    gh pr review "$pr" --approve || { err "approve failed for PR #$pr"; return 1; }
    local j fail=0 pend=1
    for ((i = 0; i < 60; i++)); do # wait for checks (5min)
        j=$(gh pr view "$pr" --json statusCheckRollup) || { err "checks query failed for PR #$pr"; return 1; }
        fail=$(echo "$j" | jq '[.statusCheckRollup[]? | select(.state? == "FAILURE" or .state? == "ERROR" or .conclusion? == "FAILURE" or .conclusion? == "CANCELLED" or .conclusion? == "TIMED_OUT" or .conclusion? == "ACTION_REQUIRED")] | length')
        pend=$(echo "$j" | jq '[.statusCheckRollup[]? | select(.state? == "EXPECTED" or .state? == "PENDING" or .state? == "IN_PROGRESS" or .state? == "QUEUED" or (.status? != null and .status? != "COMPLETED"))] | length')
        ((fail > 0)) && { err "checks failed for PR #$pr"; return 1; }
        ((pend == 0)) && break; sleep 5
    done
    ((pend == 0)) || { err "checks timed out for PR #$pr"; return 1; }
    local merged=false
    for ((i = 0; i < 30; i++)); do gh pr merge "$pr" --admin --merge && { merged=true; break; }; sleep 10; done
    $merged || { err "merge timed out for PR #$pr"; return 1; }
    log "Merged PR #$pr ($key)"; unset "prMap[$key]"
}

for mod in "k8s.io/apimachinery" "k8s.io/api" "k8s.io/client-go"; do
    for key in "${!prMap[@]}"; do [[ "$key" == *":$mod" ]] && { merge "$key" || true; }; done
done
for key in "${!prMap[@]}"; do merge "$key" || true; done
if ((FAILURES > 0)); then log "ERROR: $FAILURES PR(s) failed"; exit 1; fi
log "All dependabot PRs processed successfully"
