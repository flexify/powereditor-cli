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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/spin"
	"github.com/dommmel/goshopping/shopify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var collectionId int
var namespace string

// collectionCmd represents the "collection export" command
var collectionCmd = &cobra.Command{
	Use:   "collection <collection id>",
	Short: "export a collection's power-editor content",

	// Do all the error handling pre run
	PreRunE: func(cmd *cobra.Command, args []string) error {

		// Check for required API credentials
		errorMsg := checkGlobalRequiredFlags()

		// Check for required collection ID
		if len(args) < 1 {
			errorMsg = append(errorMsg, "a collection ID is required as an argument")
		} else {
			var err error
			collectionId, err = strconv.Atoi(args[0])
			return err
		}

		if len(errorMsg) == 1 {
			return errors.New(errorMsg[0])
		} else if len(errorMsg) > 0 {
			return errors.New("\n - " + strings.Join(errorMsg, "\n - "))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		s := spin.New("  \033[36m Scanning collection \033[m %s")
		s.Set(spin.Spin1)
		s.Start()
		defer fmt.Println("== Exported to", outputFile)

		products, err := GetProductsByCollection(collectionId, client)
		s.Stop()
		if err != nil {
			return err
		}

		var output Output
		for i, product := range products {
			progress := fmt.Sprintf("%d of %d", i, len(products))
			s = spin.New("  \033[36m Fetching product " + progress + "\033[m %s")
			s.Set(spin.Spin1)
			s.Start()

			metafields, _ := GetMetafieldsByProduct(*product.Id, viper.GetString("export.namespace"), client)

			// Add this product if it has metafields that should be exported or if the default product information should be included
			exportThisProduct := len(metafields) > 0 || viper.GetBool("export.include-product-info")

			if exportThisProduct {
				globalTitleTag, globalDescriptionTag := getSeoTagsByProduct(*product.Id, client)
				outputFields := GenerateProductDataOutput(metafields)

				// Fill in the output data
				pout := &ProductOutput{
					Id:                             product.Id,
					Handle:                         product.Handle,
					Title:                          product.Title,
					MetafieldsGlobalTitleTag:       globalTitleTag,
					MetafieldsGlobalDescriptionTag: globalDescriptionTag,
					BodyHtml:                       product.BodyHtml,
					Fields:                         outputFields,
				}
				output.Products = append(output.Products, pout)
			}
			s.Stop()
		}

		writeToFile(output, outputFile)
		return nil
	},
}

func init() {
	// this is a subcommand to the "collection" command
	collectionCmd.Flags().BoolP("include-product-info", "i", false, "Include product content (titles, descriptions) in export")
	viper.BindPFlag("export.include-product-info", collectionCmd.Flags().Lookup("include-product-info"))
	exportCmd.AddCommand(collectionCmd)
}

func writeToFile(thingsToWrite interface{}, fileName string) {
	b, err := JSONMarshalIndent(thingsToWrite, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing json %s", err)
	}
	ioutil.WriteFile(fileName, b, 0644)
}

// https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and
// Todo: Break out in own package
func JSONMarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	b, err := JSONMarshal(v)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// Foo is my foo function
func GetProductsByCollection(collectionId int, client *shopify.Client) ([]*shopify.Product, error) {

	// debug := godebug.Debug("output")
	//spit := spew.ConfigState{Indent: " ", DisableCapacities: true, DisablePointerAddresses: true}

	s := fmt.Sprintf("== Exporting collection %d", collectionId)
	fmt.Println(s)
	ctx := context.Background()

	productFields := []string{"id", "handle"}
	if viper.GetBool("export.include-product-info") {
		productFields = append(productFields, "body_html", "title")
		// debug("Export fields: %s", productFields)
	}

	opt := &shopify.ProductListOptions{
		Fields:       productFields,
		CollectionId: collectionId,
	}

	products, err := client.Products.AutoPagingList(ctx, opt)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Products.List() returned error: %v", err)
		return nil, err
	}
	if len(products) == 0 {
		fmt.Fprintf(os.Stderr, "Products.List() returned no events")
	}
	return products, nil
}

func GetMetafieldsByProduct(productId int, namespace string, client *shopify.Client) ([]*shopify.Metafield, error) {
	ctx := context.Background()
	opt := &shopify.MetafieldListOptions{Namespace: namespace, Fields: []string{"id", "key", "value"}}
	metafields, _, err := client.Metafields.ListByProduct(ctx, productId, opt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Metafields.List() returned error: %v", err)
		return nil, err
	}
	if len(metafields) == 0 {
		fmt.Fprintf(os.Stderr, "Metafields.List() returned no events")
	}
	return metafields, nil
}

func GenerateProductDataOutput(fields []*shopify.Metafield) (data []*OutputField) {

	for _, field := range fields {
		out := &OutputField{
			Data: make(map[string]map[string]string),
			Key:  field.Key,
			Id:   field.Id,
		}
		rows := strings.Split(*field.Value, rowSeparator)
		for i, row := range rows {
			cols := strings.Split(row, colSeparator)
			out.Data[strconv.Itoa(i)] = make(map[string]string)
			for j, col := range cols {
				out.Data[strconv.Itoa(i)][strconv.Itoa(j)] = col
			}
		}
		data = append(data, out)
	}
	return
}

func getSeoTagsByProduct(productID int, client *shopify.Client) (globalTitleTag *string, globalDescriptionTag *string) {
	globalMetafields, _ := GetMetafieldsByProduct(productID, "global", client)
	for _, field := range globalMetafields {
		if *field.Key == "title_tag" {
			globalTitleTag = field.Value
		}
		if *field.Key == "description_tag" {
			globalDescriptionTag = field.Value
		}
	}
	return globalTitleTag, globalDescriptionTag
}
