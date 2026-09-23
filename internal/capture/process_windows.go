//go:build windows

package capture

import (
	"os/exec"
	"syscall"
)

func configureHidden(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
}
