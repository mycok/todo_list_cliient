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
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:          "view <itemID>",
	Short:        "View a specific todo item with details",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL := viper.GetString("api-root")

		return viewAction(os.Stdout, rootURL, args[0])
	},
}

func viewAction(w io.Writer, url, id string) error {
	itemID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("%w: item ID must me a number", ErrNotNumber)
	}

	item, err := getItem(url, itemID)
	if err != nil {
		return err
	}

	return printItem(w, item)
}

func printItem(w io.Writer, i item) error {
	tw := tabwriter.NewWriter(w, 14, 2, 0, ' ', 0)

	fmt.Fprintf(tw, "Task:\t%s\n", i.Task)
	fmt.Fprintf(tw, "Created at:\t%s\n", i.CreatedAt.Format(timeFormat))

	if i.Done {
		fmt.Fprintf(tw, "Completed:\t%s\n", "Yes")
		fmt.Fprintf(tw, "CompletedAt:\t%s\n", i.CompletedAt.Format(timeFormat))

		return tw.Flush()
	}

	fmt.Fprintf(tw, "Completed:\t%s\n", "No")

	return tw.Flush()
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
