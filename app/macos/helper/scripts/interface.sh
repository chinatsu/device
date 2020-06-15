#!/usr/bin/env bash
interface="$1"
bootstrap_config="$2"
pid_path=$3
mtu=1360

wireguard_go=/opt/naisdevice/bin/naisdevice-wireguard-go
jq=/opt/naisdevice/bin/jq
pkill=/usr/bin/pkill
ifconfig=/sbin/ifconfig
route=/sbin/route
ps=/bin/ps

if [[ -d "$pid_path" ]]; then
  echo "$pid_path not found, no agent running"
fi

agent_pid=$(cat $pid_path)

if [[ ! -f "$bootstrap_config" ]]; then
  echo "no bootstrap config found"
  exit 1
fi

ip=$($jq -r .deviceIP < "$bootstrap_config")
echo $ip

if $ps -P $agent_pid; then
  if ! $pgrep -f "$wireguard_go" ; then
    echo "running wireguard-go"
    $wireguard_go "$interface"
    $ifconfig "$interface" inet "${ip}/21" "$ip" add
    $ifconfig "$interface" mtu "$mtu"
    $ifconfig "$interface" up
    $route -q -n add -inet "${ip}/21" -interface "$interface"
  else
    echo "wireguard-go already running"
  fi
else
  if $pgrep -f "$wireguard_go $interface"; then
    echo "stopping wireguard-go"
    $pkill -f "$wireguard_go $interface"
  else
    echo "wireguard-go already stopped"
  fi
fi
