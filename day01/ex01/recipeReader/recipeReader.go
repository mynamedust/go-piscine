package recipeReader

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ItemType struct {
	Name  string `xml:"itemname" json:"ingredient_name"`
	Count string `xml:"itemcount" json:"ingredient_count"`
	Unit  string `xml:"itemunit,omitempty" json:"ingredient_unit,omitempty"`
}

type CakeType struct {
	Name        string     `xml:"name" json:"name"`
	StoveTime   string     `xml:"stovetime" json:"time"`
	Ingredients []ItemType `xml:"ingredients>item" json:"ingredients"`
}

type RecipesType struct {
	XMLName xml.Name   `xml:"recipes" json:"-"`
	Cakes   []CakeType `xml:"cake" json:"cake"`
}

type DBReader interface {
	Read() (*RecipesType, error)
}

type JSONReader struct {
	File string
}

type XMLReader struct {
	File string
}

func main() {
	fileName := flag.String("f", "", "Write valid .xml/.json filename.")
	flag.Parse()
	if ext := filepath.Ext(*fileName); ext == ".xml" {
		data := new(XMLReader)
		data.File = *fileName
		xml, _ := data.Read()
		PrintJson(xml)
	} else if ext == ".json" {
		data := new(JSONReader)
		data.File = *fileName
		json, _ := data.Read()
		PrintXml(json)
	} else {
		fmt.Println("Input filename not valid. Write valid .xml/.json filename.")
		return
	}
}

func (j JSONReader) Read() (recipe *RecipesType, err error) {
	file, err := os.Open(j.File)
	if err != nil {
		fmt.Println("File open error: ", err)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("File read error: ", err)
		return
	}
	recipe = new(RecipesType)
	if err := json.Unmarshal(data, recipe); err != nil {
		fmt.Println("Unmarshaling error: ", err)
		return nil, err
	}
	return
}

func (x XMLReader) Read() (recipe *RecipesType, err error) {
	file, err := os.Open(x.File)
	if err != nil {
		fmt.Println("File read error: ", err)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("File read error: ", err)
		return
	}
	recipe = new(RecipesType)
	if err = xml.Unmarshal(data, recipe); err != nil {
		fmt.Println("Unmarshaling error: ", err)
		return nil, err
	}
	return
}

func PrintJson(recipe *RecipesType) {
	data, err := json.MarshalIndent(recipe, "", "    ")
	if err != nil {
		fmt.Println("File marshaling error: ", err)
	}
	fmt.Println(string(data))
}

func PrintXml(recipe *RecipesType) {
	data, err := xml.MarshalIndent(recipe, "", "    ")
	if err != nil {
		fmt.Println("File marshaling error: ", err)
	}
	fmt.Println(string(data))
}
