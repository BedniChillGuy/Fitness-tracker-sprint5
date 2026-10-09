package actioninfo

import (
	"fmt"
	"log"
)

// DataParser описывает тип, который умеет разбирать строку с данными
// и формировать информацию об активности.
type DataParser interface {
	Parse(datastring string) error
	ActionInfo() (string, error)
}

// Info разбирает каждую строку dataset и выводит информацию об активности.
func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		if err := dp.Parse(data); err != nil {
			log.Println(err)
			continue
		}

		info, err := dp.ActionInfo()
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println(info)
	}
}
