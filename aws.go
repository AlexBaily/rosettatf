package main

type volume struct {

}

type instance struct {
    ami, ebs_optimized, disable_api_termination, instance_type string
    key_name, subnet_id, private_ip, iam_instance_profile, root_block_device string
    vpc_security_group_ids, tags []string
    ebs_block_device []volume

}

