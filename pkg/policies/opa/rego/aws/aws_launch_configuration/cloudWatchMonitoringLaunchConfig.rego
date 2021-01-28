package accurics

cloudWatchMonitoringLaunchConfig[res.id] {
    res = input.aws_launch_configuration[_]
	res.config.enable_monitoring == false
}