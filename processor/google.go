package processor

import (
	"encoding/json"
	"fmt"
	"gsheets-slack-chatbot/model/google"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func (p *Processor) getCellContent(message string) (string, error) {
	props, err := p.getCellProps(message)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return "", err
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return "", err
	}

	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		return "", err
	}

	resp, err := srv.Spreadsheets.Values.Get(props.SheetID, props.Range).Do()
	if err != nil {
		return "", err
	}

	if len(resp.Values) == 0 {
		return "", fmt.Errorf("no data found")
	}

	for _, row := range resp.Values {
		return fmt.Sprintf("%s", row[0]), nil
	}

	return "", fmt.Errorf("unhandled error")
}

func (p *Processor) getCellProps(message string) (model.CellProps, error) {
	pattern, err := p.config.Get("GoogleSheetsCellExpr")
	if err != nil {
		return model.CellProps{}, err
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return model.CellProps{}, err
	}

	props := re.FindStringSubmatch(message)
	if len(props) != 3 {
		return model.CellProps{}, fmt.Errorf("not matched")
	}

	return model.CellProps{SheetID: props[1], Range: props[2]}, nil
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}
