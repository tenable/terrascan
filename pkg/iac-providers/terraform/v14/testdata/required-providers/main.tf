terraform {
    required_providers {
        newrelic = {
        # source is required for providers in other namespaces, to avoid ambiguity.
        source  = "newrelic/newrelic"
        version = "~> 2.1.1"
        }
    }
}