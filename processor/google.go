package processor

import "gsheets-slack-chatbot/model/google"

func (p *Processor) getCellContent(message string) (string, error) {
	return message, nil
}

func getCellProps(message string) (model.CellProps, error) {
	return model.CellProps{}, nil
}
