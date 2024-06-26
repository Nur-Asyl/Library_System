package reader

import "errors"

type Reader struct {
	FIO     string
	Address string
	NOMBIL  int
}

func NewReader(fio, address string, nombil int) (*Reader, error) {
	if err := checkFIO(fio); err != nil {
		return nil, err
	}
	if err := checkAddress(address); err != nil {
		return nil, err
	}

	return &Reader{
		FIO:     fio,
		Address: address,
		NOMBIL:  nombil,
	}, nil
}

func checkFIO(fio string) error {
	if len(fio) > 30 {
		return errors.New("FIO exceeds 30 characters")
	}
	return nil
}

func checkAddress(address string) error {
	if len(address) > 40 {
		return errors.New("address exceeeds 40 characters")
	}
	return nil
}
