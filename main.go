package main

import (
	"awsapp/uploadtos3"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to env file")
	}

	path := flag.String("path", " ", "Use this flag to give path of the document")
	flag.Parse()

	if *path == " " {
		log.Fatalln("You need to specify a file path to run this applcation.")
	}

	file, _ := os.Open(*path)
	defer file.Close()

	bytes, _ := io.ReadAll(file)
	text := string(bytes)
	uploadtos3.Summarize(text)
	return

	fmt.Println(os.Getenv("AWS_REGION"))
	uploadtos3.UploadToS3()
}
