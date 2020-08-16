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
	"github.com/giantliao/beatles-protocol/token"
	"github.com/spf13/cobra"
)

var tokenconvertoggle bool

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
		if !tokenconvertoggle {
			beatlesTk, err := token.TokenCovert(args[0])
			if err != nil {
				fmt.Println("token covert error")
				return
			}

			fmt.Println(beatlesTk)
		} else {
			beatlesTk := args[0]

			if beatlesTk[:2] != "at" {
				fmt.Println("not a correct token")
				return
			}

			fmt.Println(token.TokenRevert(beatlesTk))
		}
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
	tokenconvertCmd.Flags().BoolVarP(&tokenconvertoggle, "revert", "r", false, "revert")
}
