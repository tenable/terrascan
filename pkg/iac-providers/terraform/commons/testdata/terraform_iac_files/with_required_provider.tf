terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 4.4"
    }
  }
}

resource "aws_ecs_task_definition" "demo-ecs-task-definition" {
  family                   = "ecs-task-definition-demo"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  memory                   = "1024"
  cpu                      = "512"
  execution_role_arn       = "arn:aws:iam::123456789012:role/ecsTaskExecutionRole"
  container_definitions    =  <<TASK_DEFINITION
[
    {
        "cpu": 10,
        "command": ["sleep", "10"],
        "entryPoint": ["/"],
        "environment": [
            {"name": "VARNAME", "value": "VARVAL"}
        ],
        "essential": true,
        "image": "jenkins",
        "memory": 128,
        "name": "jenkins",
        "portMappings": [
            {
                "containerPort": 80,
                "hostPort": 8080
            }
        ],
        "resourceRequirements":[
            {
                "type":"InferenceAccelerator",
                "value":"device_1"
            }
        ]
    }
]
TASK_DEFINITION

}