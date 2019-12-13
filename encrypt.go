package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/yeka/zip"
)

func randPass(lenght int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	var b strings.Builder
	for i := 0; i < lenght; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// zipEncrypt file
func zipEncrypt(fileName string, passLenght int) {
	zipPassword := randPass(passLenght)
	zipName := fmt.Sprintf("%v.zip", fileName)
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fzip, err := os.Create(zipName)
	if err != nil {
		log.Fatalln(err)
	}

	zipw := zip.NewWriter(fzip)

	defer zipw.Close()

	w, err := zipw.Encrypt(fileName, zipPassword, zip.StandardEncryption)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(w, bytes.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	zipw.Flush()

	err = os.Remove(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(strings.Repeat("=", 72))
	fmt.Printf("\nThe passwords have been saved to: %v\nThe password to decrypt the zip file is: %v\n", zipName, zipPassword)

}
