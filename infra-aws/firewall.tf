resource "aws_security_group" "eks_cluster_sg" {
  depends_on = [aws_vpc.aws_vpc]
  vpc_id     = aws_vpc.aws_vpc.id
  tags = {
    Name = var.security_group["name"]
  }
  description = var.security_group["description"]
}

resource "aws_vpc_security_group_ingress_rule" "ingress_cluster_sg" {
  security_group_id = aws_security_group.eks_cluster_sg.id
  from_port         = var.sg_ingress_rule["from_port"]
  to_port           = var.sg_ingress_rule["to_port"]
  ip_protocol       = var.sg_ingress_rule["ip_protocol"]
  cidr_ipv4         = var.sg_ingress_rule["cidr_ipv4"]

}

resource "aws_vpc_security_group_egress_rule" "egress_cluster_sg" {
  security_group_id = aws_security_group.eks_cluster_sg.id
  ip_protocol       = var.sg_egress_rule["ip_protocol"]
  cidr_ipv4         = var.sg_egress_rule["cidr_ipv4"]
}
