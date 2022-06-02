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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete <itemID>",
	Short: "Mark a todo item as complete",
	Args: cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL := viper.GetString("api-root")

		return completeAction(os.Stdout, rootURL, args[0])
	},
}

func completeAction(w io.Writer, url, id string) error {
	itemId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("%w: item ID must me a number", ErrNotNumber)
	}

	if err := completeItem(url, itemId); err != nil {
		return err
	}

	return printCompletedItem(w, itemId)
}

func printCompletedItem(w io.Writer, id int) error {
	_, err := fmt.Fprintf(w, "Item number %d marked as complete\n", id)

	return err
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
