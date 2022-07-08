package graph

import (
	"github.com/elementsproject/glightning/glightning"
	"strconv"
	"strings"
)

type Channel struct {
	glightning.Channel
	Liquidity uint64 `json:"liquidity"`
}

func (c *Channel) computeFee(amount uint64) uint64 {
	baseFee := c.BaseFeeMillisatoshi
	result := baseFee
	proportionalFee := ((amount / 1000) * c.FeePerMillionth) / 1000
	result += proportionalFee
	return result
}

func (c *Channel) GetHop(amount uint64, delay uint) glightning.RouteHop {
	return glightning.RouteHop{
		Id:             c.Destination,
		ShortChannelId: c.ShortChannelId,
		MilliSatoshi:   amount,
		Delay:          delay,
		Direction:      c.getDirection(),
	}
}

func (c *Channel) getDirection() uint8 {
	if c.Source < c.Destination {
		return 0
	}
	return 1
}

func all(v []bool) bool {
	for _, b := range v {
		if !b {
			return false
		}
	}
	return true
}

func (c *Channel) canUse(amount uint64) bool {
	maxHtlcMsat, _ := strconv.ParseUint(strings.TrimSuffix(c.HtlcMaximumMilliSatoshis, "msat"), 10, 64)
	conditions := []bool{
		c.Liquidity >= amount,
		c.IsActive,
		maxHtlcMsat >= amount,
	}
	return all(conditions)
}
