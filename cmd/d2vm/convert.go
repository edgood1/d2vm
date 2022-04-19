// Copyright 2022 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/c2h5oh/datasize"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"go.linka.cloud/d2vm"
	"go.linka.cloud/d2vm/pkg/docker"
	"go.linka.cloud/d2vm/pkg/exec"
)

var (
	convertCmd = &cobra.Command{
		Use:          "convert [docker image]",
		Short:        "Convert Docker image to qcow2 vm image",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			img := args[0]
			size, err := parseSize(size)
			if err != nil {
				return err
			}
			if _, err := os.Stat(output); err == nil || !os.IsNotExist(err) {
				if !force {
					return fmt.Errorf("%s already exists", output)
				}
			}
			if debug {
				exec.Run = exec.RunStdout
			}
			if _, err := os.Stat(output); err == nil || !os.IsNotExist(err) {
				if !force {
					return fmt.Errorf("%s already exists", output)
				}
			}
			logrus.Infof("pulling image %s", img)
			if err := docker.Cmd(cmd.Context(), "image", "pull", img); err != nil {
				return err
			}
			return docker2vm.Convert(cmd.Context(), img, size, password, output)
		},
	}
)

func parseSize(s string) (int64, error) {
	var v datasize.ByteSize
	if err := v.UnmarshalText([]byte(s)); err != nil {
		return 0, err
	}
	return int64(v), nil
}

func init() {
	convertCmd.Flags().StringVarP(&output, "output", "o", output, "The output qcow2 image")
	convertCmd.Flags().StringVarP(&password, "password", "p", "root", "The Root user password")
	convertCmd.Flags().StringVarP(&size, "size", "s", "1G", "The output image size")
	convertCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable Debug output")
	convertCmd.Flags().BoolVarP(&force, "force", "f", false, "Override output qcow2 image")
	rootCmd.AddCommand(convertCmd)
}
