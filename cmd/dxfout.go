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

	"github.com/k-awata/acadbp/acadbp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dxfoutCmd represents the dxfout command
var dxfoutCmd = &cobra.Command{
	Use:     "dxfout input_file...",
	Aliases: []string{"dxf", "x"},
	Short:   "Convert input drawing files to DXF files",
	Example: `  acadbp dxfout *.dwg
  acadbp dxfout --format 2018 --dp Binary --preview *.dwg`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		batcher := acadbp.NewBatcher(viper.GetString("accorepath"))
		cobra.CheckErr(batcher.CheckAccore())

		files, err := acadbp.ExpandGlobPattern(args)
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

		scr, err := acadbp.CreateTempFile(
			"*.scr",
			"_.saveas DXF P "+acadbp.BtoYN(viper.GetBool("dxf.preview"))+
				" V "+viper.GetString("dxf.format")+
				" "+viper.GetString("dxf.dp")+
				" \ny\nfiledia "+acadbp.Bto10(viper.GetBool("filedia"))+"\n",
			viper.GetString("encoding"),
		)
		cobra.CheckErr(err)
		defer os.Remove(scr)

		batcher.RunForEach(scr, files, ".dxf")
	},
}

func init() {
	rootCmd.AddCommand(dxfoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dxfoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dxfoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	dxfoutCmd.Flags().StringP("format", "f", "2013", "file format version")
	dxfoutCmd.Flags().StringP("dp", "d", "16", "decimal places of accuracy (0 to 16 or Binary)")
	dxfoutCmd.Flags().BoolP("preview", "p", false, "save thumbnail preview image")
	viper.BindPFlag("dxf.format", dxfoutCmd.Flags().Lookup("format"))
	viper.BindPFlag("dxf.dp", dxfoutCmd.Flags().Lookup("dp"))
	viper.BindPFlag("dxf.preview", dxfoutCmd.Flags().Lookup("preview"))
}
