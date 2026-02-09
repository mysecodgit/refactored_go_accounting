package money

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	CurrentValueScale int64 = 100_000 // 5 decimals
	PreviousValueScale int64 = 100_000 // 5 decimals
	QtyScale   int64 = 100_000 // 5 decimals
	RateScale  int64 = 100_000 // 5 decimals
	MoneyScale int64 = 100     // cents

	MaxQty        = 9_999_999_999.99999
	MaxRate       = 9_999_999_999_999.99999
	MaxTotalMoney = 99_999_999_999_999.99 // 14 digits before decimal

)

/*
  ---------- helpers ----------
*/

func decimalPlaces(s string) int {
	s = strings.ReplaceAll(s, ",", "")
	if !strings.Contains(s, ".") {
		return 0
	}
	return len(s) - strings.Index(s, ".") - 1
}

func parseFloatStrict(s string) (float64, error) {
	if s == "" {
		return 0, errors.New("empty value")
	}

	// Remove thousands separators
	s = strings.ReplaceAll(s, ",", "")

	return strconv.ParseFloat(s, 64)
}

/*
  ---------- qty ----------
*/

func ParsePreviousValue(previousValueStr string) (int64, error) {
	previousValue, err := parseFloatStrict(previousValueStr)
	if err != nil {
		return 0, errors.New("invalid previous value")
	}

	if decimalPlaces(previousValueStr) > 5 {
		return 0, errors.New("previous value cannot have more than 5 decimals")
	}

	// if previousValue < 0 || previousValue > MaxQty {
	// 	return 0, errors.New("previous value out of allowed range")
	// }

	return int64(math.Round(previousValue * float64(PreviousValueScale))), nil
}

func ParseCurrentValue(currentValueStr string) (int64, error) {
	currentValue, err := parseFloatStrict(currentValueStr)
	if err != nil {
		return 0, errors.New("invalid current value")
	}

	if decimalPlaces(currentValueStr) > 5 {
		return 0, errors.New("current value cannot have more than 5 decimals")
	}

	return int64(math.Round(currentValue * float64(CurrentValueScale))), nil
}


/*
  ---------- qty ----------
*/

func ParseQty(qtyStr string) (int64, error) {
	qty, err := parseFloatStrict(qtyStr)
	if err != nil {
		return 0, errors.New("invalid qty")
	}

	if decimalPlaces(qtyStr) > 5 {
		return 0, errors.New("qty cannot have more than 5 decimals")
	}

	if qty < 0 || qty > MaxQty {
		return 0, errors.New("qty out of allowed range")
	}

	return int64(math.Round(qty * float64(QtyScale))), nil
}

/*
  ---------- rate ----------
*/

func ParseRate(rateStr string) (int64, error) {
	rate, err := parseFloatStrict(rateStr)
	if err != nil {
		return 0, errors.New("invalid rate")
	}

	if decimalPlaces(rateStr) > 5 {
		return 0, errors.New("rate cannot have more than 5 decimals")
	}

	if rate < 0 || rate > MaxRate {
		return 0, errors.New("rate out of allowed range")
	}

	return int64(math.Round(rate * float64(RateScale))), nil
}

func ParseUSDAmount(amountStr string) (int64, error) {
	amount, err := parseFloatStrict(amountStr)
	if err != nil {
		return 0, errors.New("invalid amount")
	}

	if decimalPlaces(amountStr) > 2 {
		return 0, errors.New("amount cannot have more than 2 decimals")
	}

	// TODO : implement max amount

	return int64(math.Round(amount * float64(MoneyScale))), nil
}

/*
  ---------- total ----------
*/

func CalculateTotalCents(qtyScaled, rateScaled int64) (int64, error) {
	totalCents := int64(math.Round(
		(float64(qtyScaled) * float64(rateScaled)) / 100_000_000,
	))

	if totalCents < 0 {
		return 0, errors.New("total cannot be negative")
	}

	return totalCents, nil
}

/*
  ---------- one-shot ----------
*/

type LineInput struct {
	Qty  string
	Rate string
	PreviousValue *string
	CurrentValue *string
}

type LineResult struct {
	QtyScaled  int64
	RateScaled int64
	TotalCents int64
	PreviousValueScaled *int64
	CurrentValueScaled *int64
}

func ConvertLineInput(in LineInput) (*LineResult, error) {

	var previousValueScaled *int64 = nil
	var currentValueScaled *int64 = nil

	if in.PreviousValue != nil {
		previousValue, err := ParsePreviousValue(*in.PreviousValue)
		if err != nil {
			return nil, err
		}
		previousValueScaled = &previousValue
	}
	if in.CurrentValue != nil {
		currentValue, err := ParseCurrentValue(*in.CurrentValue)
		if err != nil {
			return nil, err
		}
		currentValueScaled = &currentValue
	}

	qtyScaled, err := ParseQty(in.Qty)
	if err != nil {
		return nil, err
	}

	rateScaled, err := ParseRate(in.Rate)
	if err != nil {
		return nil, err
	}

	totalCents, err := CalculateTotalCents(qtyScaled, rateScaled)
	if err != nil {
		return nil, err
	}



	return &LineResult{
		QtyScaled:  qtyScaled,
		RateScaled: rateScaled,
		TotalCents: totalCents,
		PreviousValueScaled: previousValueScaled,
		CurrentValueScaled: currentValueScaled,
	}, nil
}


func FormatMoneyFromCents(cents int64) string {
	return strconv.FormatFloat(
		float64(cents)/float64(MoneyScale),
		'f',
		2,
		64,
	)
}

func FormatScaled5(value int64) string {
	return strconv.FormatFloat(
		float64(value)/100_000,
		'f',
		-1, // no trailing zeros unless needed
		64,
	)
}
