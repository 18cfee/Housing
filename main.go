package main

import (
	"fmt"
	"math"
)

func main() {
	months := calculateYearsToProduceYield(0.1375124, 0.03375, 0.0325, 0.00625, 30)
	fmt.Println(months)
}

func calculateYearsToProduceYield(yield float64, origIntRate float64, newIntRate float64, buyDownPoints float64, years int) float64 {
	origSched := getSched(origIntRate, years)
	newSched := getSched(newIntRate, years)
	origTCash := buyDownPoints
	newTCash := 0.0
	months := years * 12
	diffPayment := origSched[0].payment - newSched[0].payment
	monthlyYield := math.Pow((yield + 1.0), (1.0 / 12.0))
	for i := 0; i < months; i++ {
		origTCash *= monthlyYield
		newTCash *= monthlyYield
		newTCash += diffPayment
		if newTCash-newSched[i].balance-origTCash+origSched[i].balance > 0 {
			return (float64(i) + 1.0)
		}
	}
	return 0.0
}

type month struct {
	interest  float64
	principal float64
	balance   float64
	payment   float64
}

func getSched(interestRate float64, years int) []month {
	monthlyPay := getMonthlyPayment(years, interestRate)
	months := years * 12
	sched := make([]month, months)
	rest := interestRate / 12
	debt := 1.0
	for i := 0; i < months; i++ {
		interest := rest * debt
		principalPay := monthlyPay - interest
		debt -= principalPay
		sched[i] = month{interest: interest, principal: principalPay, balance: debt, payment: monthlyPay}
	}
	return sched
}

func getMonthlyPayment(years int, interestRate float64) float64 {
	lowerBound := interestRate * 1.000001 / 12
	upperBound := (interestRate + 1.0/float64(years)) / 12
	res := 1.0
	deadband := 0.0000000000001
	for true {
		monthlyPay := (lowerBound + upperBound) / 2
		res = runScenarioWithMonthlyPay(years, interestRate/12, monthlyPay)
		if res < 0.0 {
			upperBound = monthlyPay
		} else {
			lowerBound = monthlyPay
		}
		if res < deadband && res > -deadband {
			return monthlyPay
		}
	}
	return 0.0
}

func runScenarioWithMonthlyPay(years int, interestRate float64, monthlyPay float64) float64 {
	months := years * 12
	debt := 1.0
	for i := 0; i < months; i++ {
		debt = debt*(1.0+interestRate) - monthlyPay
	}
	return debt
}
