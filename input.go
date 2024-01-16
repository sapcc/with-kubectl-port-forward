/*******************************************************************************
*
* Copyright 2022 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

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
