package constant

import (
	"testing"
)

func TestConstant(t *testing.T) {

	cases := []struct {
		in, want int32
	}{
		{2, StarHyperGiant},
		{4, StarSuperGiant},
		{8, StarBrightGiant},
		{16, StarGiant},
		{32, StarSubGiant},
		{128, StarDwarf},
		{256, StarSubDwarf},
		{512, StarWhiteDwarf},
		{1024, StarRedDwarf},
		{2048, StarBrownDwarf},
	}

	for _, c := range cases {
		if c.in != c.want {
			t.Fail()
		}

	}
}
