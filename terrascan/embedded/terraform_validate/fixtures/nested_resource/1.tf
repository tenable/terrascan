resource "aws_instance" "foo" {

    nested_resource {
        value = 1
        value2 = 2
    }

    tags {
        value = 1
    }

}

resource "aws_elb" "bar" {

    tags {
        value = 1
    }

}