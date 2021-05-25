package accurics

{{.prefix}}ecrmaketagsimmutable[con.id]{
    con = input.aws_ecr_repository[_]
    con.config.image_tag_mutability == "MUTABLE"
}

{{.prefix}}ecrmaketagsimmutable[con.id]{
    con = input.aws_ecr_repository[_]
    object.get(con.config, "image_tag_mutability", "undefined") == "undefined"
}