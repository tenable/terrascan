package accurics

caPoolMaxPathLenZero[capools.id] {
  capools := input.google_privateca_ca_pool[_]
  issuancepolicy := capools.config.issuance_policy[_]
  baselinevalues := issuancepolicy.baseline_values[_]
  caoptions := baselinevalues.ca_options[_]
  caoptions.max_issuer_path_length == 0
  not caoptions.zero_max_issuer_path_length == true
  not caoptions.zero_max_issuer_path_length == false
}
