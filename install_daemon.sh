#!/usr/bin/env sh

stashd=$(which stashd)
rc=$?
if [ $rc != 0 ]; then
    echo "GOBIN is not part of your PATH, it is recommended that you set it"
    stashd=$(ls $GOPATH/bin/stashd)
    rc=$?
    if [ $rc != 0 ]; then
        echo "stashd not found, is your GOPATH set?"
        exit 1
    fi
fi

cat > "$HOME/.stash/rundaemon.sh" <<- EOM
#!/usr/bin/env sh
$stashd >> $HOME/.stash/stashd.log 2>&1
EOM
chmod +x $HOME/.stash/rundaemon.sh

# UPSTART daemon file:

if [ -d "$HOME/.config/upstart" ]; then
    echo "Making upstart daemon file at ~/.config/upstart/stashd.conf"

    cat > "$HOME/.config/upstart/stashd.conf" <<- EOM
pre-start script
    mkdir -p $HOME/.stash
end script

respawn
respawn limit 15 5

start on startup

exec $HOME/.stash/rundaemon.sh
EOM
    start stashd


# SYSTEMD daemon file:

elif [ -d "$HOME/.config/systemd" ]; then
    echo "Making systemd daemon file at ~/.config/systemd/system/stashd.service"

    cat > "$HOME/.config/systemd/system/stashd.service" <<- EOM
[Unit]
Description=Stash Backup Daemon
[Service]
PIDFile=/var/run/stashd.pid
ExecStartPre=/bin/rm -f /var/run/stashd.pid
ExecStart=$HOME/.stash/rundaemon.sh
Restart=on-abort
[Install]
WantedBy=multi-user.target
EOM
    systemctl enable stashd


# LAUNCHD daemon file:

elif [ -d "$HOME/Library/LaunchAgents" ]; then
    echo "Making launchd daemon file at ~/Library/LaunchAgents/stashd.plist"

    cat > "$HOME/Library/LaunchAgents/stashd.plist" <<- EOM
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>KeepAlive</key>
    <true/>
    <key>Label</key>
    <string>stashd</string>
    <key>ProgramArguments</key>
    <array>
        <string>$HOME/.stash/rundaemon.sh</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>WorkingDirectory</key>
    <string>$HOME/.stash</string>
    <key>StandardErrorPath</key>
    <string>$HOME/.stash/stashd.log</string>
    <key>StandardOutPath</key>
    <string>$HOME/.stash/stashd.log</string>
</dict>
</plist>
EOM
    launchctl load "$HOME/Library/LaunchAgents/stashd.plist"

else
    echo "Your OS does not support launchd, systemd, or upstart."
    echo "You can manually start the daemon via 'stashd'"
fi
