package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "fmt"
)

type volume struct {

}

type instance struct {
    ami, ebs_optimized, disable_api_termination, instance_type string
    key_name, subnet_id, private_ip, iam_instance_profile, root_block_device string
    vpc_security_group_ids, tags []string
    ebs_block_device []volume

}

var sess *session.Session

func init() {
    sess = session.Must(session.NewSession(&aws.Config{
        Region: aws.String("eu-west-1")}))
}

func queryEc2(instanceId string) (string) {
    svc := ec2.New(sess)
    input := &ec2.DescribeInstancesInput{
        InstanceIds: []*string{
            aws.String(instanceId),
        },
    }
    result, err := svc.DescribeInstances(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
    }
    instance := result.Reservations[0].Instances
    return instance[0].GoString()
}
