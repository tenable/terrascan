resource "aws_iam_access_key" "noAccessKeyForRootAccount" {
  user    = "root"
  pgp_key = "keybase:some_person_that_exists"
  status  = "Inactive"
}
