package processor

import (
	"fmt"
	"gsheets-slack-chatbot/model/google"
	"regexp"
)

func (p *Processor) getCellContent(message string) (string, error) {
	props, err := p.getCellProps(message)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", props.SheetID, props.Range), nil
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
