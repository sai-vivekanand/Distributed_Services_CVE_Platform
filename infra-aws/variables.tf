variable "cluster_eks" {
  type        = any
  description = "EKS Cluster Configuration"

}

variable "security_group" {
  type        = any
  description = "Security Group Configuration"
}

variable "sg_ingress_rule" {
  type        = any
  description = "Security Group Ingress Rule Configuration"
}

variable "sg_egress_rule" {
  type        = any
  description = "Security Group Egress Rule Configuration"
}

variable "ebs_csi_irsa_config" {
  type        = any
  description = "IRSA EBS Configuration"
}

variable "key_for_ebs_volume" {
  type        = any
  description = "KMS Key for EBS Volume"
}

variable "aws_vpc" {
  type        = any
  description = "VPC Configuration"
}

variable "aws_public_subnet1" {
  type        = any
  description = "Public Subnet 1 Configuration"
}

variable "aws_public_subnet2" {
  type        = any
  description = "Public Subnet 2 Configuration"
}

variable "aws_public_subnet3" {
  type        = any
  description = "Public Subnet 3 Configuration"
}

variable "aws_private_subnet1" {
  type        = any
  description = "Private Subnet 1 Configuration"
}


variable "aws_internet_gateway" {
  type        = any
  description = "Internet Gateway Configuration"
}

variable "aws_route_table_public_subnets" {
  type        = any
  description = "Route Table for Public Subnets Configuration"

}

variable "aws_nat_eip_private_subnet1" {
  type        = any
  description = "EIP for NAT Gateway Configuration"
}

variable "aws_nat_private_subnet1" {
  type        = any
  description = "NAT Gateway for Private Subnet 1 Configuration"
}

variable "aws_route_table_private_subnet1" {
  type        = any
  description = "Route Table for Private Subnet 1 Configuration"
}

variable "aws_private_subnet2" {
  type        = any
  description = "Private Subnet 2 Configuration"
}

variable "aws_nat_eip_private_subnet2" {
  type        = any
  description = "EIP for NAT Gateway Configuration"
}

variable "aws_nat_private_subnet2" {
  type        = any
  description = "NAT Gateway for Private Subnet 2 Configuration"
}

variable "aws_route_table_private_subnet2" {
  type        = any
  description = "Route Table for Private Subnet 2 Configuration"

}
variable "aws_private_subnet3" {
  type        = any
  description = "Private Subnet 3 Configuration"
}

variable "aws_nat_eip_private_subnet3" {
  type        = any
  description = "EIP for NAT Gateway Configuration"
}

variable "aws_nat_private_subnet3" {
  type        = any
  description = "NAT Gateway for Private Subnet 3 Configuration"
}

variable "aws_route_table_private_subnet3" {
  type        = any
  description = "Route Table for Private Subnet 3 Configuration"
}

variable "eks_cluster_role" {
  type        = any
  description = "EKS Cluster Role Configuration"

}

variable "eks_node_role" {
  type        = any
  description = "EKS Node Role Configuration"
}

variable "aws_iam_policy_document_eks" {
  type        = any
  description = "IAM Policy Document Configuration"

}

variable "aws_iam_policy_document_eks_node" {
  type        = any
  description = "IAM Policy Document Configuration"

}

variable "postgres_ha" {
  type        = any
  description = "Postgres HA Configuration"
}

variable "kafka_config" {
  type        = any
  description = "Kafka Configuration"
}

variable "namespace1" {
  type        = string
  description = "Namespace 1 Configuration"
}

variable "namespace2" {
  type        = string
  description = "Namespace 1 Configuration"
}

variable "namespace3" {
  type        = string
  description = "Namespace 1 Configuration"
}

variable "namespace_autoscaler" {
  type        = string
  description = "Namespace Autoscaler Configuration"
}

variable "sc_config" {
  type        = any
  description = "Storage Class Configuration"

}

variable "kubernetes_autoscaler" {
  type        = any
  description = "Kubernetes Autoscaler Configuration"

}


variable "eks_autoscaler_policy_doc" {
  type        = any
  description = "EKS Autoscaler Policy Document Configuration"

}

variable "eks_autoscaler_policy" {
  type        = any
  description = "EKS Autoscaler Policy Configuration"
}

variable "github_username" {
  type        = string
  description = "Github Username"
}
variable "github_pat" {
  type        = string
  description = "Github Personal Access"
}

variable "github_chart_url" {
  type        = string
  description = "Github Chart URL"
}

variable "docker_hub_registry" {
  type        = string
  description = "Docker Hub Registry"
}

variable "docker_hub_username" {
  type        = string
  description = "Docker Hub Username"

}

variable "docker_hub_password" {
  type        = string
  description = "Docker Hub Password"
}

variable "docker_hub_email" {
  type        = string
  description = "Docker Hub Email"
}

variable "autoscaler_config" {
  type        = any
  description = "Autoscaler Configuration"

}
variable "namespaces" {
  type        = any
  description = "Namespaces Configuration"
}

variable "namespace_istio_system" {
  type        = string
  description = "Namespace Istio System"
}

variable "sa_cert_manager" {
  type        = string
  description = "Service account for cert manager"
}

variable "sa_external_dns" {
  type        = string
  description = "Service account for external dns"
}
