//go:build !windows

package capture

import "os/exec"

func configureHidden(c *exec.Cmd) {}
