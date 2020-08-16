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
	"github.com/btcsuite/btcutil/base58"
	"github.com/spf13/cobra"
	"github.com/status-im/keycard-go/hexutils"
)

// tokenconvertCmd represents the tokenconvert command
var tokenconvertCmd = &cobra.Command{
	Use:   "tokenconvert",
	Short: "convert token to a base58 code",
	Long:  `convert token to a base58 code`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("please input token with hash code")
			return
		}

		tokenbyte := hexutils.HexToBytes(args[0])

		fmt.Println("at" + base58.Encode(tokenbyte))
	},
}

func init() {
	rootCmd.AddCommand(tokenconvertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tokenconvertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tokenconvertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
