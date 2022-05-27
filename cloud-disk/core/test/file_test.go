package test

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"testing"
)

// file chunk size, current: 100M
const chunkSize = 100 * 1024 * 1024

// divide file into chunks
func TestGenerateChunkFile(t *testing.T) {

	// get the detailed file information through os.Stat func
	fileInfo, err := os.Stat("test.dsf")
	if err != nil {
		t.Fatal(err)
	}

	// get the chunk num
	chunkNum := math.Ceil(float64(fileInfo.Size()) / chunkSize)

	myFile, err := os.OpenFile("test.dsf", os.O_RDONLY, 0666)
	defer myFile.Close()
	if err != nil {
		t.Fatal(err)
	}

	b := make([]byte, chunkSize)

	// each loop will create a file slice
	for i := 0; i < int(chunkNum); i++ {

		// specify the starting point of the read file
		myFile.Seek(int64(i*chunkSize), 0)

		// at the last time, the file still unread may smaller than chunkSize
		if chunkSize > fileInfo.Size()-int64(i*chunkSize) {
			b = make([]byte, fileInfo.Size()-int64(i*chunkSize))
		}

		myFile.Read(b)

		// store file slice
		f, err := os.OpenFile("./"+strconv.Itoa(i)+".chunk", os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		f.Write(b)
		f.Close()
	}
	myFile.Close()
}

// merge file chunk
func TestMergeChunkFile(t *testing.T) {
	myFile, err := os.OpenFile("test2.dsf", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	// get the detailed file information through os.Stat func
	fileInfo, err := os.Stat("test.dsf")
	if err != nil {
		t.Fatal(err)
	}

	// get the chunk num
	chunkNum := math.Ceil(float64(fileInfo.Size()) / chunkSize)

	for i := 0; i < int(chunkNum); i++ {
		f, err := os.OpenFile("./"+strconv.Itoa(i)+".chunk", os.O_RDONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		myFile.Write(b)
		f.Close()
	}
	myFile.Close()
}

// file consistency check
func TestCheckConsisitecy(t *testing.T) {

	// get the first file's information
	file1, err := os.OpenFile("test.dsf", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}

	b1, err := ioutil.ReadAll(file1)
	if err != nil {
		t.Fatal(err)
	}

	file2, err := os.OpenFile("test2.dsf", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}

	b2, err := ioutil.ReadAll(file2)
	if err != nil {
		t.Fatal(err)
	}

	s1 := fmt.Sprintf("%x", md5.Sum(b1))
	s2 := fmt.Sprintf("%x", md5.Sum(b2))

	fmt.Println(s1, s2, s1 == s2)
}
