package accurics 

{{.prefix}}cloudfrontNoSecureCiphers[retVal]{
    cloudfront = input.aws_cloudfront_distribution[_]
    some i
    certificate = cloudfront.config.viewer_certificate[i]
    certificate.cloudfront_default_certificate = false
    not minimumAllowedProtocolVersion(certificate.minimum_protocol_version)
    traverse := sprintf("viewer_certificate[%d].minimum_protocol_version", [i])
    retVal := { "Id": cloudfront.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "viewer_certificate.minimum_protocol_version", "AttributeDataType": "string", "Expected": "TLSv1.2", "Actual": certificate.minimum_protocol_version }
}

minimumAllowedProtocolVersion(currentVersion) {
    currentVersion == "TLSv1.1"
}

minimumAllowedProtocolVersion(currentVersion) {
    currentVersion == "TLSv1.2"
}