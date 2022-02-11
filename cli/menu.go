package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"tel-note/internal/config"
	"tel-note/internal/service/basic_info"
	"tel-note/internal/service/contact"
	"tel-note/internal/service/identity"
	"tel-note/internal/storage"
)

func ShowMenu(MainData *storage.AllContact, MainCity *storage.AllCities) {
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("*** Main menu: ***\n", "Please select:") //todo: Why fmt.Println don't split contents line
	fmt.Println("-------------------------------------------------------------")
	fmt.Print("'n|N' 	 	|	new record\n")
	fmt.Print("'l|L' 	 	|	list of contact\n")
	fmt.Print("'f|F' 	 	|	find one contact by id\n")
	fmt.Print("'fc|FC'		|	find contact, contain some character\n")
	fmt.Print("'e|E' 		|	find and edit contact by contact id\n")
	fmt.Print("'d|D' 		|	delete contact by id\n")
	fmt.Print("'dm|DM'		|	delete multi contact by id(s)\n")
	fmt.Print("'da|DA'		|	delete all contacts\n")
	fmt.Println("-------------------------------------------------------------")
	if identity.IsAdmin {
		fmt.Print("'data|DATA'	|	insert some sample's contacts\n")
		fmt.Print("'nc|NC'		|	insert new city\n")
		fmt.Print("'lc|LC'		|	list of cities\n")
		fmt.Print("'ec|EC'		|	edit city by id\n")
		fmt.Println("-------------------------------------------------------------")
	}
	runMenu(MainData, MainCity)
}

