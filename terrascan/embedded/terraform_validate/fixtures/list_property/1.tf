resource "aws_autoscaling_group" "asg" {
  tag {
    key                 = "Name"
    value               = "My ASG"
    propagate_at_launch = true
  }
}

resource "aws_autoscaling_group" "asg2" {
  tag {
    key                 = "CostCode"
    value               = "1234"
    propagate_at_launch = true
  }
  tag {
    key                 = "Name"
    value               = "My ASG"
    propagate_at_launch = true
  }
}
