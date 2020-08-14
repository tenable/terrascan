resource "type1" "resource1" {
	test3 = 1 + 2
	test1 = "hello"
	test2 = 5
	arr = [1, 2, 3, 4]
	hyphen-test = 3
	temp = "${1 + 2} %{if local.test2 < 3}\"4\n\"%{endif}"
	temp2 = "${"hi"} there"
		quoted = "\"quoted\""
		squoted = "'quoted'"
	}

resource "type2" "resource2"{
	other = {
		num = local.test2 + 5
		thing = [for x in local.arr: x * 2]
		"${local.test3}" = 4
		3 = 1
		"local.test1" = 89
		"a.b.c[\"hi\"][3].*" = 3
		loop = "This has a for loop: %{for x in local.arr}x,%{endfor}"
		a.b.c = "True"
	}
}

resource "type3" "resource3" {
	heredoc = <<-EOF
		This is a heredoc template.
		It references ${local.other.3}
	EOF
	simple = 4 - 2
	cond = test3 > 2 ? 1: 0
	heredoc2 = <<EOF
		Another heredoc, that
		doesn't remove indentation
		${local.other.3}
		%{if true ? false : true}"gotcha"\n%{else}4%{endif}
	EOF
}

resource "terraform_remote_state" "remote" {
	backend = "s3"
	config = {
		profile = var.profile
		region  = var.region
		bucket  = "mybucket"
		key     = "mykey"
	}
}

variable "profile" {}

variable "region" {
	default = "us-east-1"
}
