package accurics

{{.prefix}}s3EnforceUserACL[retVal] {
    bucket := input.aws_s3_bucket[_]

    bucket_policies_set := { policy_id | policy_id := split(input.aws_s3_bucket_policy[_].id, "." )[1] }

    not bucket_policies_set[split(bucket.id, ".")[1]]

	rc = "cmVzb3VyY2UgImF3c19zM19idWNrZXRfcG9saWN5IiAiIyNyZXNvdXJjZV9uYW1lIyNQb2xpY3kiIHsKICBidWNrZXQgPSAiJHthd3NfczNfYnVja2V0LiMjcmVzb3VyY2VfbmFtZSMjLmlkfSIKCiAgcG9saWN5ID0gPDxQT0xJQ1kKewogICJWZXJzaW9uIjogIjIwMTItMTAtMTciLAogICJTdGF0ZW1lbnQiOiBbCiAgICB7CiAgICAgICJTaWQiOiAiIyNyZXNvdXJjZV9uYW1lIyMtcmVzdHJpY3QtYWNjZXNzLXRvLXVzZXJzLW9yLXJvbGVzIiwKICAgICAgIkVmZmVjdCI6ICJBbGxvdyIsCiAgICAgICJQcmluY2lwYWwiOiBbCiAgICAgICAgewogICAgICAgICAgIkFXUyI6IFsKICAgICAgICAgICAgImFybjphd3M6aWFtOjojI2Fjb3VudF9pZCMjOnJvbGUvIyNyb2xlX25hbWUjIyIsCiAgICAgICAgICAgICJhcm46YXdzOmlhbTo6IyNhY291bnRfaWQjIzp1c2VyLyMjdXNlcl9uYW1lIyMiCiAgICAgICAgICBdCiAgICAgICAgfQogICAgICBdLAogICAgICAiQWN0aW9uIjogInMzOkdldE9iamVjdCIsCiAgICAgICJSZXNvdXJjZSI6ICJhcm46YXdzOnMzOjo6JHthd3NfczNfYnVja2V0LiMjcmVzb3VyY2VfbmFtZSMjLmlkfS8qIgogICAgfQogIF0KfQpQT0xJQ1kKfQ=="
    decode_rc = base64.decode(rc)
    replaced_resource_name := replace(decode_rc, "##resource_name##", bucket.name)

    traverse = ""
    retVal := { "Id": bucket.id, "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "resource", "Expected": base64.encode(replaced_resource_name), "Actual": null }
}