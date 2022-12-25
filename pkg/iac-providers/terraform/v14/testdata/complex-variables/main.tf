resource "complex_var_resource" "complex" {
  stringList = var.stringList
  intList = var.intList
  floatList = var.floatList
  boolList = var.boolList
  list_no_type = var.list_no_type
  setVar = var.set_var
  tupleVar = var.tuple_var
  listTuple = var.list_of_tuple
  tupleVarComplex = var.tuple_var_complex
  objectVar = var.object_var
  mapVar = var.map_var
  mapVarComplex = var.map_var_complex
  objectList = var.object_list
  objectListComplex = var.object_list_complex
}