package main

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type Catalog struct {
	Shop struct {
		Name   string `xml:"name"`
		Offers struct {
			Offer []Offer `xml:"offer"`
		} `xml:"offers"`
	} `xml:"shop"`
}

type Offer struct {
	Available   bool     `xml:"available,attr"`
	GroupID     int      `xml:"group_id,attr"`
	ID          int      `xml:"id,attr"`
	URL         string   `xml:"url"`
	Price       int      `xml:"price"`
	OldPrice    int      `xml:"old_price"`
	Currency    string   `xml:"currencyId"`
	Pictures    []string `xml:"picture"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Vendor      string   `xml:"vendor"`
	Sku         string   `xml:"vendorCode"`
	CategoryID  int      `xml:"categoryId"`
	Params      []Param  `xml:"param"`
}

type Param struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type csvDataForProduct struct {
	Name   string
	Season string
}

type Product struct {
	Sku         string
	Name        string
	Seasons     []string
	Available   bool
	GroupID     int
	ID          int
	URL         string
	Price       int
	OldPrice    int
	Currency    string
	Pictures    []string
	OldXmlName  string
	Description string
	Vendor      string
	CategoryID  int
	Styles      []string
	Kind        string
	Sizes       []string
	Color       string
	Collections []string
	Composition string
	Materials   []string
	Dimensions  string
}

const (
	StyleParam       = "Стиль"
	KindParam        = "Вид"
	SizeParam        = "Размер"
	ColorParam       = "Цвет"
	CollectionParam  = "Коллекция"
	CompositionParam = "Состав"
	MaterialParam    = "Материал"
	DimensionsParam  = "Замеры"
)

func main() {
	var csvFilePath string
	var xmlFilePath string
	flag.StringVar(&csvFilePath, "csv", "pkg/data_base/Files/2024.02.13.csv", "CSV file path")
	flag.StringVar(&xmlFilePath, "xml", "pkg/data_base/Files/export_rozetka.xml", "XML file path")
	flag.Parse()
	if csvFilePath == "" || xmlFilePath == "" {
		fmt.Println("Usage: program_name -csv <csv_file_path> -xml <xml_file_path>")
		return
	}

	csvBytes, err := OpenFile(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	xmlBytes, err := OpenFile(xmlFilePath)
	if err != nil {
		log.Fatal(err)
	}

	csvFileData, err := Parsecsv(csvBytes)
	if err != nil {
		log.Fatal(err)
	}

	xmlFileData, err := Parsexml(xmlBytes)
	if err != nil {
		log.Fatal(err)
	}

	skirtsCategoryId := 28
	skirts := MapCategory(csvFileData, xmlFileData, skirtsCategoryId)
	for _, skirt := range skirts {
		fmt.Println(skirt)
	}
}

func OpenFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Parsecsv(data []byte) (map[string]csvDataForProduct, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = ';'
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	mapData := make(map[string]csvDataForProduct)
	for _, row := range records {
		name := row[2]
		sku := row[4]
		season := row[9]
		mapData[sku] = csvDataForProduct{
			Name:   name,
			Season: season,
		}
	}
	return mapData, nil
}

func Parsexml(data []byte) (*Catalog, error) {
	var catalog Catalog
	err := xml.NewDecoder(bytes.NewReader(data)).Decode(&catalog)
	if err != nil {
		return nil, err
	}
	return &catalog, nil
}

func MapCategory(csvData map[string]csvDataForProduct, catalog *Catalog, productCategoryId int) []*Product {
	skirts := make([]*Product, 0)
	mSkirts := make(map[string]*Product)
	for _, offer := range catalog.Shop.Offers.Offer {
		if offer.CategoryID == productCategoryId {
			if v, ok := mSkirts[offer.Sku]; ok {
				v.Seasons = append(v.Seasons, csvData[offer.Sku].Season)
				v.Styles = append(v.Styles, GetParam(offer.Params, StyleParam))
				v.Sizes = append(v.Sizes, GetParam(offer.Params, SizeParam))
				v.Collections = append(v.Collections, GetParam(offer.Params, CollectionParam))
				v.Materials = append(v.Materials, GetParam(offer.Params, MaterialParam))
				continue
			}
			product := &Product{
				Sku:         offer.Sku,
				Name:        csvData[offer.Sku].Name,
				Seasons:     []string{csvData[offer.Sku].Season},
				Available:   offer.Available,
				GroupID:     offer.GroupID,
				ID:          offer.ID,
				URL:         offer.URL,
				Price:       offer.Price,
				OldPrice:    offer.OldPrice,
				Currency:    offer.Currency,
				Pictures:    offer.Pictures,
				OldXmlName:  offer.Name,
				Description: offer.Description,
				Vendor:      offer.Vendor,
				CategoryID:  offer.CategoryID,
				Styles:      []string{GetParam(offer.Params, StyleParam)},
				Kind:        GetParam(offer.Params, KindParam),
				Sizes:       []string{GetParam(offer.Params, SizeParam)},
				Color:       GetParam(offer.Params, ColorParam),
				Collections: []string{GetParam(offer.Params, CollectionParam)},
				Composition: GetParam(offer.Params, CompositionParam),
				Materials:   []string{GetParam(offer.Params, MaterialParam)},
				Dimensions:  GetParam(offer.Params, DimensionsParam),
			}
			skirts = append(skirts, product)
			mSkirts[offer.Sku] = product
		}
	}
	return skirts
}

func GetParam(params []Param, key string) string {
	for _, param := range params {
		if param.Name == key {
			return param.Value
		}
	}
	return ""
}
