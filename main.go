package main

import "flag"
import "fmt"
import "os"

func main() {

    service := flag.String("service", "", "REQUIRED: The AWS Service to parse.")
    flag.Parse()

    if *service == "" {
        fmt.Println("Error, missing required arguement:")
        flag.PrintDefaults()
        os.Exit(1)
    }

}
