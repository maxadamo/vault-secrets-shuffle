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
	//contentByte, err := ioutil.ReadFile(fileName)
	//content := string(contentByte)
	if err != nil {
		log.Fatal(err)
	}

	fzip, err := os.Create(zipName)
	if err != nil {
		log.Fatalln(err)
	}

	zipw := zip.NewWriter(fzip)

	defer zipw.Close()

	w, err := zipw.Encrypt(fileName, zipPassword, zip.AES256Encryption)
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
	fmt.Printf("\npassword file saved as %v and encrypt with password %v\n", zipName, zipPassword)

}
