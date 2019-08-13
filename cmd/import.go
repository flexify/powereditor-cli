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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/caarlos0/spin"
	"github.com/dommmel/goshopping/shopify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileName string
var allowedPrimaryKeys map[string]bool

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a previously exported data dump",

	// Do all the error handling pre run
	PreRunE: func(cmd *cobra.Command, args []string) error {

		// Check for required API credentials
		errorMsg := checkGlobalRequiredFlags()

		key := viper.GetString("import.primary-key")
		if !allowedPrimaryKeys[key] && key != "id" {
			errorMsg = append(errorMsg, "primary key '"+key+"' is not valid")
		}
		// Check for required collection ID
		if len(args) < 1 {
			errorMsg = append(errorMsg, "path to data file required as an argument")
		} else {
			fileName = args[0]
			if _, err := os.Stat(fileName); os.IsNotExist(err) {
				errorMsg = append(errorMsg, fmt.Sprintf("Can't access file. %v", err))
			}
		}

		if len(errorMsg) == 1 {
			return errors.New(errorMsg[0])
		} else if len(errorMsg) > 0 {
			return errors.New("\n - " + strings.Join(errorMsg, "\n - "))
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		client := GetClient()
		file, _ := ioutil.ReadFile(fileName)
		var data Output
		json.Unmarshal(file, &data)
		// Loop over all products
		for i, p := range data.Products {
			progress := fmt.Sprintf("%d of %d", i, len(data.Products))
			s := spin.New("  \033[36m Importing product " + progress + "\033[m %s")
			s.Set(spin.Spin1)
			s.Start()

			var productId *int

			// Get ID of the product whose metafields will be updated
			if key := viper.GetString("import.primary-key"); allowedPrimaryKeys[key] {

				// Get the key (handle or title) as string
				f := reflect.ValueOf(p).Elem().FieldByName(strings.Title(key))
				keyValue := reflect.Indirect(f).String()

				var noproducterr error
				productId, noproducterr = getProductIdByProperty(key, keyValue, client)

				if noproducterr != nil {
					fmt.Printf("Skipping %s:%s", key, keyValue)
					continue
				}

				fmt.Printf("%s: %s => id: %d\n", key, keyValue, *productId)

			} else {
				productId = p.Id
			}

			// debug("Product id: %d", *productId)
			// debug("Body: %s", *p.BodyHtml)
			updatedProduct := &shopify.Product{
				Id: productId,
				// Handle:     p.Handle,
				Title:      p.Title,
				BodyHtml:   p.BodyHtml,
				Metafields: AssembleMetafieldData(p.Fields, client),
			}

			_, err := client.Products.Edit(context.Background(), updatedProduct)
			// debug("resp %s", resp)
			if err != nil {
				fmt.Errorf("%s", err)
			}
			s.Stop()
		}
	},
}

func AssembleMetafieldData(fields []*OutputField, client *shopify.Client) (metafields []*shopify.Metafield) {

	for _, field := range fields {
		var rowsToMerge []string
		for _, row := range field.Data {
			colsToMerge := getSliceOfMapValue(row)
			rowsToMerge = append(rowsToMerge, strings.Join(colsToMerge, colSeparator))
		}
		metafieldValue := strings.Join(rowsToMerge, rowSeparator)

		valueType := "string"
		ns := viper.GetString("import.namespace")
		out := &shopify.Metafield{
			Namespace: &ns,
			Key:       field.Key,
			Id:        field.Id,
			Value:     &metafieldValue,
			ValueType: &valueType,
		}

		metafields = append(metafields, out)
	}
	return
}

func getProductIdByProperty(propertyName string, propertyValue string, client *shopify.Client) (*int, error) {

	ctx := context.Background()
	opt := &shopify.ProductListOptions{Fields: []string{"id", "metafields"}}

	f := reflect.Indirect(reflect.ValueOf(opt)).FieldByName(strings.Title(propertyName))
	f.SetString(propertyValue)

	// if there is no ID look up products by handle or title
	products, _, err := client.Products.List(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("Can't find product with %s '%s': %s'", propertyName, propertyValue, err)
	}
	if len(products) > 1 {
		return nil, fmt.Errorf("Found more than on product for %s '%s': %s'", propertyName, propertyValue, err)
	}
	if len(products) == 0 {
		return nil, fmt.Errorf("Found no product  '%s': %s'", propertyName, propertyValue, err)
	}
	//spew.Dump(products[0])
	return products[0].Id, nil
}

func init() {
	allowedPrimaryKeys = map[string]bool{"handle": true, "title": true}
	RootCmd.AddCommand(importCmd)
	importCmd.Flags().BoolP("metafields-only", "m", false, "Don't import product titles or descriptions")
	// importCmd.Flags().BoolP("dry-run", "d", false, "Do not import but show a list of updates that would happen")
	importCmd.Flags().StringP("primary-key", "1", "id", `Possible values are "id", "handle" and "title"`)
	viper.BindPFlag("import.primary-key", importCmd.Flags().Lookup("primary-key"))
	viper.BindPFlag("import.metafields-only", importCmd.Flags().Lookup("metafields-only"))
}
