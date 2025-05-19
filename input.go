// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
)

func splitArgs() (argsForPortForward, cmdline []string) {
	doubleDashIndex := -1
	for idx, arg := range os.Args {
		if idx == 0 {
			continue
		}
		if arg == "--" {
			doubleDashIndex = idx
			break
		}
		if arg == "--help" {
			usage(0)
		}
	}

	if doubleDashIndex == -1 {
		usageError("missing `--` in argument list")
	}
	if doubleDashIndex == 1 {
		usageError("missing arguments for `kubectl port-forward` (you need to put something before `--`)")
	}
	if doubleDashIndex == len(os.Args)-1 {
		usageError("missing command line (you need to put something after `--`)")
	}

	return os.Args[1:doubleDashIndex], os.Args[doubleDashIndex+1:]
}

func usage(status int) {
	fmt.Fprintf(os.Stderr, "usage: %s <port-forward-arg>... -- <command> <arg>...", os.Args[0])
	os.Exit(status)
}

func usageError(msg string) {
	fmt.Fprintln(os.Stderr, "ERROR: "+msg)
	usage(1)
}
