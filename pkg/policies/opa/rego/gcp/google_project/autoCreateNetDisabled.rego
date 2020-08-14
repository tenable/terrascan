package accurics

autoCreateNetDisabled[api.id]
{
     api := input.google_project[_]
     not api.config.auto_create_network == false
}

