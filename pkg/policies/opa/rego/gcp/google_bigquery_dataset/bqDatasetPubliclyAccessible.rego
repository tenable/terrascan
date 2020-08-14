package accurics

bqDatasetPubliclyAccessible[api.id]{
    api := input.google_bigquery_dataset[_]
    data := api.config.access[_]
    data.special_group == "allAuthenticatedUsers"
}