package operator

import "errors"

type Operator string

const (
	OperatorAdd  Operator = "+"
	OperatorMult Operator = "*"
)

func (o Operator) Merge(x, y int64) int64 {
	switch o {
	case OperatorAdd:
		return x + y
	case OperatorMult:
		return x * y
	}

	panic("should not be reached")
}

func (o Operator) Identity() int64 {
	switch o {
	case OperatorAdd:
		return 0
	case OperatorMult:
		return 1
	}

	panic("should not be reached")
}

func ParseOperator(opStr string) (Operator, error) {
	switch opStr {
	case string(OperatorAdd):
		return OperatorAdd, nil
	case string(OperatorMult):
		return OperatorMult, nil
	default:
		return "", errors.New("not valid operator")
	}
}
