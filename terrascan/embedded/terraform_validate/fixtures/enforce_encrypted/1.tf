
resource "aws_db_instance_valid" "foo1" {
    storage_encrypted = "True"
}

resource "aws_db_instance_invalid" "foo2" {
    storage_encrypted = "False"
}

resource "aws_db_instance_invalid2" "foo3" {
    wibble = 1
}

resource "aws_instance_valid" "bizz1" {
    ebs_block_device {
        encrypted = "True"
    }
}

resource "aws_instance_invalid" "bizz2" {
    ebs_block_device {
        encrypted = "False"
    }
}

resource "aws_instance_invalid2" "bizz3" {
   ebs_block_device {
        wibble = 1
    }
}

resource "aws_ebs_volume_valid" "bar1" {
    encrypted = "True"
}

resource "aws_ebs_volume_invalid" "bar2" {
    encrypted = "False"
}

resource "aws_ebs_volume_invalid2" "bar3" {
    wibble = 1
}
