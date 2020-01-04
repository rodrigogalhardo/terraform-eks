variable "namespace_name" {
  default = "batpay-eks-service"
  type    = "string"
}

variable "nginx_pod_name" {
  default = "batpay-eks-service"
  type    = "string"
}

variable "nginx_pod_image" {
  default = "nginx:latest"
  type    = "string"
}
