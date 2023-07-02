### go.mod

```shell
go get github.com/senior-cyber/utility-systemd
```

### main.go

```go
package main

import (
	"github.com/senior-cyber/utility-systemd/systemd"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) >= 4 && strings.ToUpper(os.Args[1]) == "INSTALL" {
		_cli, _ := filepath.Abs(os.Args[0])
		systemdConfigFile, _ := filepath.Abs(os.Args[2])
		appConfigFile, _ := filepath.Abs(os.Args[3])
		installError := systemd.Install(systemdConfigFile, appConfigFile)
		if installError != nil {
			panic(installError)
		} else {
			log.Println(":/>sudo systemctl status " + filepath.Base(_cli))
			return
		}
	}

	if len(os.Args) >= 3 && strings.ToUpper(os.Args[1]) == "UNINSTALL" {
		systemdConfigFile, _ := filepath.Abs(os.Args[2])
		uninstallError := systemd.Uninstall(systemdConfigFile)
		if uninstallError != nil {
			panic(uninstallError)
		} else {
			log.Println("uninstalled")
		}
		return
	}
}
```

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