/*
Copyright © 2022 mycok <github.com/mycok>

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
	"io"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List All Todo Items",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL := viper.GetString("api-root")

		return listAction(os.Stdout, rootURL)
	},
}

func listAction(w io.Writer, url string) error {
	items, err := getAll(url)

	if err != nil {
		return err
	}

	return printItems(w, items)
}

func printItems(w io.Writer, items []item) error {
	tw := tabwriter.NewWriter(w, 3, 2, 0, ' ', 0)

	for i, item := range items {
		done := "𝘅"

		if item.Done {
			done = "✅"
		}

		fmt.Fprintf(tw, "%s\t%d\t%s\t\n", done, i+1, item.Task)
	}

	return tw.Flush()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
