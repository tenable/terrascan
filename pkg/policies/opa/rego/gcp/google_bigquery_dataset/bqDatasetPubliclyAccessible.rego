package accurics

bqDatasetPubliclyAccessible[api.id]{
    api := input.google_bigquery_dataset[_]
    access := api.config.access[_]
    access.special_group == "allAuthenticatedUsers"
}