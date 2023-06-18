package hackbrowserdata

import (
	"testing"
)

func TestChromium_Init(_ *testing.T) {
}

func BenchmarkChromium_Init(b *testing.B) {
	chromium := browsers[Chrome]
	for i := 0; i < b.N; i++ {
		if err := chromium.Init(); err != nil {
			b.Fatal(err)
		}
	}
}
