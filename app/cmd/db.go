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
	"fmt"
	"github.com/giantliao/beatles-master/app/cmdclient"
	"github.com/giantliao/beatles-master/app/cmdcommon"
	"github.com/giantliao/beatles-master/config"
	"github.com/spf13/cobra"
	"log"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "show db",
	Long:  `show db`,
	Run: func(cmd *cobra.Command, args []string) {

		dbs := config.GetCBtlm().GetDbs()

		if len(args) == 0 {
			fmt.Println(dbs)
			return
		}

		if _, err := cmdcommon.IsProcessStarted(); err != nil {
			log.Println(err)
			return
		}

		if len(args) > 1 {
			fmt.Println("parameter error")
			return
		}
		dbName := ""
		for i := 0; i < len(dbs); i++ {
			if dbs[i] == args[0] {
				dbName = args[0]
				break
			}
		}

		if len(dbName) == 0 {
			fmt.Println("db name not correct")
			return
		}

		var param []string
		param = append(param, dbName)

		cmdclient.StringOpCmdSend("", cmdcommon.CMD_DB_SHOW, param)
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
