package cephalopod

import "day6/operator"

type CephalopodRow struct {
	nums     []int64
	operator operator.Operator
}

func (c CephalopodRow) Calculate() int64 {
	result := c.operator.Identity()
	for _, x := range c.nums {
		result = c.operator.Merge(result, x)
	}

	return result
}

type CephalopodRowBuilder struct {
	row *CephalopodRow
}

func Builder() *CephalopodRowBuilder {
	return &CephalopodRowBuilder{
		row: &CephalopodRow{},
	}
}

func (c *CephalopodRowBuilder) Nums(nums ...int64) *CephalopodRowBuilder {
	if c.row == nil {
		c.row = &CephalopodRow{}
	}

	c.row.nums = append(c.row.nums, nums...)
	return c
}

func (c *CephalopodRowBuilder) Operator(op operator.Operator) *CephalopodRowBuilder {
	if c.row == nil {
		c.row = &CephalopodRow{}
	}

	c.row.operator = op
	return c
}

func (c *CephalopodRowBuilder) Build() CephalopodRow {
	return *c.row
}
