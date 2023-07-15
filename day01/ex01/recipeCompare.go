package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"recipeReader"
)

func main() {
	oldName := flag.String("old", "", "Write valid .xml/.json filename.")
	newName := flag.String("new", "", "Write valid .xml/.json filename.")
	flag.Parse()
	oldData, newData := dataCreate(*oldName), dataCreate(*newName)
	if oldData == nil || newData == nil {
		fmt.Println("Write valid .xml/.json filename.")
		return
	}
	recipesCompare(oldData, newData)
}

func dataCreate(name string) (recipe *recipeReader.RecipesType) {
	var err error
	ext := filepath.Ext(name)
	if ext == ".xml" {
		data := new(recipeReader.XMLReader)
		data.File = name
		recipe, err = data.Read()
	} else if ext == ".json" {
		data := new(recipeReader.JSONReader)
		data.File = name
		recipe, err = data.Read()
	}
	if err != nil {
		return nil
	}
	return recipe
}

func recipesCompare(old, new *recipeReader.RecipesType) {
	for _, elem := range old.Cakes {
		if newCake := findCake(elem.Name, new.Cakes); newCake == nil {
			fmt.Printf("REMOVED cake \"%s\"\n", elem.Name)
		}
	}
	for _, elem := range new.Cakes {
		oldCake := findCake(elem.Name, old.Cakes)
		if oldCake == nil {
			fmt.Printf("ADDED cake \"%s\"\n", elem.Name)
			continue
		}
		checkTime(elem.Name, elem.StoveTime, oldCake.StoveTime)
		checkIngredients(elem.Name, elem, *oldCake)

	}
}

func findCake(cakeName string, cakesArr []recipeReader.CakeType) (newCake *recipeReader.CakeType) {
	for _, elem := range cakesArr {
		if elem.Name == cakeName {
			return &elem
		}
	}
	return nil
}

func checkTime(cakeName, time1, time2 string) {
	if time1 != time2 {
		fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", cakeName, time1, time2)
	}
}

func checkIngredients(cakeName string, cake1, cake2 recipeReader.CakeType) {
	for _, elem := range cake1.Ingredients {
		if newItem := findIngredient(elem, cake2.Ingredients); newItem == nil {
			fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", elem.Name, cakeName)
		}
	}
	for _, elem := range cake2.Ingredients {
		newItem := findIngredient(elem, cake1.Ingredients)
		if newItem == nil {
			fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", elem.Name, cakeName)
			continue
		}
		checkIngredientProps(cakeName, elem, *newItem)
	}
}

func findIngredient(item recipeReader.ItemType, itemArr []recipeReader.ItemType) *recipeReader.ItemType {
	for _, elem := range itemArr {
		if elem.Name == item.Name {
			return &elem
		}
	}
	return nil
}

func checkIngredientProps(cakeName string, item1, item2 recipeReader.ItemType) {
	if item1.Count != item2.Count {
		fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n",
			item1.Name, cakeName, item1.Count, item2.Count)
	}
	if item1.Unit == "" && item2.Unit != "" {
		fmt.Printf("ADDED unit \"%s\" for ingredient \"%s\" for cake  \"%s\"\n",
			item2.Unit, item1.Name, cakeName)
	} else if item1.Unit != "" && item2.Unit == "" {
		fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake  \"%s\"\n",
			item1.Unit, item1.Name, cakeName)
	} else if item1.Unit != item2.Unit {
		fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n",
			item1.Name, cakeName, item1.Unit, item2.Unit)
	}
}
