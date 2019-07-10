package main

import "flag"
import "fmt"
import "os"

func main() {

    service := flag.String("service", "", "REQUIRED: The AWS Service to parse.")
    instance := flag.String("instance-id", "", "Required if using EC2 instances.")
    flag.Parse()

    if *service == "" {
        fmt.Println("Error, missing required arguement:")
        flag.PrintDefaults()
        os.Exit(1)
    }

    ec2 := queryEc2(*instance)
    fmt.Println(ec2)
}
