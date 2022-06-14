/*
Copyright © 2022 K.Awata

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initconfCmd represents the initconf command
var initconfCmd = &cobra.Command{
	Use:   "initconf accoreconsole_path",
	Short: "Create acadbp config file to home directory",
	Long:  "Create acadbp config file to home directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		abs, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if !strings.EqualFold(filepath.Base(abs), "accoreconsole.exe") {
			fmt.Fprintln(os.Stderr, "given path is not to accoreconsole.exe")
			return
		}
		viper.Set("accorepath", abs)
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		conf := filepath.Join(home, ".acadbp.yaml")
		if _, err := os.Stat(conf); err != nil {
			if _, err := os.Create(conf); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
		if err := viper.WriteConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initconfCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initconfCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initconfCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
