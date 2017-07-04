This cli tool is a side kick of the [Power Editor](https://apps.shopify.com/power-editor) Shopify app.
It lets you export and import Power-Editor data (basically any metafields and textual product data) in a clean format.

## Restrictions
At the moment the tool only works with product related data (collections, pages and articles will follow)

## Use cases
* Sync product data between shops (power-editor data, product titles and descriptions)
* Bulk editing (e.g. translating into another language) of this data

## Usage
You need to set up a private app for the shopify store you want to use this tool with. See [this part of the Shopify docs](https://help.shopify.com/api/getting-started/authentication/private-authentication#generate-credentials-from-the-shopify-admin) for a short how to.
You need to provde the `powereditor_cli` with the key and password you obtained that way, alongside with your store name. You can do this via command line options but the easiest way is to use a `config.yml` file

`config.yml`

```yaml
export:
  key: iiiiiiiiii
  password: jjjjjjjjjj
  store: my-first-store
  namespace: power-editor

import:
  key: xxxxxxxxxx
  password: yyyyyyyyyyyy
  store: my-other-store
  namespace: test
```

put this file in an empty folder and `cd` into it. Then run the tool from the command line

### Export data

```
powereditor_cli collection export 12345678
```
Replace 12345678 above with the ID of the collection you'd like to export.
This will export the data to `output.json` within the same folder.

### Import data

```
powereditor_cli import output.json 
```

## More options

For more options see

```
powereditor-cli help export collection
powereditor-cli help import
```

## Example of exported data


```json
{
  "products": [
    {
      "id": 117563710,
      "handle": "multi-channelled-assymetric-capability",
      "body_html": "<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quid est igitur, cur ita semper deum appellet Epicurus beatum et aeternum? Conferam tecum, quam cuique verso rem subicias; Parvi enim primo ortu sic iacent, tamquam omnino sine animo sint. Apparet statim, quae sint officia, quae actiones. Quae cum dixisset paulumque institisset, Quid est?</p>",
      "title": "Clown1",
      "fields": [
        {
          "id": 32038808337,
          "key": "link1",
          "data": {
            "0": {
              "0": "fsdfsadf"
            }
          }
        },
        {
          "id": 32005501521,
          "key": "multimulit",
          "data": {
            "0": {
              "0": "naja",
              "1": "aha"
            }
          }
        },
        {
          "id": 32038808401,
          "key": "products",
          "data": {
            "0": {
              "0": "multi-channelled-assymetric-capability"
            },
            "1": {
              "0": "phased-explicit-architecture"
            },
            "2": {
              "0": "right-sized-clear-thinking-parallelism"
            }
          }
        },
        {
          "id": 32018829457,
          "key": "single",
          "data": {
            "0": {
              "0": "2nd"
            }
          }
        },
        {
          "id": 32038808273,
          "key": "single1",
          "data": {
            "0": {
              "0": "1st"
            }
          }
        },
        {
          "id": 31869769041,
          "key": "tabs",
          "data": {
            "0": {
              "0": "aha",
              "1": "AAA",
              "2": "<p>AAAAA</p>",
              "3": "AAAaaaa"
            },
            "1": {
              "0": "Mein Dingsd",
              "1": "false",
              "2": "falselll",
              "3": "<p>soso jajajaj</p>"
            }
          }
        },
        {
          "id": 32344699729,
          "key": "test",
          "data": {
            "0": {
              "0": "full"
            },
            "1": {
              "0": "left"
            }
          }
        }
      ]
    },
    {
      "id": 117563712,
      "handle": "enterprise-wide-upward-trending-hardware",
      "body_html": "<p>So this is a product.</p>\n<p>The text you see here is a Product Description. Every product has a price, a weight, a picture and a description. To edit the description of this product or to create a new product you can go to the <a href=\"/admin/products\">Products Tab</a> of the administration menu.</p>\n<p>Once you have mastered the creation and editing of products you will want your products to show up on your Shopify site. There is a two step process to do this.</p>\n<p>First you need to add your products to a Collection. A Collection is an easy way to group products together. If you go to the <a href=\"/admin/custom_collections\">Collections Tab</a> of the administration menu you can begin creating collections and adding products to them.</p>\n<p>Second you’ll need to create a link from your shop’s navigation menu to your Collections. You can do this by going to the <a href=\"/admin/links\">Navigations Tab</a> of the administration menu and clicking on “Add a link”.</p>\n<p>Good luck with your shop!</p>",
      "title": "Not a clown",
      "fields": [
        {
          "id": 33741946065,
          "key": "single",
          "data": {
            "0": {
              "0": "1st"
            }
          }
        }
      ]
    }
  ]
}
```