module "security" {
  source = "./live/security"
  providers = {
    aws = aws.sec
  }
  # Set some vars
}

module "log" {
  depends_on = [module.security]
  source     = "./live/log"
  providers = {
    aws = aws.log
  }
  # Set some vars
}