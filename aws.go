package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "fmt"
    "strconv"
    "reflect"
)

//EBS Volume struc for captring
type volume struct {

}

//EC2 Instance Volumes Struct
type instance struct {
    ami, ebs_optimized, disable_api_termination, instance_type string
    key_name, subnet_id, private_ip, iam_instance_profile, root_block_device string
    vpc_security_group_ids, tags []string
    ebs_block_device []volume

}

//String function to convert the instance Struct to a terraform string
func (inst instance) String() string {
    r := reflect.ValueOf(inst)
    s := fmt.Sprintf("resource \"aws_instance\" \"test\" {\n")
    for i := 0; i < r.Type().NumField(); i++ {
        if r.Field(i).String() != "" {
            s += fmt.Sprintf("    %s = \"%s\"\n", r.Type().Field(i).Name, r.Field(i))
        }
    }
    s += "}"
    return s
}

//Create the session var that will be used throught the package
var sess *session.Session

//Initliase the session in init()
func init() {
    sess = session.Must(session.NewSession(&aws.Config{
        Region: aws.String("eu-west-1")}))
}

//Query EC2 for the information about the instance ID and return the instance String.
func queryEc2(instanceId string) (string) {
    svc := ec2.New(sess)
    //Set config for getting instance information
    input := &ec2.DescribeInstancesInput{
        InstanceIds: []*string{
            aws.String(instanceId),
        },
    }
    //get instance information and deal with errors.
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
    //These two blocks are due to the fact that you cannot get api termination via normal describe
    attrInput := &ec2.DescribeInstanceAttributeInput{
        Attribute:  aws.String("disableApiTermination"),
        InstanceId: aws.String(instanceId),
    }
    //Get the attribute input
    attrResult, err := svc.DescribeInstanceAttribute(attrInput)
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
    //EC2 Instane are returned as reservations which can have multiple EC2 instances;
    //In this case we are only looking at the first reservation and instance in that reservation.
    r := result.Reservations[0].Instances[0]
    //Get the boolean value of disable API termination.
    attrR := attrResult.DisableApiTermination.Value
    //Create the instance struct based on the information from the EC2 instance.
    //Need to deference the attributes.
    instanceStruct := instance{
        ami: *r.ImageId,
        ebs_optimized: strconv.FormatBool(*r.EbsOptimized),
        disable_api_termination: strconv.FormatBool(*attrR),
        subnet_id: *r.SubnetId,
    }

    return instanceStruct.String()
}
