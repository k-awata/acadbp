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
	Example: `  acadbp script example.scr *.dwg`,
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		batcher := acadbp.NewBatcher(viper.GetString("accorepath"))
		cobra.CheckErr(batcher.CheckAccore())

		files, err := acadbp.ExpandGlobPattern(args[1:])
		cobra.CheckErr(err)

		if viper.GetString("log") != "" {
			log, err := os.OpenFile(
				viper.GetString("log"),
				os.O_WRONLY|os.O_CREATE|os.O_APPEND,
				os.ModePerm,
			)
			cobra.CheckErr(err)
			defer log.Close()
			batcher.SetOutput(log)
		}

		// Read file, or stdio if arg is "-"
		scr := ""
		if args[0] == "-" {
			scr, err = acadbp.CreateTempFile("*.scr", acadbp.StdinToString(), viper.GetString("encoding"))
			cobra.CheckErr(err)
		} else {
			scr, err = filepath.Abs(args[0])
			cobra.CheckErr(err)
		}

		batcher.RunForEach(scr, files, "")
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
