package main

import (
	"fmt"
)

type ItemType struct {
	Name string `xml:"itemname" json:"ingredient_name"`
	Count string `xml:"itemcount" json:"ingredient_count"`
	Unit string `xml:"itemunit,omitempty" json:"ingredient_unit,omitempty"`
}

type CakeType struct {
	Name string `xml:"name" json:"name"`
	StoveTime string `xml:"stovetime" json:"time"`
	Ingredients []ItemType `xml:"ingredients>item" json:"ingredients"`
}

type RecipesType struct {
	Cake []CakeType `xml:"cake" json:"cake"`
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

}

func (r JSONReader) Read() (recipe *RecipesType, err error) {
}

func (r XMLReader) Read() (recipe *RecipesType, err error) {
}