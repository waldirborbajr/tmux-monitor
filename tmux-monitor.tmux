#!/usr/bin/env bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

get_tmux_option() {
  local option=$1
  local default_value=$2
  local option_value=$(tmux show-option -gqv "$option")
  if [ -z "$option_value" ]; then
    echo "$default_value"
  else
    echo "$option_value"
  fi
}

set_tmux_option() {
  local option=$1
  local value=$2
  tmux set-option -gq "$option" "$value"
}

# docker_monitor() {
#   "$CURRENT_DIR/tmux-monitor s"
# }

main() {
  local status_right=$(get_tmux_option "status-right" "")
  local new_status_right="#($CURRENT_DIR/tmux-monitor s) $status_right"
  set_tmux_option "status-right" "$new_status_right"

  # Set a shorter status-interval to update more frequently
  set_tmux_option "status-interval" "30"

  # Force refresh of tmux status bar
  tmux refresh-client -S

  # Start the monitor script in the background
  tmux run-shell -b "$CURRENT_DIR/tmux-monitor s"
}

main
