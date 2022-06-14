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

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:     "publish input_file...",
	Aliases: []string{"p"},
	Short:   "Publish input drawing files with specified page setup",
	Long:    "Publish input drawing files with specified page setup",
	Example: `  acadbp publish --setup-file setup.dwg --setup-name Setup1 *.dxf`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		files := acadbp.ExpandGlobPattern(args)

		// Create dsd file
		tmpl := "[DWF6Version]\r\nVer=1\r\n[DWF6MinorVersion]\r\nMinorVer=1\r\n"
		if viper.IsSet("publish.dsd") {
			var err error
			tmpl, err = acadbp.ReadTemplateDsd(viper.GetString("publish.dsd"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
		trg, err := acadbp.CreateDsdTarget(
			viper.GetString("publish.type"),
			viper.GetString("publish.multi"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		shts, err := acadbp.CreateDsdSheets(
			files,
			viper.GetString("publish.setup-name"),
			viper.GetString("publish.setup-file"),
			viper.GetString("publish.layout"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		dsd, err := acadbp.CreateTempFile("*.dsd", tmpl+trg+shts)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		scr, err := acadbp.CreateTempFile("*.scr", "_.-publish "+dsd+"\r\n")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		bat, err := acadbp.CreateBatContents(viper.GetString("accorepath"), scr, viper.GetString("log"), nil)
		if err != nil {
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
	rootCmd.AddCommand(publishCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// publishCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// publishCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	publishCmd.Flags().StringP("dsd", "d", "", "template dsd file")
	publishCmd.Flags().StringP("setup-name", "s", "", "page setup name")
	publishCmd.Flags().StringP("setup-file", "f", "", "drawing file that includes page setup")
	publishCmd.Flags().StringP("layout", "l", "Model", "layout name or model to publish")
	publishCmd.Flags().StringP("type", "t", "pdf", "output type (plotter|dwf|dwfx|pdf)")
	publishCmd.Flags().StringP("multi", "m", "", "multi-sheet file name")

	viper.BindPFlag("publish.dsd", publishCmd.Flags().Lookup("dsd"))
	viper.BindPFlag("publish.setup-name", publishCmd.Flags().Lookup("setup-name"))
	viper.BindPFlag("publish.setup-file", publishCmd.Flags().Lookup("setup-file"))
	viper.BindPFlag("publish.layout", publishCmd.Flags().Lookup("layout"))
	viper.BindPFlag("publish.type", publishCmd.Flags().Lookup("type"))
	viper.BindPFlag("publish.multi", publishCmd.Flags().Lookup("multi"))
}
