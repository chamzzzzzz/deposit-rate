package main

type BOC struct {
}

func (i *BOC) Name() string {
	return "BOC"
}

func (i *BOC) GetDepositRate() (*DepositRate, error) {
	return nil, nil
}