func runMenu(MainData *storage.AllContact, MainCity *storage.AllCities) {
	for {
		var userInput string
		fmt.Scanln(&userInput)
		//fmt.Println(userInput) todo: when negative value inserted, maybe because of define uint for id; The positive value is entered again, why?
		//admin actions only:
		if identity.IsAdmin {
			switch userInput {
			case "data", "DATA":
				contact.FillSimpleDataInMainData(MainData)
				ShowMenu(MainData, MainCity)
			case "nc", "NC":
				var inputCity string
				fmt.Println("insert city name:")
				fmt.Scanln(&inputCity)
				basic_info.NewCity(MainCity, inputCity)
				ShowMenu(MainData, MainCity)
			case "lc", "LC":
				dataJSON, _ := json.MarshalIndent(MainCity, "", "  ")
				fmt.Println(string(dataJSON))
				fmt.Println("-------------------------------------------------------------")
				for _, data := range MainCity.CityData {
					fmt.Printf("%3v | %-15s \n", data.Id, data.Name)
				}
				ShowMenu(MainData, MainCity)
			case "ec", "EC":
				var inputID uint
				var inputName string
				fmt.Println("insert city id:")
				fmt.Scanln(&inputID)
				fmt.Println("insert new city name:")
				fmt.Scanln(&inputName)
				if basic_info.EditCityByID(MainCity, inputID, inputName).State {
					fmt.Println("City changed ...")
				}
				ShowMenu(MainData, MainCity)
			}
		}
		switch userInput {
		case "N", "n":
			//todo: “To use function, or to use method, that is the question”
			var firstName, lastName, tel, cellphone, description string
			fmt.Println("Please inter your first name:")
			fmt.Scanln(&firstName)
			//input lastname
			fmt.Println("Please inter your last name:")
			fmt.Scanln(&lastName)
			//input tel
			fmt.Println("Please inter your tel:")
			fmt.Scanln(&tel)
			//input cellphone
			fmt.Println("Please inter your cellphone:")
			fmt.Scanln(&cellphone)
			//input description lastID +
			fmt.Println("Please inter your description:")
			fmt.Scanln(&description)
			contact.NewContact(MainData, firstName, lastName, tel, cellphone, description)
			//(*storage.AllContact).AddContact(MainData)
			fmt.Println(">> New record done <<")
			ShowMenu(MainData, MainCity)
		case "L", "l":
			dataJSON, _ := json.MarshalIndent(MainData, "", "  ")
			fmt.Println(string(dataJSON))
			fmt.Println("-------------------------------------------------------------")
			for _, data := range MainData.ContactData {
				fmt.Printf("%3v | %-15s | %-35v | %-5v | %-15v | %-5v\n", data.Id, data.FirstName, data.LastName, data.Tel, data.Cellphone, data.Description)
			}
			ShowMenu(MainData, MainCity)
		case "F", "f":
			var id uint
			fmt.Println("Please insert your contact ID:")
			fmt.Scanln(&id)
			result, isFound := contact.FindContactByID(MainData, id)
			if isFound {
				fmt.Println(result)
			} else {
				fmt.Println("not found")
			}
			ShowMenu(MainData, MainCity)
		case "FC", "fc":
			var insertChar string
			fmt.Println("insert character(s):")
			fmt.Scanln(&insertChar)
			resultData, resultCount := contact.FindContactByChar(MainData, insertChar)
			fmt.Println(resultData)
			if resultCount == 0 {
				fmt.Println("not found")
			} else {
				fmt.Println(resultCount, "record(s) found")
			}
			ShowMenu(MainData, MainCity)
		case "E", "e":
			var (
				insertContactID                                  uint
				firstName, lastName, tel, cellphone, description string
			)
			fmt.Println("Please insert your contact id you want to edit:")
			fmt.Scanln(&insertContactID)
			result, isFound := contact.FindContactByID(MainData, insertContactID)
			if isFound {
				fmt.Println(result)
			} else {
				fmt.Println("not found")
				ShowMenu(MainData, MainCity)
			}
			fmt.Println("new first name:")
			fmt.Scanln(&firstName)
			fmt.Println("new last name:")
			fmt.Scanln(&lastName)
			fmt.Println("new description:")
			fmt.Scanln(&tel)
			fmt.Println("new tel:")
			fmt.Scanln(&cellphone)
			fmt.Println("new cellphone:")
			fmt.Scanln(&description)
			editedContact := storage.Contact{Person: &storage.Person{FirstName: firstName, LastName: lastName}, Tel: tel, Cellphone: cellphone, Description: description}
			if state := contact.EditContactByID(MainData, editedContact, insertContactID); state.State != false {
				fmt.Println("contact edited successful")
			} else {
				println("contact not changed")
			}
			ShowMenu(MainData, MainCity)
		case "D", "d":
			var confirmDel string
			fmt.Println("*** important, be careful, you are deleting a contact ***")
			fmt.Println("do you want to continue? (yes or no)")
			fmt.Scanln(&confirmDel)
			if strings.ToLower(confirmDel) == config.OkStatus {
				var deleteID uint
				fmt.Println("insert your contact id that you want to delete:")
				fmt.Scanln(&deleteID)
				status := contact.DeleteContactByID(MainData, deleteID)
				switch status.State {
				case true:
					fmt.Printf("contact with id number:  %d deleted.", deleteID)
				case false:
					fmt.Println("something wrong")
				}
			}
			ShowMenu(MainData, MainCity)
		case "DA", "da":
			var confirmDel string
			fmt.Println("*** important, be careful, you are deleting all of contacts ***")
			fmt.Println("are you sure? (yes or no)")
			fmt.Scanln(&confirmDel)
			if strings.ToLower(confirmDel) == config.OkStatus {
				resultStatus := contact.DeleteAll(MainData)
				switch resultStatus.State {
				case true:
					fmt.Println("All contact deleted")
				case false:
					fmt.Println("something wrong")
				}
			}
		case "DM", "dm":
			var confirmDel string
			fmt.Println("*** important, be careful, you are deleting contact(s) ***")
			fmt.Println("do you want to continue? (yes or no)")
			fmt.Scanln(&confirmDel)
			if strings.ToLower(confirmDel) == config.OkStatus {
				var deleteIDS string
				var status config.ResponseStatus
				fmt.Println("insert your contact id(s) that you want to delete, separate id's by ',':")
				fmt.Scanln(&deleteIDS)
				idPack := strings.Split(deleteIDS, ",")
				for _, delID := range idPack {
					uintDelID, _ := strconv.ParseUint(delID, 10, 8)
					if uint(uintDelID) <= 0 {
						fmt.Println("not equal kind")
						break
					} else {
						status = contact.DeleteContactByID(MainData, uint(uintDelID))
					}
				}
				switch status.State {
				case true:
					fmt.Printf(deleteIDS)
					fmt.Println("contact(s) deleted.")
				case false:
					fmt.Println("something wrong")
				}
			}
			ShowMenu(MainData, MainCity)
		//something wrong:
		default:
			fmt.Println("bad input, please insert character of menu list")
		}
	}
}
