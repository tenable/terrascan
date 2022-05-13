//go:build windows
// +build windows

/*
    Copyright (C) 2022 Tenable, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package cli

import (
	"os"

	"golang.org/x/sys/windows"
)

// In order for colored output to work on Windows, we need to explicitly
// enable the virtual terminal processing console mode.  This is needed for
// zap (log) output that goes to the console as well as program output.
func init() {
	var originalMode uint32

	stdout := windows.Handle(os.Stdout.Fd())

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}
