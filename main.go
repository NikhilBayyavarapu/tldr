package main

import (
	"awsapp/uploadtos3"
	"flag"
	"fmt"
	"log"

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

	fmt.Println(*path)
	uploadtos3.UploadToS3(*path)
}
