#!/usr/bin/env bash
#
# emergency-halt.sh — convenience wrapper around the x/validator emergency
# levers (hard halt / soft freeze) so an operator can stop and resume the chain
# with a single command instead of remembering the full tx flag set.
#
# Underlying node CLI (already built into dscd):
#   dscd tx validator halt-chain    [reason] [--height H]   # hard consensus halt
#   dscd tx validator resume-chain                          # clear a hard halt
#   dscd tx validator freeze-chain  [reason]                # soft freeze (reject non-emergency txs)
#   dscd tx validator unfreeze-chain                        # lift a soft freeze
#   dscd query validator emergency-status                   # inspect current state
#
# Only the dedicated halt-admin key (types.GetHaltAdmin) may sign these.
#
# Usage:
#   scripts/emergency-halt.sh status
#   scripts/emergency-halt.sh halt      ["reason"] [--height H]
#   scripts/emergency-halt.sh resume
#   scripts/emergency-halt.sh freeze    ["reason"]
#   scripts/emergency-halt.sh unfreeze
#
# Configuration (override via env):
#   DSCD             node binary            (default: dscd)
#   NODE             tendermint RPC         (default: https://testnet-val.decimalchain.com:443/rpc/)
#                    NOTE: dscd requires an explicit host:PORT — use :443 for the HTTPS gateway,
#                    and keep the trailing slash on the /rpc/ path.
#   CHAIN_ID         chain id               (default: decimal_202020-221213  — testnet)
#   FROM             halt-admin key name    (default: haltadmin)
#   KEYRING_BACKEND  keyring backend        (default: test)
#   FEES             tx fees                (default: empty — let the chain auto-charge)
#                    Emergency msgs are fee-free at the msg level (only a tiny per-byte
#                    commission applies). Leaving FEES empty makes dscd submit a zero-fee
#                    tx, and the node's FeeDecorator deducts the exact commission from the
#                    signer's base-coin balance. Do NOT hardcode a small fee: the required
#                    commission = TxByteFee * bytes / delPrice, and with delPrice ~0.0028
#                    USD/tdel on testnet that is several tdel — a low fixed fee is rejected
#                    with "insufficient funds to pay for fees" (FeeLessThanCommission).
#                    Set FEES=<n>tdel only to override with an explicit (sufficient) fee.
#   GAS              gas                    (default: 400000)
#   YES              set to 1 to skip the interactive confirmation
#
set -euo pipefail

DSCD="${DSCD:-dscd}"
NODE="${NODE:-https://testnet-val.decimalchain.com:443/rpc/}"
CHAIN_ID="${CHAIN_ID:-decimal_202020-221213}"
FROM="${FROM:-haltadmin}"
KEYRING_BACKEND="${KEYRING_BACKEND:-test}"
FEES="${FEES:-}"
GAS="${GAS:-400000}"
YES="${YES:-0}"

red()   { printf '\033[31m%s\033[0m\n' "$*"; }
green() { printf '\033[32m%s\033[0m\n' "$*"; }
bold()  { printf '\033[1m%s\033[0m\n' "$*"; }

tx_flags() {
  local fee_flag=""
  # Only pass --fees when explicitly overridden; empty fees trigger the node's
  # "free" emergency path, which auto-charges the exact per-byte commission.
  [ -n "$FEES" ] && fee_flag="--fees $FEES"
  echo --from "$FROM" \
       --chain-id "$CHAIN_ID" \
       --node "$NODE" \
       --keyring-backend "$KEYRING_BACKEND" \
       $fee_flag \
       --gas "$GAS" \
       --broadcast-mode block \
       -y
}

query_flags() {
  echo --node "$NODE" --chain-id "$CHAIN_ID" -o json
}

confirm() {
  [ "$YES" = "1" ] && return 0
  local prompt="$1"
  red   "  chain-id : $CHAIN_ID"
  red   "  node     : $NODE"
  red   "  from     : $FROM"
  printf '%s ' "$prompt [type YES to proceed]:"
  read -r ans
  [ "$ans" = "YES" ] || { echo "aborted."; exit 1; }
}

status() {
  bold "Emergency status ($CHAIN_ID):"
  # shellcheck disable=SC2046
  "$DSCD" query validator emergency-status $(query_flags)
}

do_halt() {
  local reason="" height_flag=()
  while [ $# -gt 0 ]; do
    case "$1" in
      --height) height_flag=(--height "$2"); shift 2 ;;
      *)        reason="$1"; shift ;;
    esac
  done
  bold "HARD HALT — the chain will STOP producing blocks."
  echo "Recovery requires every validator to restart with --unsafe-skip-halt,"
  echo "then the halt-admin to broadcast 'resume'."
  confirm "Schedule hard halt${reason:+ (\"$reason\")}?"
  # shellcheck disable=SC2046
  "$DSCD" tx validator halt-chain "$reason" "${height_flag[@]}" $(tx_flags)
  green "halt-chain broadcast. Check: $0 status"
}

do_resume() {
  bold "RESUME — clear a scheduled/active hard halt."
  echo "If the chain already halted, validators must first restart with --unsafe-skip-halt"
  echo "so this tx can be included."
  confirm "Clear hard halt?"
  # shellcheck disable=SC2046
  "$DSCD" tx validator resume-chain $(tx_flags)
  green "resume-chain broadcast. Check: $0 status"
}

do_freeze() {
  local reason="${1:-}"
  bold "SOFT FREEZE — blocks keep being produced, but all non-emergency txs are rejected."
  echo "Reversible with 'unfreeze' — no node restart needed."
  confirm "Soft-freeze the chain${reason:+ (\"$reason\")}?"
  # shellcheck disable=SC2046
  "$DSCD" tx validator freeze-chain "$reason" $(tx_flags)
  green "freeze-chain broadcast. Check: $0 status"
}

do_unfreeze() {
  bold "UNFREEZE — the chain accepts transactions again."
  confirm "Lift the soft freeze?"
  # shellcheck disable=SC2046
  "$DSCD" tx validator unfreeze-chain $(tx_flags)
  green "unfreeze-chain broadcast. Check: $0 status"
}

usage() {
  sed -n '2,40p' "$0"
  exit 1
}

cmd="${1:-}"; shift || true
case "$cmd" in
  status)   status ;;
  halt)     do_halt "$@" ;;
  resume)   do_resume ;;
  freeze)   do_freeze "$@" ;;
  unfreeze) do_unfreeze ;;
  *)        usage ;;
esac
