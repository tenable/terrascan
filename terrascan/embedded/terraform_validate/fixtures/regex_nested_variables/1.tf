resource "aws_instance" "foo" {
    tags {
        CPM_Service_wibble = 1
    }
}