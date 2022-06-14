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

	"github.com/k-awata/acadbp/acadbp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dwgoutCmd represents the dwgout command
var dwgoutCmd = &cobra.Command{
	Use:     "dwgout input_file...",
	Aliases: []string{"dwg", "w"},
	Short:   "Convert input drawing files to DWG files",
	Long:    "Convert input drawing files to DWG files",
	Example: `  acadbp dwgout *.dxf
  acadbp dwgout --format 2010 *.dwg`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		files := acadbp.ExpandGlobPattern(args)

		scr, err := acadbp.CreateTempFile("*.scr", "_.saveas "+viper.GetString("dwg.format")+" \r\nY\r\n", viper.GetBool("sjis"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		bat, err := acadbp.CreateBatContents(viper.GetString("accorepath"), scr, viper.GetString("log"), files)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		if err := acadbp.CreateEmptyFiles(files, ".dwg"); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		if err := acadbp.RunBat(bat); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(dwgoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dwgoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dwgoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	dwgoutCmd.Flags().StringP("format", "f", "2013", "file format version")
	viper.BindPFlag("dwg.format", dwgoutCmd.Flags().Lookup("format"))
}
