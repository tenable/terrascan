resource "aws_load_balancer_policy" "elbWeakCipher" {
  load_balancer_name = "some-name"
  policy_name        = "wu-tang-ssl"
  policy_type_name   = "SSLNegotiationPolicyType"

  policy_attribute {
    name  = "ECDHE-RSA-RC4-SHA"
    value = "true"
  }
}
