package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "fmt"
    "strconv"
    "reflect"
    "strings"
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

//Create a list of the optional EC2 Values.
var (
    ec2Optional = [4]string{"iam_instance_profile", "vpc_security_groups_ids",
        "tags", "ebs_block_device"}
)

//String function to convert the instance Struct to a terraform string
func (inst instance) String() string {
    r := reflect.ValueOf(inst)
    s := fmt.Sprintf("resource \"aws_instance\" \"test\" {\n")
    for i := 0; i < r.Type().NumField(); i++ {
        if r.Field(i).String() != "" &&  r.Field(i).Type().Kind() != reflect.Slice {
            s += fmt.Sprintf("    %s = \"%s\"\n", r.Type().Field(i).Name, r.Field(i))
        } else if r.Field(i).Type().Kind() == reflect.Slice {
            fmt.Println("Slice found")
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
        instance_type: *r.InstanceType,
        key_name: *r.KeyName,
        private_ip: *r.PrivateIpAddress,
    }

    //Get a reflection of the AWS SDK instance so we can check the fields
    //Here we can check for optional values such as IAM instance profile or SGs
    rInst := reflect.ValueOf(*r)

    //Relect Value.Type() interface 'FieldByName' will return a StructField and a boolean
    //The boolean will determine whether the field is found.
    _, iamProfileBool := rInst.Type().FieldByName("IamInstanceProfile")
    if iamProfileBool {
        profileName := strings.Split(*r.IamInstanceProfile.Arn, "/")
        instanceStruct.iam_instance_profile = profileName[1]
    }
    //Creating the string slice to house the SG array
    sgArray := make([]string, len(r.SecurityGroups))
    for i := 0; i < len(r.SecurityGroups); i++ {
        sgArray[i] = *r.SecurityGroups[i].GroupId
    }
    instanceStruct.vpc_security_group_ids = sgArray

    return instanceStruct.String()
}
