package systemd

import (
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/senior-cyber/utility-systemd/dto"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//go:embed systemd.service
var systemd string

const systemdPath = "/etc/systemd/system"

func Install(systemConfigFile string, appConfigFile string) error {
	_cli, _ := filepath.Abs(os.Args[0])
	_workspace := filepath.Dir(_cli)

	config, configError := readConfig(systemConfigFile)
	if configError != nil {
		println(configError.Error())
		return configError
	}

	if config.Name == "" {
		println("systemd is required")
		return errors.New("systemd is required")
	}

	var _user = "root"
	if config.User != "" {
		_user = config.User
	}

	var _group = "root"
	if config.Group != "" {
		_group = config.Group
	}

	_name := config.Name
	_systemd := systemd
	_systemd = strings.ReplaceAll(_systemd, "{{name}}", _name)
	_systemd = strings.ReplaceAll(_systemd, "{{workspace}}", _workspace)
	_systemd = strings.ReplaceAll(_systemd, "{{cli}}", _cli)
	_systemd = strings.ReplaceAll(_systemd, "{{user}}", _user)
	_systemd = strings.ReplaceAll(_systemd, "{{group}}", _group)
	_systemd = strings.ReplaceAll(_systemd, "{{config}}", appConfigFile)
	_systemdError := os.WriteFile(filepath.Join(systemdPath, _name+".service"), []byte(_systemd), 0755)
	if _systemdError != nil {
		println(_systemdError.Error())
		return _systemdError
	}
	cmdSudo1 := exec.Command("sudo", "systemctl", "daemon-reload")
	_, _ = cmdSudo1.CombinedOutput()
	time.Sleep(1 * time.Second)
	cmdSudo2 := exec.Command("sudo", "systemctl", "enable", _name)
	_, _ = cmdSudo2.CombinedOutput()
	time.Sleep(1 * time.Second)
	cmdSudo3 := exec.Command("sudo", "systemctl", "start", _name)
	_, _ = cmdSudo3.CombinedOutput()
	time.Sleep(1 * time.Second)
	return nil
}

func Uninstall(systemConfigFile string) error {
	config, configError := readConfig(systemConfigFile)
	if configError != nil {
		println(configError.Error())
		return configError
	}

	if config.Name == "" {
		println("systemd is required")
		return errors.New("systemd is required")
	}

	_name := config.Name

	cmdSudo1 := exec.Command("sudo", "systemctl", "stop", _name)
	_, _ = cmdSudo1.CombinedOutput()
	time.Sleep(1 * time.Second)
	cmdSudo2 := exec.Command("sudo", "systemctl", "disable", _name)
	_, _ = cmdSudo2.CombinedOutput()
	time.Sleep(1 * time.Second)
	_ = os.Remove(filepath.Join(systemdPath, _name+".service"))
	cmdSudo3 := exec.Command("sudo", "systemctl", "daemon-reload")
	_, _ = cmdSudo3.CombinedOutput()
	time.Sleep(1 * time.Second)
	return nil
}

func readConfig(_config string) (*dto.SystemdDto, error) {
	configData, configDataError := os.ReadFile(_config)
	if configDataError != nil {
		return nil, configDataError
	}
	var config dto.SystemdDto
	configError := json.Unmarshal(configData, &config)
	if configError != nil {
		return nil, configError
	}
	return &config, nil
}
