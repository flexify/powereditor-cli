// Copyright © 2017 flexify.net
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/dommmel/goshopping/shopify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const rowSeparator string = "<!--|row|-->"
const colSeparator string = "<!--|col|-->"

// init global flags
var cfgFile, outputFile string

//var debug = Debug("cli")

var RootCmd = &cobra.Command{
	Use:   "powereditor-cli",
	Short: "powereditor-cli is a tool to export/import power-editor content from/to Shopify",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is \"config.yml\")")
	RootCmd.PersistentFlags().StringP("key", "k", "", "shopify api key. This will override what is in your config.yml")
	RootCmd.PersistentFlags().StringP("password", "p", "", "shopify api password. This will override what is in your config.yml")
	RootCmd.PersistentFlags().StringP("store", "s", "", "your shopify domain. This will override what is in your config.yml")
	RootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "output.json", "the file the results are written to")
	RootCmd.PersistentFlags().StringP("namespace", "n", "power-editor", "the metafield namespace. This will override what is in your config.yml")
	viper.BindPFlag("export.namespace", RootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("import.namespace", RootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("export.key", RootCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("export.password", RootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("export.store", RootCmd.PersistentFlags().Lookup("store"))
	viper.BindPFlag("import.key", RootCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("import.password", RootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("import.store", RootCmd.PersistentFlags().Lookup("store"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func checkGlobalRequiredFlags() []string {

	// Check for required API credentials
	var errorMsg []string
	if viper.GetString("export.key") == "" {
		errorMsg = append(errorMsg, "api key is required")
	}
	if viper.GetString("export.password") == "" {
		errorMsg = append(errorMsg, "api password is required")
	}
	if viper.GetString("export.store") == "" {
		errorMsg = append(errorMsg, "store domain is required")
	}
	return errorMsg
}

func GetClient() *shopify.Client {
	return shopify.NewPrivateClient(nil, viper.GetString("export.key"), viper.GetString("export.password"), viper.GetString("export.store"))
}

func getSliceOfMapValue(m map[string]string) []string {
	// Preserve Order of map entries
	// See: https://blog.golang.org/go-maps-in-action#TOC_7.
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var v []string

	for _, k := range keys {
		v = append(v, m[k])
	}
	return v
}

/* DATA FORMAT */

type Output struct {
	Products []*ProductOutput `json:"products"`
}

type OutputField struct {
	Id   *int                         `json:"id"`
	Key  *string                      `json:"key"`
	Data map[string]map[string]string `json:"data"`
}

type ProductOutput struct {
	Id                             *int           `json:"id"`
	Handle                         *string        `json:"handle"`
	BodyHtml                       *string        `json:"body_html,omitempty"`
	Title                          *string        `json:"title,omitempty"`
	MetafieldsGlobalTitleTag       *string        `json:"metafields_global_title_tag,omitempty"`
	MetafieldsGlobalDescriptionTag *string        `json:"metafields_global_description_tag,omitempty"`
	Fields                         []*OutputField `json:"fields,omitempty"`
}
