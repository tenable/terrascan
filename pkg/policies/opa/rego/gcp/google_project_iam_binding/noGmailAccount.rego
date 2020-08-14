package accurics

noGmailAccount[member.id] {
  member := input.google_project_iam_binding[_]
  mail := member.config.members[_]
  contains(mail, "gmail.com")
}
