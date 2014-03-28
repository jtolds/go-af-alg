package sha1

import (
	"crypto/rand"
	"crypto/sha1"
	"io"
	"testing"
)

func TestSHA1(t *testing.T) {
	for i := 0; i < 100; i++ {
		buf := make([]byte, 10*1024-i)
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			t.Fatal(err)
		}

		expected := sha1.Sum(buf)
		got, err := SHA1(buf)
		if err != nil {
			t.Fatal(err)
		}

		if expected != got {
			t.Fatal("exp:%x got:%x", expected, got)
		}
	}
}

func TestSHA1Writer(t *testing.T) {
	hash := sha1.New()

	for i := 0; i < 100; i++ {
		ohash, err := New()
		if err != nil {
			t.Fatal(err)
		}
		hash.Reset()
		buf := make([]byte, 10*1024-i)
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			t.Fatal(err)
		}

		if _, err := ohash.Write(buf); err != nil {
			t.Fatal(err)
		}
		if _, err := hash.Write(buf); err != nil {
			t.Fatal(err)
		}

		var got, exp [20]byte

		hash.Sum(exp[:0])
		got, err = ohash.Sum()
		if err != nil {
			t.Fatal(err)
		}

		if got != exp {
			t.Fatal("exp:%x got:%x", exp, got)
		}
	}
}

type shafunc func([]byte)

func benchmarkSHA1(b *testing.B, length int64, fn shafunc) {
	buf := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		b.Fatal(err)
	}
	b.SetBytes(length)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(buf)
	}
}

func BenchmarkSHA1Large(b *testing.B) {
	benchmarkSHA1(b, 1024*1024, func(buf []byte) { SHA1(buf) })
}

func BenchmarkSHA1Large_stdlib(b *testing.B) {
	benchmarkSHA1(b, 1024*1024, func(buf []byte) { sha1.Sum(buf) })
}

func BenchmarkSHA1Small(b *testing.B) {
	benchmarkSHA1(b, 1, func(buf []byte) { SHA1(buf) })
}

func BenchmarkSHA1Small_stdlib(b *testing.B) {
	benchmarkSHA1(b, 1, func(buf []byte) { sha1.Sum(buf) })
}
