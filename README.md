# psb
pretty-safe-backup ;)

Backups are extremely important, and there's tons of great software available that come with different ways to customize your backups. One that inspired this project is Rsnapshot, it's perfect as long as the computer is always running, otherwise things get a little weird... I'm hoping this will fill the gap.  

## Setup
> Dependencies: openssh, rsync\
> Build dependencies: golang, make, goxc

#### Build
```sh
make build
```

#### Install
```sh
make install
```

#### Remote server
Skip this section if your backup directory is mounted locally.

> Dependencies: openssh, rsync, [psb-rotatorc](//github.com/orange-lightsaber/psb-rotatorc), [psb-rotatord](//github.com/orange-lightsaber/psb-rotatord)\
> SSH keys must be generated and authorized for passwordless login.

Ensure dependencies are met, and psb-rotatord is properly configured.

Add the following line to sudoers file, this is necessary to allow Rsync to maintain file ownership during transfers. Replace "psbuser" with the username of the SSH user.
```
psbuser ALL= NOPASSWD:/usr/bin/rsync
```

Then back to the host...

#### Enable and start service
An example Systemd service file can be found in [./examples](examples)
```sh
sudo cp ./examples/psb.service /etc/systemd/system/psb.service
systemctl enable psb.service
systemctl start psb.service
```

#### Configuration
Configuration files are writen and edited internally, if created manually, any edits in the future may require that you manually delete the config's rsync script in */etc/xdg/psb/*, which will be necessary until the configuration editor is up, and any edits to existing run configs (regardless of method) will require a restart of the rotator daemon for the changes to take effect and to drop the old data from memory.

Run configs, by default, go in */etc/xdg/psb/run/*.

- enabled: True to enable, or false to disable.

- compatibility-key: A string, unique to each config and is regenerated each time the application writes/edits the run config. If an existing run config is edited, it is required to restart the rotator daemon on the remote destination if you are using one.
- name: The name field in the config must match the name of the file(not including extension ".toml"), and should contain no spaces.
- description: A short description of the backup operation.
- source: Absolute path to source directory.
- includes: Add paths relative to source to directories or files to make exceptions to an excluded directory, excepts wildcard as well.
- excludes: Add paths relative to source to directories or files, excepts wildcard as well.
- backup-directory: Optional. Overrides the path to the backups directory for current operation.
- remote-host: Only for remote backup destination. Address to remote backup server.
- username: Only for remote backup destination. Username of SSH user on remote backup server.
- port: Only for remote backup destination. SSH port to remote backup server.
- private-key: Only for remote backup destination. Path to the private key used to authenticate SSH communication with remote backup server.
- frequency: Number of minutes to wait between snapshots.
- delay: Number of minutes to wait between adding the most recent snapshot to initial rotation. Works best in increments of whatever frequency is set to. If delay is less than frequency, every snapshot gets archived.
- initial: Number of days to keep timed snapshots.
- daily: Number of months to keep daily snapshots.
- monthly: Number of months to keep monthly snapshots.
- yearly: Number of years to keep yearly snapshots.

Example run config, *home.toml*:
```
enabled = true
compatibility-key = "d2520db8729506f8850d71d87d86a1f8"
name = "home"
description = "Emergency Hoth evacuation backups"
source = "/home/leia"
includes = ["excluded_directory/overriden_directory_to_keep"]
excludes = [".cache/*", "excluded_directory"]
backup-directory = ""
remote-host = "138.831.42.24"
username = "remoteuser"
port = "22"
private-key = "/home/leia/.ssh/id_rsa"
frequency = 5
delay = 60
initial = 7
daily = 2
monthly = 10
yearly = 2
```