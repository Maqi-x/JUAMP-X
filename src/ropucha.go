package main

import (
	"math/rand"
	"strings"
	"time"
)

func handleRopucha() {
	Talk([][2]interface{}{
		{"Kasjer", "Witaj w najlepszym sklepie w całej dzielnicy!"},
		{"Kasjer", "Co podać?"},
		{"Ty", "Hm..."},
	}, map[string]string{
		"Kasjer": "magenta",
		"Ty":     "green",
	})

	list := func() {
		PrintLine("Dostępne produkty:")
		for name, info := range products {
			price := info.Price
			PrintColor(Sprintf("%.2f - <bold>%s</bold>: %s", price, name, func() string {
				if price < wallet {
					return colorCodes["green"] + "✔" + colorCodes["reset"]
				} else {
					return colorCodes["red"] + "❌" + colorCodes["reset"]
				}
			}()))
		}
		Sep()
	}

	help := func() {
		PrintS("Komendy:")
		PrintColor("<bold>info {nazwa przedmiotu}</bold> - wyświetla informacje o przedmiocie")
		PrintColor("<bold>list</bold> - wyświetla liste dostępnych przedmiotów, wraz z cenami")
		PrintColor("<bold>buy {nazwa przedmiotu}</bold> - kupuje przedmiot")
		PrintColor("<bold>help</bold> - wyświetla pomoc")
		PrintColor("<x><bold>exit</bold></x> - opuszcza sklep")
	}
	help()

	products = make(map[string]struct {
		Price       float64
		Description string
	})
	for key, v := range productS {
		products[normStr(key)] = v
	}
	// ---------------------------------------------- main loop ----------------------------------------- \\
	for {
		inp := Prompt(">>> ")

		if inp == "help" {
			help()
		} else if inp == "list" {
			list()
			continue
		} else if inp == "exit" {
			break
		} else if strings.HasPrefix(inp, "info") {
			if len(inp) < 6 {
				Talk([][2]interface{}{
					{"Kasjer", "Podaj nazwe przedmiotu po komendzie `info`, a chętnie powiem ci o nim więcej"},
				}, map[string]string{
					"Kasjer": "magenta",
				})
				Println("\r")
				continue
			}

			key := normStr(strings.TrimSpace(inp[5:]))
			v, exists := products[key]
			if !exists {
				Talk([][2]interface{}{
					{"Kasjer", "Nie ma takiego przedmiotu!"},
					{"Kasjer", "Uzyj `list` by wyświetlić dostępne przedmioty"},
				}, map[string]string{
					"Kasjer": "magenta",
				})
				Println("\r")
				continue
			}
			if key == "gazeta" {
				v.Description = Sprintf(v.Description, TOWN)
			}
			PrintColor(Sprintf("<x><bold>Info o przedmiocie:</x></bold> %s", key))
			PrintColor(Sprintf("<x><bold>Opis:</x></bold> %s", v.Description))
			PrintColor(Sprintf("<x><bold>Cena:</x></bold> %.2f", v.Price))
		} else if strings.HasPrefix(inp, "buy") {
			if len(inp) < 5 {
				Println("Podaj prawidłową nazwę przedmiotu po komendzie `buy`")
				continue
			}

			key := strings.TrimSpace(inp[4:])
			v, exists := products[normStr(key)]
			if !exists {
				Println("Nie ma takiego przedmiotu!")
				continue
			}

			ok := cPay(v.Price)
			if !ok {
				Println("Wygląda na to, że nie stać cię na ten produkt...")
				Println("Spróbuj ponownie później")
				if bank > v.Price {
					PrintS("Wskazówka:")
					PrintColor("W <bold>banku</bold> masz wystarczającą ilość pieniędzy, by zakupić ten produkt!")
					PrintColor("Możesz udać się do banku i <bold>wypłacić</bold> pieniądze, a następnie tu wrócić")
				}
			} else {
				if key == "gazeta" && assortment["newspapers"] > 0 {
					assortment["newspapers"]--
				} else if key == "bułka" && assortment["buns"] > 0 {
					assortment["buns"]--
				} else if key == "pizza" && assortment["pizzas"] > 0 {
					assortment["pizzas"]--
				} else if key == "buty" && assortment["shoes"] > 0 {
					assortment["shoes"]--
				} else {
					Talk([][2]interface{}{
						{"Kasjer", "Przykro mi, ale ten produkt już się wyprzedał. Musisz poczekać do dostaw"},
						{"Ty", "Niestety..."},
					}, map[string]string{
						"Kasjer": "magenta",
						"Ty":     "green",
					})
					Println()
					continue
				}
				buy(key)
			}
		} else {
			Println("Nieznana komenda.")
		}
	}
	// ----------------------------------------------------------------------------------------------- \\
	Back()
}

func buy(name string) {
	name = normStr(name)
	switch name {
	case "bułka":
		hungryAdd(20)
		loading(2, "Jedzienie")
		chm("ROPUCHA")
		time.Sleep(500 * ms)
	case "pizza":
		hungryAdd(40)
		loading(3, "Jedzienie")
		chm("ROPUCHA")
		time.Sleep(500 * ms)
	case "gazeta":
		Sep()
		Println("Gazeta:")
		for _, line := range strings.Split(papers[rand.Intn(len(papers))], "\n") {
			Println(line)
			time.Sleep(300 * ms)
		}
		Sep()
		if tutStep[0] == 1 {
			tutStep = []int{2, 0}
		}
	}
}
