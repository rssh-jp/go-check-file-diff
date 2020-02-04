package checkfilediff

import (
	"bufio"
	"hash/crc64"
	"io"
	"log"
	"math/rand"
)

const (
	bufSize          = 1 * 1024 * 1024
	randomAccessSize = int(bufSize / 100)
)

func IsSame(f1, f2 io.Reader) (bool, error) {
	r1 := bufio.NewReader(f1)
	r2 := bufio.NewReader(f2)

	var b1, b2 []byte
	b1 = make([]byte, bufSize)
	b2 = make([]byte, bufSize)

	for {
		n1, err1 := r1.Read(b1)
		n2, err2 := r2.Read(b2)
		if err1 == io.EOF && err2 == io.EOF {
			break
		}

		if err1 != nil {
			return false, err1
		}

		if err2 != nil {
			return false, err2
		}

		if !isSame(b1[:n1], b2[:n2]) {
			return false, nil
		}
	}

	return true, nil
}

func IsMaybeSame(f1, f2 io.Reader) (bool, error) {
	r1 := bufio.NewReader(f1)
	r2 := bufio.NewReader(f2)

	var b1, b2 []byte
	b1 = make([]byte, bufSize)
	b2 = make([]byte, bufSize)

	for {
		n1, err1 := r1.Read(b1)
		n2, err2 := r2.Read(b2)
		if err1 == io.EOF && err2 == io.EOF {
			break
		}

		if err1 != nil {
			return false, err1
		}

		if err2 != nil {
			return false, err2
		}

		if randomAccessSize > n1 {
			if !isSame(b1[:n1], b2[:n2]) {
				return false, nil
			}
		} else {
			if !isMaybeSame(b1[:n1], b2[:n2]) {
				return false, nil
			}
		}
	}

	return true, nil
}

func isSame(a1, a2 []byte) bool {
	for i := 0; i < len(a1); i++ {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

func isSame2(a1, a2 []byte) bool {
	table := crc64.MakeTable(10)
	ui1 := crc64.Checksum(a1, table)
	ui2 := crc64.Checksum(a2, table)

	if ui1 != ui2 {
		return false
	}

	return true
}

func isMaybeSame(a1, a2 []byte) bool {
	for i := 0; i < randomAccessSize; i++ {
		index := rand.Intn(len(a1))
		if a1[index] != a2[index] {
			return false
		}
	}
	return true
}

func Fu(f1 io.Reader) {
	r1 := bufio.NewReader(f1)
	b1 := make([]byte, 1024)
	for {
		n1, err1 := r1.Read(b1)
		log.Println(n1, err1)
		if err1 != nil {
			break
		}
	}
}
