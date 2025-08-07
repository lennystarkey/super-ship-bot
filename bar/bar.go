package bar

import "strings"

var (
	bStart  = "<:border_start:1403112918861615305>"
	bMiddle = "<:border_middle:1403112917494399117>"
	bEnd    = "<:border_end:1403112915489525933>"

	r = assetSet{
		startEnd: "<:red_start_end:1403112904584335531>",
		start:    "<:red_start:1403112901971152896>",
		middle:   "<:red_middle:1403112898880082012>",
		end:      "<:red_end:1403112931456979075>",
		bEnd:     "<:red_end_border_end:1403112896308707359>",
	}
	y = assetSet{
		startEnd: "<:yellow_start_end:1403112914604396746>",
		start:    "<:yellow_start:1403112912066973706>",
		middle:   "<:yellow_middle:1403112910498037810>",
		end:      "<:yellow_end:1403112905884307547>",
		bEnd:     "<:yellow_end_border_end:1403112908572852435>",
	}
	g = assetSet{
		startEnd: "<:green_start_end:1403112929485914253>",
		start:    "<:green_start:1403113067000369306>",
		middle:   "<:green_middle:1403112925631086763>",
		end:      "<:green_end:1403113060444541119>",
		bEnd:     "<:green_end_border_end:1403112922468843661>",
	}
)

type assetSet struct {
	startEnd string
	start    string
	middle   string
	end      string
	bEnd     string
}

func Generate(p float64) string {
	var c assetSet
	switch {
	case p < 0.33:
		c = r
	case p < 0.66:
		c = y
	default:
		c = g
	}
	n := int(p*17 + 0.5)
	switch {
	case n <= 0:
		return bStart + strings.Repeat(bMiddle, 15) + bEnd
	case n == 1:
		return c.startEnd + strings.Repeat(bMiddle, 15) + bEnd
	case n >= 17:
		return c.start + strings.Repeat(c.middle, 15) + c.end
	default:
		return c.start + strings.Repeat(c.middle, n-2) + c.end + strings.Repeat(bMiddle, 16-n) + bEnd
	}
}
