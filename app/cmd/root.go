// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-master/wallet"
	"github.com/giantliao/beatles-master/webserver"
	"github.com/howeyc/gopass"
	"os"

	"github.com/giantliao/beatles-master/app/cmdcommon"
	"github.com/giantliao/beatles-master/config"

	"github.com/giantliao/beatles-master/app/cmdservice"

	"github.com/spf13/cobra"
	"log"
)

var (
	cmdconfigfilename string
)

var keypassword string

func inputpassword() (password string, err error) {
	passwd, err := gopass.GetPasswdPrompt("Please Enter Password: ", true, os.Stdin, os.Stdout)
	if err != nil {
		return "", err
	}

	if len(passwd) < 1 {
		return "", errors.New("Please input valid password")
	}

	return string(passwd), nil
}

func inputChoose() (choose string, err error) {
	c, err := gopass.GetPasswdPrompt("Do you reinit config[yes/no]: ", true, os.Stdin, os.Stdout)
	if err != nil {
		return "", err
	}

	return string(c), nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "btlm",
	Short: "start beatles master in current shell",
	Long:  `start beatles master in current shell`,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := cmdcommon.IsProcessCanStarted()
		if err != nil {
			log.Println(err)
			return
		}

		InitCfg()
		cfg := config.GetCBtlm()
		cfg.Save()

		if cfg.EthAccessPoint == "" {
			log.Println("please init first")
			return
		}

		if keypassword == "" {
			if keypassword, err = inputpassword(); err != nil {
				log.Println(err)
				return
			}
		}

		err = wallet.LoadWallet(keypassword)
		if err != nil {
			panic("load wallet failed")
		}

		go webserver.StartWebDaemon()
		go db.GetMinersDb().TimeOut()

		cmdservice.GetCmdServerInst().StartCmdService()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func InitCfg() {
	if cmdconfigfilename != "" {
		cfg := config.LoadFromCfgFile(cmdconfigfilename)
		if cfg == nil {
			return
		}
	} else {
		config.LoadFromCmd(cfginit)
	}

}

func cfginit(bc *config.BtlMasterConf) *config.BtlMasterConf {
	cfg := bc
	if remoteethaccesspoint != "" {
		cfg.EthAccessPoint = remoteethaccesspoint
	}
	if remotetrxaccesspoint != "" {
		cfg.TrxAccessPoint = remotetrxaccesspoint
	}

	if remotebtlcaccesspoint != ""{
		cfg.BTLCAccessPoint = remotebtlcaccesspoint
	}

	if btlccontractaddr != ""{
		cfg.BTLCoinAddr = btlccontractaddr
	}

	return cfg

}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringVarP(&cmdconfigfilename, "config-file-name", "c", "", "configuration file name")
	rootCmd.Flags().StringVarP(&remoteethaccesspoint, "eth", "e", "", "eth access point")
	rootCmd.Flags().StringVarP(&remotetrxaccesspoint, "trx", "t", "", "tron network access point")
}
