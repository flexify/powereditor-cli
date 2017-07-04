package cmd

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/dommmel/goshopping/shopify"
)

var spit spew.ConfigState

func init() {
	spit = spew.ConfigState{Indent: " ", DisableCapacities: true, DisablePointerAddresses: true}
}

type ProductInput struct {
	Id     *int                 `json:"id"`
	Handle *string              `json:"handle"`
	Fields []*shopify.Metafield `json:"fields"`
}

func TestJsonParsing(t *testing.T) {
	fields := []byte(`{
    "id": 3894060097,
    "handle": "blackroll-med-45",
    "fields": [
      {
        "id": 34506107137,
        "key": "accordion",
        "value": "GRÖSSE & GEWICHT<!--|col|--><ul>\n  <li>45 cm x 15 cm, 158 g&nbsp;</li>\n</ul><!--|row|-->LIEFERUMFANG<!--|col|--><ul>\n  <li>1 x BLACKROLL® MED 45 grün</li>\n  <li>1 x BLACKROLL® DVD</li>\n  <li>1 x BLACKROLL® Übungskarte&nbsp;</li>\n</ul><!--|row|-->MADE IN GERMANY<!--|col|--><ul>\n  <li>Höchste Produktqualität</li>\n  <li>Markenrechtlich geschütztes Produkt</li>\n  <li>Qualitätsmanagement nach DIN ISO 9001:2000&nbsp;</li>\n</ul><!--|row|-->PRODUKTIONSINFORMATIONEN<!--|col|--><ul>\n  <li>Umweltfreundliche und energieschonende Produktion</li>\n  <li>Material zu 100 % recyclefähig</li>\n  <li>Frei von Treibgasen</li>\n  <li>Frei von anderen chemischen Treibmitteln&nbsp;</li>\n</ul><!--|row|-->HYGIENE<!--|col|--><ul>\n  <li>Geruchlos</li>\n  <li>Wasserunlöslich</li>\n  <li>Einfach zu reinigen</li>\n  <li>Einfach zu sterilisieren&nbsp;</li>\n</ul>"
      },
      {
        "id": 34895481857,
        "key": "products",
        "value": "ball-1<!--|row|-->blackroll-duoball-12<!--|row|-->blackroll-mat<!--|row|-->blackroll-mini"
      },
      {
        "id": 34506107201,
        "key": "testimonial",
        "value": "<p>&nbsp;„Die Anwendungsmöglichkeiten der BLACKROLL® bieten eine ganze Reihe von Vorteilen. Die Option, die Druckstärke durch entsprechende Entlastungstechniken individuell zu variieren, erlaubt ein breites therapeutisches Spektrum. Ausführliche Tests mit Patienten zeigen hervorragende Ergebnisse.“ &nbsp;</p>\n<p><strong>Dr. biol. hum. Robert Schleip, Direktor Fascia Research Project, Universität Ulm</strong></p><!--|col|-->https://cdn.shopify.com/s/files/1/0429/1421/t/11/assets/Robert_Schleip-902301865.jpeg?8776001892895034102"
      },
      {
        "id": 34524200961,
        "key": "uebungen",
        "value": "Übung<!--|col|--><h3><strong>Ausführung</strong></h3>\n<ul>\n  <li>Oberkörper absenken, dabei Hände vom Boden nehmen und hinter dem Kopf verschränken.</li>\n  <li>Langsam auf der BLACKROLL® in Richtung Brustwirbelsäule rollen, indem die Beugung im Kniegelenk verstärkt wird.&nbsp;</li>\n</ul>\n<p><strong>WEITERE ÜBUNGEN FINDEST DU HIER.</strong> &nbsp;</p><!--|col|-->https://cdn.shopify.com/s/files/1/0429/1421/t/11/assets/BLACKROLL_MED4511-3581107480.jpeg?6250174402157654771<!--|col|-->https://www.blackroll.com/de/uebungen"
      },
      {
        "id": 34507710529,
        "key": "video",
        "value": "XGBQkxcM8DI"
      }
    ]
  }`)
	var pin ProductInput
	if err := json.Unmarshal(fields, &pin); err != nil {
		log.Fatal(err)
	}
	// for _, dataPoint := range p.Data {
	// 	rows := strings.Split(dataPoint["value"].(string), "<!--|row|-->")
	// 	for i, row := range rows {
	// 		cols := strings.Split(row, "<!--|col|-->")
	// 		for j, col := range cols {
	// 			log.Printf("Row %d, Col %d: %s", i, j, col)
	// 		}
	// 	}
	// }
	//

	// init output data structure for Product

	outfields := GenerateProductDataOutput(pin.Fields)
	_, err := JSONMarshalIndent(outfields, "", "  ")
	if err != nil {
		t.Errorf("Error generating output")
	}
}
