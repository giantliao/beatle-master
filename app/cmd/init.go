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
	"github.com/giantliao/beatles-master/app/cmdcommon"
	"github.com/giantliao/beatles-master/config"

	"github.com/spf13/cobra"
	"log"
)

var (
	remotetrxaccesspoint  string
	remoteethaccesspoint  string
	remotebtlcaccesspoint string
	btlccontractaddr      string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init beatles master",
	Long:  `init beatles master`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		_, err = cmdcommon.IsProcessCanStarted()
		if err != nil {
			log.Println(err)
			return
		}

		InitCfg()

		if remoteethaccesspoint == "" || remotebtlcaccesspoint == "" || btlccontractaddr == "" {
			fmt.Println("eth access point and btlc access point must set")
			return
		}

		cfg := config.GetCBtlm()

		cfg.Save()

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")
	initCmd.Flags().StringVarP(&remoteethaccesspoint, "eth", "e", "", "eth access point")
	initCmd.Flags().StringVarP(&remotetrxaccesspoint, "trx", "t", "", "tron network access point")
	initCmd.Flags().StringVarP(&remotebtlcaccesspoint, "btlc", "c", "", "btlc access point")
	initCmd.Flags().StringVarP(&btlccontractaddr, "baddr", "b", "", "btlc contract address")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
