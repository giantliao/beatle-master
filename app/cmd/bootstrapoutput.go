/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"github.com/giantliao/beatles-master/app/cmdclient"
	"github.com/giantliao/beatles-master/app/cmdcommon"
	"github.com/spf13/cobra"
	"log"
)

var (
	bootstrapOutputFile string
	//bootstrapToGithub bool
)

// outputCmd represents the output command
var bootstrapoutputCmd = &cobra.Command{
	Use:   "output",
	Short: "output bootstrap to a file ",
	Long:  `output bootstrap to a file, include bootstrap server, miners...`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := cmdcommon.IsProcessStarted(); err != nil {
			log.Println(err)
			return
		}

		var param []string
		param = append(param, bootstrapOutputFile)

		cmdclient.StringOpCmdSend("", cmdcommon.CMD_BOOTSTRAP_LIST, param)

	},
}

func init() {
	bootstrapCmd.AddCommand(bootstrapoutputCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// outputCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// outputCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	bootstrapoutputCmd.Flags().StringVarP(&bootstrapOutputFile, "filename", "f", "bootstrap.list", "output a bootstrap server list")
	//bootstrapoutputCmd.Flags().BoolVarP(&bootstrapToGithub,"togithub","t",false,"push bootstrap to github")

}
