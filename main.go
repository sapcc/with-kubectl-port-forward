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
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sapcc/go-bits/errext"
)

func main() {
	// parse command line
	portForwardArgs, subcommandArgs := splitArgs()

	// setup async machinery
	ctx, cancel := signal.NotifyContext(context.Background(), // once one subprocess returns, we cancel this ctx to reap the other one
		os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup
	errChan := make(chan error, 2)          // collects errors from the subprocess; any error from one will terminate the whole program
	portReadableChan := make(chan struct{}) // signals that the ports are established and the subcommand can start

	// run subprocesses
	wg.Add(2)
	go func() {
		defer wg.Done()
		runKubectlPortForward(ctx, portForwardArgs, errChan, portReadableChan)
	}()
	go func() {
		defer wg.Done()
		runSubcommand(ctx, subcommandArgs, errChan, portReadableChan)
	}()

	// once either exits, cancel the other
	err := <-errChan
	cancel()
	wg.Wait()

	if err == nil {
		os.Exit(0)
	} else if exitErr, ok := errext.As[*exec.ExitError](err); ok {
		os.Exit(exitErr.ExitCode())
	} else {
		fmt.Fprintln(os.Stderr, "error: "+err.Error())
		os.Exit(1)
	}
}
