/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
)

// saveCmd represents the save command
var dbsaveCmd = &cobra.Command{
	Use:   "save",
	Short: "save db",
	Long:  `save db`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("please set a database name")
			return
		}

		dbname := ""
		dbs := config.GetDbs()

		for i := 0; i < len(dbs); i++ {
			if args[0] == dbs[i] {
				dbname = args[0]
			}
		}

		if dbname == "" {
			fmt.Println("no database name")
			return
		}

		var param []string
		param = append(param, dbname)

		cmdclient.StringOpCmdSend("", cmdcommon.CMD_DB_SAVE, param)

	},
}

func init() {
	dbCmd.AddCommand(dbsaveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
