package checkfilediff

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

const (
	temporary = "tmp/"
)

const (
	testTypeNone = iota + 1
	testTypeCache
)

var (
	ErrCouldNotMatchFileSize = errors.New("Could not match file size")
)

var (
	testType int
)

func preprocess() {
	os.Mkdir(temporary, 0755)
}

func postprocess() {
	switch testType {
	case testTypeNone:
		err := os.RemoveAll(temporary)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestMain(m *testing.M) {
	flag.IntVar(&testType, "t", 1, "type. 1 is none, 2 is cache. default 1")
	flag.Parse()

	preprocess()

	defer postprocess()

	m.Run()
}

func isExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

func createFile(path string, filesize, seed int64) error {
	switch testType {
	case testTypeCache:
		if isExist(path) {
			return nil
		}
	}

	r := rand.New(rand.NewSource(seed))

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer fd.Close()

	const size = 1024 * 1024

	var worksize int64

	for {
		if worksize >= filesize {
			break
		}

		s := int64(size)

		if s+worksize > filesize {
			s = filesize - worksize
		}

		buf := make([]byte, s)

		r.Read(buf)

		n, err := fd.Write(buf)
		if err != nil {
			return err
		}

		if n != len(buf) {
			return ErrCouldNotMatchFileSize
		}

		worksize += s
	}

	return nil
}

type test interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

func testSuccess(t test, filepath1, filepath2 string, size int64) {
	createFile(filepath1, size, 10)
	createFile(filepath2, size, 10)

	switch testType {
	case testTypeNone:
		defer os.Remove(filepath1)
		defer os.Remove(filepath2)
	}

	s := time.Now()

	fd1, err := os.Open(filepath1)
	if err != nil {
		t.Fatal(err)
	}

	defer fd1.Close()

	fd2, err := os.Open(filepath2)
	if err != nil {
		t.Fatal(err)
	}

	defer fd2.Close()

	same, err := IsSame(fd1, fd2)
	if err != nil {
		t.Fatal(err)
	}

	if !same {
		t.Error("Not match file")
	}

	t.Log(time.Now().Sub(s))
}

func TestIsSameSuccess(t *testing.T) {
	t.Run("10 Byte", func(t *testing.T) {
		const filepath1 = temporary + "byte1"
		const filepath2 = temporary + "byte2"

		const size = 10

		testSuccess(t, filepath1, filepath2, size)
	})

	t.Run("1 Killo byte", func(t *testing.T) {
		const filepath1 = temporary + "killo1"
		const filepath2 = temporary + "killo2"

		const size = 1 * 1024

		testSuccess(t, filepath1, filepath2, size)
	})

	t.Run("1 Mega byte", func(t *testing.T) {
		const filepath1 = temporary + "mega1"
		const filepath2 = temporary + "mega2"

		const size = 1 * 1024 * 1024

		testSuccess(t, filepath1, filepath2, size)
	})

	t.Run("500 Mega byte", func(t *testing.T) {
		const filepath1 = temporary + "halfgiga1"
		const filepath2 = temporary + "halfgiga2"

		const size = 512 * 1024 * 1024

		testSuccess(t, filepath1, filepath2, size)
	})

	t.Run("1 Giga byte", func(t *testing.T) {
		const filepath1 = temporary + "giga1"
		const filepath2 = temporary + "giga2"

		const size = 1 * 1024 * 1024 * 1024

		testSuccess(t, filepath1, filepath2, size)
	})
}

func BenchmarkSame(b *testing.B) {
	b.Run("10 Byte", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			const filepath1 = temporary + "byte1"
			const filepath2 = temporary + "byte2"

			const size = 10

			testSuccess(b, filepath1, filepath2, size)
		}
	})

	b.Run("1 Killo Byte", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			const filepath1 = temporary + "killo1"
			const filepath2 = temporary + "killo2"

			const size = 1 * 1024

			testSuccess(b, filepath1, filepath2, size)
		}
	})

	b.Run("1 Mega Byte", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			const filepath1 = temporary + "mega1"
			const filepath2 = temporary + "mega2"

			const size = 1 * 1024 * 1024

			testSuccess(b, filepath1, filepath2, size)
		}
	})

	b.Run("500 Mega Byte", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			const filepath1 = temporary + "halfgiga1"
			const filepath2 = temporary + "halfgiga2"

			const size = 512 * 1024 * 1024

			testSuccess(b, filepath1, filepath2, size)
		}
	})

	b.Run("1 Giga Byte", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			const filepath1 = temporary + "giga1"
			const filepath2 = temporary + "giga2"

			const size = 1 * 1024 * 1024 * 1024

			testSuccess(b, filepath1, filepath2, size)
		}
	})
}
