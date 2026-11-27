package suffixarray

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestGoLeanSuffixarrayBuckets(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SUFFIXARRAY_BUCKET_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for row := 0; ; row++ {
		var c struct {
			Profile, FSize, BSize   int
			Text, Frequency, Bucket []int64
			Want                    string
		}
		if err := rows.Decode(&c); err == io.EOF {
			if row == 0 {
				t.Fatal("empty suffixarray bucket fixtures")
			}
			t.Logf("twelve pinned declarations matched %d triple rows", row)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		var got string
		switch c.Profile {
		case 0:
			got = checkBucket(t, c.Text, c.Frequency, c.Bucket, c.FSize, c.BSize,
				freq_8_32, bucketMin_8_32, bucketMax_8_32)
		case 1:
			got = checkBucket(t, c.Text, c.Frequency, c.Bucket, c.FSize, c.BSize,
				freq_8_64, bucketMin_8_64, bucketMax_8_64)
		case 2:
			got = checkBucket(t, c.Text, c.Frequency, c.Bucket, c.FSize, c.BSize, freq_32, bucketMin_32, bucketMax_32)
		case 3:
			got = checkBucket(t, c.Text, c.Frequency, c.Bucket, c.FSize, c.BSize, freq_64, bucketMin_64, bucketMax_64)
		default:
			t.Fatal("unknown bucket profile")
		}
		if got != c.Want {
			t.Fatalf("row %d profile %d FSize=%d BSize=%d text=%v: native=%q want=%q",
				row, c.Profile, c.FSize, c.BSize, c.Text, got, c.Want)
		}
	}
}

type symbol interface{ ~byte | ~int32 | ~int64 }
type counter interface{ ~int32 | ~int64 }

func checkBucket[T symbol, W counter](t *testing.T, text, f, b []int64, fSize, bSize int,
	frequency func([]T, []W, []W) []W, minimum, maximum func([]T, []W, []W)) string {
	t.Helper()
	input := convert[T](text)
	var out []string
	for mode := range 3 {
		freq, bucket := convert[W](f), convert[W](b)
		if fSize >= 0 {
			freq = freq[:fSize]
		}
		if bSize >= 0 {
			bucket = bucket[:bSize]
		}
		status := observeBucket(func() string {
			switch mode {
			case 0:
				result := frequency(input, freq, bucket)
				source := freq
				if freq == nil {
					source = bucket
				}
				if cap(result) != cap(source) || cap(result) > 0 &&
					&result[:cap(result)][0] != &source[:cap(source)][0] {
					t.Fatal("frequency returned a different backing slice")
				}
				return bucketWire(result)
			case 1:
				minimum(input, freq, bucket)
			case 2:
				maximum(input, freq, bucket)
			}
			return "ok"
		})
		out = append(out, status+"/"+bucketWire(freq)+"/"+bucketWire(bucket))
	}
	return strings.Join(out, "|")
}

func convert[T symbol](data []int64) []T {
	if data == nil {
		return nil
	}
	result := make([]T, len(data))
	for i, n := range data {
		result[i] = T(n)
	}
	return result
}

func observeBucket(call func() string) (report string) {
	defer func() {
		if failure := recover(); failure != nil {
			if err, ok := failure.(runtime.Error); ok &&
				(strings.HasPrefix(err.Error(), "runtime error: index out of range") ||
					strings.HasPrefix(err.Error(), "runtime error: slice bounds out of range")) {
				report = "panic"
				return
			}
			panic(failure)
		}
	}()
	return call()
}

func bucketWire[T counter](view []T) string {
	if view == nil {
		return "nil"
	}
	data := view[:cap(view)]
	var runs []string
	for i := 0; i < len(data); {
		j := i + 1
		for j < len(data) && data[j] == data[i] {
			j++
		}
		runs = append(runs, fmt.Sprintf("%d.%d", data[i], j-i))
		i = j
	}
	return fmt.Sprintf("%d:%s", len(view), strings.Join(runs, ","))
}
