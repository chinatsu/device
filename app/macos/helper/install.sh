#!/usr/bin/env bash
self_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

user=$(/bin/ls -l /dev/console | /usr/bin/awk '{print $3}')
config_dir="/Users/${user}/Library/Application Support/naisdevice"

daemons=(io.naisdevice.{interface,sync})
for d in "${daemons[@]}"; do
  template="${self_dir}/daemon_templates/${d}.plist"
  destination="/Library/LaunchDaemons/${d}.plist"

  launchctl list | grep -q "$d" && launchctl unload "$destination"
  sed -e "s#CONFIG_DIR#${config_dir}#" "$template" > "$destination"
  launchctl load "$destination"
  echo "Installed service $d"
done
