package accurics

weakCipherSuitesEnabled[api.id]
{
    api := input.google_compute_ssl_policy[_]
    not api.config.min_tls_version == "TLS_1_2"
}