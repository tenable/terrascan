package accurics

{{.prefix}}athenaDatabaseEncrypted[athena.id]{
    athena = input.aws_athena_database[_]
	object.get(athena.config, "encryption_configuration", "undefined") = ["undefined", []][_]
}