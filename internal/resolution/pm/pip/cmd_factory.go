package pip

import (
	"os/exec"
	"strings"
)

type ICmdFactory interface {
	MakeCreateVenvCmd(file string) (*exec.Cmd, error)
	MakeInstallCmd(command string, file string) (*exec.Cmd, error)
	MakeCatCmd(file string) (*exec.Cmd, error)
	MakeListCmd(command string) (*exec.Cmd, error)
	MakeShowCmd(command string, list []string) (*exec.Cmd, error)
}

type IExecPath interface {
	LookPath(file string) (string, error)
}

type ExecPath struct {
}

func (_ ExecPath) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

type CmdFactory struct {
	execPath IExecPath
}

func (cmdf CmdFactory) MakeCreateVenvCmd(fpath string) (*exec.Cmd, error) {
	python, err := cmdf.execPath.LookPath("python3")
	pythonCommand := "python3"
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found in ") {
			// Python 3 not found, try Python
			python, err = cmdf.execPath.LookPath("python")
			pythonCommand = "python"
		}

		// If error still is != nil, return
		if err != nil {
			return nil, err
		}
	}

	return &exec.Cmd{
		Path: python,
		Args: []string{pythonCommand, "-m", "venv", fpath, "--clear", "--system-site-packages"},
	}, nil
}

func (cmdf CmdFactory) MakeInstallCmd(command string, file string) (*exec.Cmd, error) {
	path, err := cmdf.execPath.LookPath(command)

	return &exec.Cmd{
		Path: path,
		Args: []string{command, "install", "-r", file},
	}, err
}

func (cmdf CmdFactory) MakeCatCmd(file string) (*exec.Cmd, error) {
	path, err := cmdf.execPath.LookPath("cat")

	return &exec.Cmd{
		Path: path,
		Args: []string{"cat", file},
	}, err
}

func (cmdf CmdFactory) MakeListCmd(command string) (*exec.Cmd, error) {
	path, err := cmdf.execPath.LookPath(command)

	return &exec.Cmd{
		Path: path,
		Args: []string{"pip", "list"},
	}, err
}

func (cmdf CmdFactory) MakeShowCmd(command string, list []string) (*exec.Cmd, error) {
	path, err := cmdf.execPath.LookPath(command)

	args := []string{command, "show"}
	args = append(args, list...)

	return &exec.Cmd{
		Path: path,
		Args: args,
	}, err
}
