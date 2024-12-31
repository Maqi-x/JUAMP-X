package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

/****** STRUCTS ******/
type Player struct {
	Hungry int    `toml:"hungry"`
	Place  string `toml:"place"`
	Town   string `toml:"town"`
	Name   string `toml:"name"`
	Age    int    `toml:"age"`
}

type Game struct {
	Wallet       float64 `toml:"wallet"`
	Bank         float64 `toml:"bank"`
	Autosave     int     `toml:"autosave"`
	First        bool    `toml:"first"`
	TutorialStep []int   `toml:"tutorialStep"`
}

type Assortment struct {
	Buns       int `toml:"buns"`
	Pizzas     int `toml:"pizzas"`
	Newspapers int `toml:"newspapers"`
	Shoes      int `toml:"shoes"`
}

type SaveData struct {
	Player     Player     `toml:"player"`
	Game       Game       `toml:"game"`
	Assortment Assortment `toml:"assortment"`
}

/****** STRUCTS ******/

func loadSave(name string) (int, error) {
	filePath := "saves/" + name + ".toml"
	var save SaveData

	_, err := toml.DecodeFile(filePath, &save)
	if err != nil {
		err := os.MkdirAll("saves", 0755)
		if err != nil {
			return 1, err // Error while creating the saves directory
		}

		file, err := os.Create(filePath)
		if err != nil {
			return 1, err // Error while creating the save file
		}
		defer func() {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}()

		err = toml.NewEncoder(file).Encode(defaultSave)
		if err != nil {
			return 1, err // Error while writing default save
		}

		/* DEFAULT */
		initData()
		/* DEFAULT */

		// Return exit code 1 (file did not exist and was created with default values)
		return 1, nil
	}

	/****************** VALIDATE *****************8*/
	if save.Player.Hungry < 1 || save.Player.Hungry > 100 || save.Player.Place == "" || save.Player.Town == "" ||
		save.Player.Name == "" || save.Game.Wallet < 0 || save.Game.Bank < 0 || save.Game.Autosave < -1 ||
		len(save.Game.TutorialStep) != 2 ||
		save.Assortment.Buns < 0 || save.Assortment.Pizzas < 0 || save.Assortment.Newspapers < 0 || save.Assortment.Shoes < 0 || age < 5 {

		file, err := os.Create(filePath)
		if err != nil {
			return 2, err
		}
		defer func() {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}()
		/****************** VALIDATE *****************8*/

		err = toml.NewEncoder(file).Encode(defaultSave)
		if err != nil {
			return 2, err
		}

		initData()

		// Return exit code 2 (invalid structure, reset file to defaults)
		return 2, nil
	}

	hungry = save.Player.Hungry
	PLACE = save.Player.Place
	TOWN = save.Player.Town
	NAME = save.Player.Name
	wallet = save.Game.Wallet
	bank = save.Game.Bank
	autosave = save.Game.Autosave
	first = save.Game.First
	tutStep = save.Game.TutorialStep
	age = save.Player.Age
	assortment = map[string]int{
		"buns":       save.Assortment.Buns,
		"pizzas":     save.Assortment.Pizzas,
		"newspapers": save.Assortment.Newspapers,
		"shoes":      save.Assortment.Shoes,
	}

	return 0, nil
}

func saveSave(name string) error {
	if tmpses {
		return nil
	}
	/****** DATA ******/
	saveData := SaveData{
		Player: Player{
			Hungry: hungry,
			Place:  PLACE,
			Town:   TOWN,
			Name:   NAME,
			Age:    age,
		},
		Game: Game{
			Wallet:       wallet,
			Bank:         bank,
			Autosave:     autosave,
			First:        first,
			TutorialStep: tutStep,
		},
		Assortment: Assortment{
			Buns:       assortment["buns"],
			Pizzas:     assortment["pizzas"],
			Newspapers: assortment["newspapers"],
			Shoes:      assortment["shoes"],
		},
	}
	/****** DATA ******/

	filePath := "saves/" + name + ".toml"

	err := os.MkdirAll("saves", 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	encoder := toml.NewEncoder(file)
	err = encoder.Encode(saveData)
	if err != nil {
		return err
	}

	return nil
}

func initData() {
	/* DEFAULT */
	hungry = defaultSave.Player.Hungry
	PLACE = defaultSave.Player.Place
	TOWN = defaultSave.Player.Town
	NAME = defaultSave.Player.Name
	wallet = defaultSave.Game.Wallet
	bank = defaultSave.Game.Bank
	autosave = defaultSave.Game.Autosave
	first = defaultSave.Game.First
	tutStep = defaultSave.Game.TutorialStep
	age = defaultSave.Player.Age
	assortment = map[string]int{
		"buns":       defaultSave.Assortment.Buns,
		"pizzas":     defaultSave.Assortment.Pizzas,
		"newspapers": defaultSave.Assortment.Newspapers,
		"shoes":      defaultSave.Assortment.Shoes,
	}
	/* DEFAULT */
}
