// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webp

import (
	"bytes"
	_ "image/png"
	"io/ioutil"
	"testing"
)

type tTester struct {
	Filename string
	Lossless bool
	Quality  float32 // 0 ~ 100
	MaxDelta int
}

var tTesterList = []tTester{
	tTester{
		Filename: "video-001.png",
		Lossless: false,
		Quality:  90,
		MaxDelta: 5,
	},
	tTester{
		Filename: "1_webp_ll.png",
		Lossless: false,
		Quality:  90,
		MaxDelta: 5,
	},
	tTester{
		Filename: "2_webp_ll.png",
		Lossless: true,
		Quality:  90,
		MaxDelta: 0,
	},
	tTester{
		Filename: "3_webp_ll.png",
		Lossless: false,
		Quality:  90,
		MaxDelta: 5,
	},
	tTester{
		Filename: "4_webp_ll.png",
		Lossless: true,
		Quality:  90,
		MaxDelta: 0,
	},
	tTester{
		Filename: "5_webp_ll.png",
		Lossless: false,
		Quality:  75,
		MaxDelta: 15,
	},
}

func TestEncode(t *testing.T) {
	for i, v := range tTesterList {
		img0, err := loadImage(testdataDir + v.Filename)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}

		buf := new(bytes.Buffer)
		err = Encode(buf, img0, &Options{
			Lossless: v.Lossless,
			Quality:  v.Quality,
		})
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}

		img1, err := Decode(buf)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}

		// Compare the average delta to the tolerance level.
		want := v.MaxDelta
		if got := averageDelta(img0, img1); got > want {
			t.Fatalf("%d: average delta too high; got %d, want <= %d", i, got, want)
		}
	}
}

// BenchmarkEncode benchmarks the encoding of an image.
func BenchmarkEncode(b *testing.B) {
	img, err := loadImage(testdataDir + "1_webp_ll.png")
	if err != nil {
		b.Fatal(err)
	}
	s := img.Bounds().Size()
	b.SetBytes(int64(s.X * s.Y * 4))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(ioutil.Discard, img, nil)
	}
}
