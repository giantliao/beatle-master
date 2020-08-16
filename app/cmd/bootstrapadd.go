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
	"github.com/giantliao/beatles-protocol/token"
	"log"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var (
	bootstrapOwner       string
	bootstrapRepository  string
	bootstrapPath        string
	bootstrapReadToken   string
	bootstrapCommitName  string
	bootstrapCommitEmail string
)

var bootstrapaddCmd = &cobra.Command{
	Use:   "add",
	Short: "add bootstrap server",
	Long:  `add bootstrap server`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := cmdcommon.IsProcessStarted(); err != nil {
			log.Println(err)
			return
		}

		if bootstrapOwner == "" || bootstrapRepository == "" || bootstrapPath == "" ||
			bootstrapReadToken == "" || bootstrapCommitEmail == "" || bootstrapCommitName == "" {
			log.Println("please enter bootstrap server")
			return
		}

		if bootstrapReadToken[:2] != "at" {
			var err error
			bootstrapReadToken, err = token.TokenCovert(bootstrapReadToken)
			if err != nil {
				log.Println("bootstrap readtoken error")
				return
			}
		}

		var param []string
		param = append(param, bootstrapOwner, bootstrapRepository, bootstrapPath, bootstrapReadToken, bootstrapCommitName, bootstrapCommitEmail)

		cmdclient.StringOpCmdSend("", cmdcommon.CMD_BOOTSTRAP_ADD, param)

	},
}

func init() {
	bootstrapCmd.AddCommand(bootstrapaddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	bootstrapaddCmd.Flags().StringVarP(&bootstrapOwner, "owner", "o", "", "repository owner")
	bootstrapaddCmd.Flags().StringVarP(&bootstrapRepository, "repository", "r", "", "repository name")
	bootstrapaddCmd.Flags().StringVarP(&bootstrapPath, "path", "p", "", "file path")
	bootstrapaddCmd.Flags().StringVarP(&bootstrapReadToken, "token", "t", "", "token for read bootstrap file")
	bootstrapaddCmd.Flags().StringVarP(&bootstrapCommitName, "commitname", "n", "", "name for commit")
	bootstrapaddCmd.Flags().StringVarP(&bootstrapCommitEmail, "commitemail", "e", "", "email for commit")

}
