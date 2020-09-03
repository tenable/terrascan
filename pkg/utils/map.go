package utils

// InterfaceToMapStringInterface converts instances of map[interface{}interface{} to map[string]interface{} within a nested
// map data structure such as json
func InterfaceToMapStringInterface(iface interface{}) interface{} {
	switch ifaceObj := iface.(type) {
	case []interface{}:
		for i := range ifaceObj {
			ifaceObj[i] = InterfaceToMapStringInterface(ifaceObj[i])
		}
		return iface
	case map[interface{}]interface{}:
		mapData := make(map[string]interface{})
		for k := range ifaceObj {
			mapData[k.(string)] = InterfaceToMapStringInterface(ifaceObj[k])
		}
		return mapData
	}
	return iface
}
