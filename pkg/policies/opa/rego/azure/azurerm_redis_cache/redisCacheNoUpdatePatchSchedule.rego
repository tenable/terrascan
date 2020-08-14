package accurics

{{.prefix}}redisCacheNoUpdatePatchSchedule[retVal] {
    redis := input.azurerm_redis_cache[_]
    count(redis.config.patch_schedule) <= 0
    emptyPatchSchedule(redis) == true
   	rc = "ewogICJwYXRjaF9zY2hlZHVsZSI6IHsKICAgICJkYXlfb2Zfd2VlayI6ICJTdW5kYXkiLAogICAgInN0YXJ0X2hvdXJfdXRjIjogMAogIH0KfQ=="
    traverse = "patch_schedule"
    retVal := { "Id": redis.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "base64", "Expected": rc, "Actual": redis.config.patch_schedule }
}

emptyPatchSchedule(redis) = true {
    not redis.config.patch_schedule
}
emptyPatchSchedule(redis) = true {
    count(redis.config.patch_schedule) <= 0
}