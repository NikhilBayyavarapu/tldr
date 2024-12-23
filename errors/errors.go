package errors

import "log"

func HandleError(err error, function string) {
	if err != nil {
		log.Fatalf("Error in %s function. Error: %v", function, err)
	}
}
