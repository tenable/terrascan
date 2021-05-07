package accurics

ec2ebsnotoptimized[con.id] {
	con = input.aws_instance[_]
	object.get(con.config, "ebs_optimized", "undefined") == "undefined"
}
