package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	// "log"
	"context"
	client "github.com/zinclabs/sdk-go-zincsearch"
)

func main() {
	var cont int
	var registers string
	cont = 0
	registers = ""
	searchEmails("./enron_mail_20110402", &cont, &registers)
}

type Email struct {
	MessageId string
	Date      string
	From      string
	To        string
	Subject   string
	Xfolder   string
	Message   string
}

func searchEmails(dir string, cont *int, registers *string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading email dir ", err)
	}
	for _, file := range files {
		if file.IsDir() {
			searchEmails(dir+"/"+file.Name(), cont, registers)
		} else {
			scanText(dir+"/"+file.Name(), cont, registers)
			*cont += 1
			fmt.Println(*cont)
		}
	}
}

func scanText(dir string, cont *int, registers *string) {
	emptyLine := 0
	message := ""
	
	var structMail Email
	file, err := os.Open(dir)
	if err != nil {
		fmt.Println("Error reading email information ", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 && emptyLine == 0 {
			emptyLine += 1
		}

		if strings.HasPrefix(scanner.Text(), "Message-ID: ") {
			structMail.MessageId = strings.Split(scanner.Text(), ": ")[1]
		} else if strings.HasPrefix(scanner.Text(), "Date: ") {
			structMail.Date = strings.Split(scanner.Text(), ": ")[1]
		} else if strings.HasPrefix(scanner.Text(), "From: ") {
			structMail.From = strings.Split(scanner.Text(), ": ")[1]
		} else if strings.HasPrefix(scanner.Text(), "To: ") {
			structMail.To = strings.Split(scanner.Text(), ": ")[1]
		} else if strings.HasPrefix(scanner.Text(), "Subject: ") {
			structMail.Subject = strings.Split(scanner.Text(), ": ")[1]
		} else if strings.HasPrefix(scanner.Text(), "X-Folder: ") {
			structMail.Xfolder = strings.Split(scanner.Text(), ": ")[1]
		}

		if emptyLine == 1 {
			message += "\n" + scanner.Text()
		}
	}
	structMail.Message = message

	e, err := json.MarshalIndent(structMail, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}

	index := `{ "index" : { "_index" : "Emails" } }`
	*registers += index + "\n" +string(e)+ "\n"

	if *cont == 500 {

		// f, err := os.OpenFile("processedMails.ndjson", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// f.WriteString(*registers)
		// f.Close()
		// panic(1)

		// fmt.Println("500 REGISTER")
		// fmt.Println(*registers)
		user := "admin"
		pass := "123"

		auth := context.WithValue(context.Background(), client.ContextBasicAuth, client.BasicAuth{
			UserName: user,
			Password: pass,
		})
		configuration := client.NewConfiguration()
		apiClient := client.NewAPIClient(configuration)
		resp, r, err := apiClient.Document.Bulk(auth).Query(*registers).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `Document.Bulk``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		}
		// response from `Bulk`: MetaHTTPResponseRecordCount
		fmt.Fprintf(os.Stdout, "Response from `Document.Bulk`: %v\n", *resp)

		*cont = 0
		*registers = ""
	}
	

	// f, err := os.OpenFile("processedMails.ndjson", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// f.WriteString(`{ "index" : { "_index" : "Emails" } }`)
	// f.WriteString("\n")
	// f.WriteString(string(e) + "\n")
	// f.Close()
}
