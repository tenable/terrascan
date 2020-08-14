package accurics

{{.prefix}}macieIsNotAssociated[retVal] {
  check_empty(input)
  rc := "ZGF0YSAiYXdzX2NhbGxlcl9pZGVudGl0eSIgImN1cnJlbnQiIHt9CgpyZXNvdXJjZSAiYXdzX21hY2llX21lbWJlcl9hY2NvdW50X2Fzc29jaWF0aW9uIiAibWFjaWVfbWVtYmVyX2Fzc29jaWF0aW9uX25hbWUiIHsKICAgICJtZW1iZXJfYWNjb3VudF9pZCI6ICIke2RhdGEuYXdzX2NhbGxlcl9pZGVudGl0eS5jdXJyZW50LmFjY291bnRfaWR9Igp9"
  traverse = ""
  retVal := { "Id": "no_macie_member_account_association", "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "resource", "Expected": rc, "Actual": null }
}

check_empty(macie_input) = true {
	not macie_input.aws_macie_member_account_association
}

check_empty(macie_input) = true {
	count(macie_input.aws_macie_member_account_association) <= 0
}