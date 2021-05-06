variable "filename" {
  type = string
}

resource "null_resource" "example" {
  container_definitions = templatefile(
    var.filename,
    {
      foo = "bar"
    }
  )
}