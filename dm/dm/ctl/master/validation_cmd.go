// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package master

import "github.com/spf13/cobra"

func NewValidationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "validation",
		Short:  "operate or query validation task",
		Hidden: true, // hide all validation command before formally GA
	}
	cmd.AddCommand(
		NewStartValidationCmd(),
		NewStopValidationCmd(),
		NewQueryValidationErrorCmd(),
		NewQueryValidationStatusCmd(),
		NewIgnoreValidationErrorCmd(),
		NewResolveValidationErrorCmd(),
		NewClearValidationErrorCmd(),
	)
	return cmd
}
