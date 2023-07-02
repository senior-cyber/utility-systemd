### systemd.json

```json
{
  "name": "my-service",
  "user": "root",
  "group": "root"
}
```

### install / uninstall

```shell
./cli install systemd.json config.json
./cli uninstall systemd.json
```