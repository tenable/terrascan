resource "datadog_monitor" "foo" {
  tags = ["baz:biz"]
}

resource "datadog_monitor" "bar" {
  tags = ["baz:biz", "foo:bar"]
}
