# psb
pretty-safe-backup ;)

psb must be ran as a superuser, an example Systemd service file can be found in [./examples](examples).

### Build
```sh
make build
```

### Install
```sh
make install
```

### Enable and start service
```sh
sudo cp ./examples/psb.service /etc/systemd/system/psb.service
systemctl enable psb.service
systemctl start psb.service
```