package main

import (
	s "BetterString"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func saveExists(filename string) bool {
	_, err := os.Stat("saves/" + filename + ".toml")
	return err == nil
}

func copyF(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		err := srcFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() {
		err := dstFile.Close()
		if err != nil {

		}
	}()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func validatePath(baseDir, path string) (string, error) {
	cleanedPath := filepath.Clean(path)
	fullPath := filepath.Join(baseDir, cleanedPath)
	if !strings.HasPrefix(fullPath, filepath.Clean(baseDir)+string(filepath.Separator)) {
		return "", errors.New("odwołanie do katalogu spoza saves/ jest zabronione")
	}

	return fullPath, nil
}

func SaveMenager() (string, bool) {
	PrintColor("<bold><x>Menażer save'ów</x><bold>")
	PrintColor("To miejsce, gdzie możesz przejrzeć swoje zapisane save'y")
	PrintColor("I je usuwać/kopiować lub zmieniać im nazwy")
	Println()
	help := func() {
		PrintS("Komendy:")
		PrintColor("<bold>list</bold> - wyświetla liste save'ów")
		PrintColor("<bold>copy/cp {nazwa save} {nazwa kopii}</bold> - kopiuje save")
		PrintColor("<bold>rename/rn {nazwa save} {nowa nazwa}</bold> - zmienia nazwa save")
		PrintColor("<bold>delete/del {nazwa save}</bold> - usuwa save")
		PrintColor("<bold>help</bold> - wyświetla pomoc")
		PrintColor("<bold>play</bold> - rozpoczyna gre na wybranym save")
		PrintColor("<bold><x>exit</x></bold> - opuszcza save menager")
	}
	help()
	for {
		inp := s.New(Prompt(">>> ")).TrimSpace()
		switch {
		case inp.ToLower().String() == "list":
			listSaves()
		case inp.ToLower().HasPrefix("copy") || inp.ToLower().HasPrefix("cp"):
			if inp.HasPrefix("copy") {
				inp = inp.TrimPrefix("copy")
			} else if inp.HasPrefix("cp") {
				inp = inp.TrimPrefix("cp")
			}
			parts := inp.Split(" ")

			if len(parts) < 2 {
				ShowError("Musisz podać źródło i cel kopiowania save'a!")
				continue
			}

			if !saveExists(parts[0].String()) {
				ShowErrorf("Podany save (%s) nie istnieje!", parts[0].String())
				continue
			}

			err := copyF("saves/"+parts[0].Add(".toml").String(), "saves/"+parts[1].Add(".toml").String())
			if err != nil {
				ShowError("Niestety, ale wystąpił problem podczas kopiowania save'a")
				debug(err.Error())
			}
			PrintS(Sprintf("Pomyślnie skopiowano save (%s) jako %s", parts[0].String(), parts[1].String()))

		case inp.ToLower().HasPrefix("rename") || inp.ToLower().HasPrefix("rn"):
			if inp.HasPrefix("rename") {
				inp = inp.TrimPrefix("rename")
			} else if inp.HasPrefix("rn") {
				inp = inp.TrimPrefix("rn")
			}

			parts := inp.TrimSpace().Split(" ")

			debug(parts)
			debug(len(parts))
			if len(parts) < 2 {
				ShowErrorf("Podano za mało argumentów! Wymagane: nazwa obecna i nazwa nowa.")
				continue
			}
			if len(parts) > 2 {
				ShowErrorf("Podano zbyt wiele argumentów! Wymagane tylko dwa: nazwa obecna i nazwa nowa.")
				continue
			}

			oldName := parts[0].TrimSpace()
			newName := parts[1].TrimSpace()

			if oldName == "" || newName == "" {
				ShowErrorf("Nazwy plików nie mogą być puste!")
				continue
			}
			if !saveExists(oldName.String()) {
				ShowErrorf("Podany save (%s) nie istnieje!", oldName.String())
				continue
			}

			err := os.Rename("saves/"+oldName.Add(".toml").String(), "saves/"+newName.Add(".toml").String())
			if err != nil {
				ShowError("Nie można zmienić nazwy save'a")
				debug(err.Error())
			} else {
				PrintS("Pomyślnie zmieniono nazwę save")
			}
		case inp.ToLower().HasPrefix("delete") || inp.ToLower().HasPrefix("del"):
			if inp.HasPrefix("delete") {
				inp = inp.TrimPrefix("delete").TrimSpace()
			} else if inp.HasPrefix("del") {
				inp = inp.TrimPrefix("del").TrimSpace()
			}
			parts, err := inp.SplitString()
			if err != nil {
				ShowError(err.Error())
				continue
			}
			var stack []interface{}
			for _, sv := range parts {
				if !saveExists(sv) {
					stack = append(stack, Sprintf("Podany save (%s) nie istnieje!", sv))
					continue
				}

				safePath, err := validatePath("saves", sv+".toml")
				if err != nil {
					stack = append(stack, err)
					continue
				}

				err = os.Remove(safePath)
				if err != nil {
					stack = append(stack, err)
				} else {
					stack = append(stack, true)
				}
			}
			handleStack(stack, parts)
		case inp.HasPrefix("play"):
			game := inp.TrimPrefix("play").TrimSpace()
			if saveExists(game.String()) {
				return game.String(), true
			} else {
				ShowErrorf("Save (%s) nie istnieje!", game.String())
				continue
			}
		case inp == "help":
			help()
		case inp == "exit":
			return "skibidi", false
		default:
			ShowError("Nieznana opcja!")
		}
	}
}

func handleStack(stack []interface{}, parts []string) {
	debug(stack)
	allTrue := true
	var errorMessages []string

	for _, val := range stack {
		switch v := val.(type) {
		case bool:
			if !v {
				allTrue = false
			}
		case string:
			allTrue = false
			errorMessages = append(errorMessages, v)
		}
	}

	if allTrue {
		PrintS("Wszytskie operacje zakończone sukcesem!")
	} else {
		PrintColor("<orange><bold>Podczas usuwania wystąpiły błędy!</bold></orange>")
		for i, err := range errorMessages {
			idk := parts[i]
			ShowErrorf("%s: %s", idk, err)
		}
		PrintS("Reszta operacji przebiegła pomyślnie!")
	}
}
