package processor

import (
	"encoding/json"
	"fmt"
	"gsheets-slack-chatbot/model/google"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func (p *Processor) getCellContent(message string) (string, error) {
	where := "processor.Processor.getCellContent(...)"

	p.log.Trace(where, "Retrieving cell properties.")
	props, err := p.getCellProps(message)
	if err != nil {
		return "", err
	}

	fileName, err := p.config.Get("GoogleCredentialsFileName")
	if err != nil {
		return "", err
	}

	p.log.Trace(where, fmt.Sprintf("Loading Google credentials file. File name: %s", fileName))
	creds, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	p.log.Trace(where, "Creating config based on Google credentials file.")
	config, err := google.ConfigFromJSON(creds, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return "", err
	}

	p.log.Trace(where, "Getting HTTP client to access Google Sheets API.")
	client, err := p.getClient(config)
	if err != nil {
		return "", err
	}

	p.log.Trace(where, "Creating service for HTTP client.")
	srv, err := sheets.New(client)
	if err != nil {
		return "", err
	}

	p.log.Trace(where, fmt.Sprintf("Requesting for cell content. Sheet ID: %s, range: %s",
		props.SheetID, props.Range))
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

	return "", fmt.Errorf("unhandled error occurred")
}

func (p *Processor) getCellProps(message string) (model.CellProps, error) {
	where := "processor.Processor.getCellProps(...)"

	pattern, err := p.config.Get("GoogleSheetsCellExpr")
	if err != nil {
		return model.CellProps{}, err
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return model.CellProps{}, err
	}

	p.log.Trace(where, "Searching for cell properties inside the message.")
	props := re.FindStringSubmatch(message)
	if len(props) != 3 {
		p.log.Trace(where, "Cell properties occurrence was not found.")
		return model.CellProps{}, fmt.Errorf("not matched")
	}
	p.log.Trace(where, fmt.Sprintf("Found cell properties. Sheet ID: %s, range: %s.", props[1], props[2]))

	return model.CellProps{SheetID: props[1], Range: props[2]}, nil
}

func (p *Processor) getClient(config *oauth2.Config) (*http.Client, error) {
	where := "processor.Processor.getClient(...)"

	fileName, err := p.config.Get("GoogleTokenFileName")
	if err != nil {
		return nil, err
	}

	p.log.Trace(where, fmt.Sprintf("Retrieving token from file: %s", fileName))
	tok, err := tokenFromFile(fileName)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), tok), nil
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
