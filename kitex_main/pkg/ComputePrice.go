package pkg

import (
	"kitex_main/model"
	"strconv"
)

func ComputePrices(Distance float64, CarType string) (price float64, err error) {
	var LxhPricingRules model.LxhPricingRules

	err = LxhPricingRules.FindLxhPricingRulesCarType(CarType)
	if err != nil {
		return 0.0, err
	}

	Price, _ := strconv.ParseFloat(LxhPricingRules.Price, 64)
	Kilometers, _ := strconv.ParseFloat(LxhPricingRules.Kilometers, 64)

	SumPaice := Price * Kilometers

	Kilometer := Distance - Kilometers

	Excess, _ := strconv.ParseFloat(LxhPricingRules.Excess, 64)

	Sum := Kilometer * Excess

	return Sum + SumPaice, nil

}
