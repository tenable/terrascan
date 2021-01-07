resource "aws_ecs_task_definition" "instanceNotInVpc" {
  family                = "service"
  network_mode          = "bridge"
  container_definitions = file("ecs/service.json")
}
