package alsa

import (
	"fmt"
	"os"
)

type Card struct {
	Path   string
	Title  string
	Number int

	fh       *os.File
	pversion pVersion
	cardinfo cardInfo
}

func (card Card) String() string {
	return card.Title
}

func OpenCards() ([]*Card, error) {
	ret := make([]*Card, 0)

	max := 3 // arbitrary
	for i := 0; i < max; i++ {
		path := fmt.Sprintf("/dev/snd/controlC%d", i)
		_, err := os.Stat(path)
		if err != nil {
			continue
		}
		max++

		fh, err := os.Open(path)
		if err != nil {
			return ret, err
		}

		card := Card{
			Path:   path,
			Number: i,
			fh:     fh,
		}

		err = ioctl(fh.Fd(), ioctl_encode(cmdRead, 4, cmdControlVersion), &card.pversion)
		if err != nil {
			return ret, err
		}

		err = ioctl(fh.Fd(), ioctl_encode(cmdRead, 376, cmdControlCardInfo), &card.cardinfo)
		if err != nil {
			return ret, err
		}

		card.Title = gstr(card.cardinfo.Name[:])
		ret = append(ret, &card)
	}

	return ret, nil
}

func CloseCards(cards []*Card) {
	for _, card := range cards {
		card.fh.Close()
	}
}
