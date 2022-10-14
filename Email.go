package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"context"
	client "github.com/zinclabs/sdk-go-zincsearch"
)

func main() {
	var cont int
	var registers []map[string]interface{}
	cont = 0
	searchEmails("./enron_mail_20110402", &cont, &registers)
	//countEmails("./enron_mail_20110402", &cont)
}



func countEmails(dir string, cont *int) {   
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading email dir ", err)
	}
	for _, file := range files {
		if file.IsDir() {
			countEmails(dir+"/"+file.Name(), cont)
		} else {

			*cont += 1
			fmt.Println(*cont)
		}
	}
}

func searchEmails(dir string, cont *int, registers *[]map[string]interface{}) {
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
		}
	}
}

func scanText(dir string, cont *int, registers *[]map[string]interface{}) {
	emptyLine := 0
	message := ""

	file, err := os.Open(dir)
	if err != nil {
		fmt.Println("Error reading email information ", err)
	}
	defer file.Close()
	record := map[string]interface{}{
		"MessageId": "",
		"Date":      "",
		"From":     "",
		"To":        "",
		"Subject":   "",
		"Xfolder":   "",
		"Message":   "",
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 && emptyLine == 0 {
			emptyLine += 1
		}

		if strings.HasPrefix(scanner.Text(), "Message-ID: ") {
			record["MessageId"] = strings.Split(scanner.Text(), ": ")[1]
		} else if strings.HasPrefix(scanner.Text(), "Date: ") {
			record["Date"] = strings.Split(scanner.Text(), ": ")[1]
		} else if strings.HasPrefix(scanner.Text(), "From: ") {
			record["From"] = strings.Split(scanner.Text(), ": ")[1]
		} else if strings.HasPrefix(scanner.Text(), "To: ") {
			record["To"] = strings.Split(scanner.Text(), ": ")[1]
		} else if strings.HasPrefix(scanner.Text(), "Subject: ") {
			record["Subject"] = strings.Split(scanner.Text(), ": ")[1]
		} else if strings.HasPrefix(scanner.Text(), "X-Folder: ") {
			record["Xfolder"] = strings.Split(scanner.Text(), ": ")[1]
		}

		if emptyLine == 1 {
			message += "\n" + scanner.Text()
		}
	}
	record["Message"] = message
	*registers = append(*registers, record)
	if err != nil {
		fmt.Println(err)
		return
	}

	index := `Emails`

	if *cont == 500 {

		user := "admin"
		pass := "123"

		auth := context.WithValue(context.Background(), client.ContextBasicAuth, client.BasicAuth{
			UserName: user,
			Password: pass,
		})
		query := *client.NewMetaJSONIngest()
		query.SetIndex(index)
		query.SetRecords(*registers)
		configuration := client.NewConfiguration()
		apiClient := client.NewAPIClient(configuration)
		resp, r, err := apiClient.Document.Bulkv2(auth).Query(query).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `Document.Bulk``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		}
		// response from `Bulk`: MetaHTTPResponseRecordCount
		fmt.Fprintf(os.Stdout, "Response from `Document.Bulk`: %v\n", *resp.RecordCount)

		*cont = 0
		*registers = nil
	}
}
