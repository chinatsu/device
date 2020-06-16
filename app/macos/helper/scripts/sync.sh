#!/bin/bash
interface="$1"
conf="$2"

wg=/opt/naisdevice/bin/naisdevice-wg
pgrep=/usr/bin/pgrep
route=/sbin/route

if $pgrep "device-agent"; then
  $wg syncconf "$interface" "$conf"
else
  echo "no device-agent, aborted syncconf"
  exit 1
fi

#declare -A desiered_routes
#for i in $(while read -r _ ips; do for ip in $ips; do [[ $ip =~ ^[0-9a-z:.]+/[0-9]+$ ]] && echo "$ip"; done; done < <($wg show "$interface" allowed-ips) | sort -nr -k 2 -t /); do
#  desiered_routes[$i]=true
#done
#
#for r in "${!desiered_routes[@]}"; do
#  r=${r%"/32"}
#  ri=$($route get "$r" | grep interface | cut -d ':' -f 2 | tr -d ' ')
#  if [[ "$ri" != "$interface" ]]; then
#    echo $route -q -n add -inet $r -interface "$interface"
#    $route -q -n add -inet $r -interface "$interface"
#  fi
#done
