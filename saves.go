package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func loadSave(name string) (int, error) {
	type Player struct {
		Hungry int    `toml:"hungry"`
		Place  string `toml:"place"`
		Town   string `toml:"town"`
		Name   string `toml:"name"`
	}

	type Money struct {
		Wallet float64 `toml:"wallet"`
		Bank   float64 `toml:"bank"`
	}

	type SaveData struct {
		Player Player `toml:"player"`
		Money  Money  `toml:"money"`
	}

	defaultSave := SaveData{
		Player: Player{
			Hungry: 100,
			Place:  "DOM",
			Town:   "_",
			Name:   "_",
		},
		Money: Money{
			Wallet: 0,
			Bank:   1000,
		},
	}

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
		hungry = defaultSave.Player.Hungry
		PLACE = defaultSave.Player.Place
		TOWN = defaultSave.Player.Town
		wallet = defaultSave.Money.Wallet
		bank = defaultSave.Money.Bank
		NAME = defaultSave.Player.Name
		/* DEFAULT */

		// Return exit code 1 (file did not exist and was created with default values)
		return 1, nil
	}

	// Validate loaded data
	if save.Player.Hungry < 1 || save.Player.Hungry > 100 || save.Player.Place == "" || save.Player.Town == "" ||
		save.Player.Name == "" || save.Money.Wallet < 0 || save.Money.Bank < 0 {
		// Invalid structure or missing values, reset to default
		file, err := os.Create(filePath)
		if err != nil {
			return 2, err // Error while resetting the save file
		}
		defer func() {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}()

		err = toml.NewEncoder(file).Encode(defaultSave)
		if err != nil {
			return 2, err // Error while writing default save
		}

		/* DEFAULT */
		hungry = defaultSave.Player.Hungry
		PLACE = defaultSave.Player.Place
		TOWN = defaultSave.Player.Town
		wallet = defaultSave.Money.Wallet
		bank = defaultSave.Money.Bank
		NAME = defaultSave.Player.Name
		/* DEFAULT */

		// Return exit code 2 (invalid structure, reset file to defaults)
		return 2, nil
	}

	hungry = save.Player.Hungry
	PLACE = save.Player.Place
	TOWN = save.Player.Town
	wallet = save.Money.Wallet
	bank = save.Money.Bank
	NAME = save.Player.Name

	return 0, nil
}

func saveSave(name string) error {
	// PATH
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

	// formating data
	content := fmt.Sprintf(`[player]
hungry = %d
place = "%s"
town = "%s"
name = "%s"

[money]
wallet = %.2f
bank = %.2f
`, hungry, PLACE, TOWN, NAME, wallet, bank)

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
