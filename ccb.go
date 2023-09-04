package main

type CCB struct {
}

func (i *CCB) Name() string {
	return "CCB"
}

func (i *CCB) GetDepositRate() (*DepositRate, error) {
	return nil, nil
}
