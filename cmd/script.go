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

	"github.com/k-awata/acadbp/acadbp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scriptCmd represents the script command
var scriptCmd = &cobra.Command{
	Use:     "script scr_file input_file...",
	Aliases: []string{"scr", "s"},
	Short:   "Run script file for each input file",
	Long:    "Run script file for each input file",
	Example: `  acadbp script example.scr *.dwg`,
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := acadbp.CheckAcCorePath(viper.GetString("accorepath")); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		files, err := acadbp.ExpandGlobPattern(args[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Read file, or stdio if arg is "-"
		scr := ""
		if args[0] == "-" {
			scr, err = acadbp.CreateTempFile("*.scr", acadbp.StdinToString(), viper.GetString("encoding"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		} else {
			scr, err = filepath.Abs(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}

		bat := acadbp.CreateBatContents(viper.GetString("accorepath"), scr, viper.GetString("log"), files)
		if err := acadbp.RunBatCommands(bat, viper.GetString("encoding")); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(scriptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scriptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scriptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
