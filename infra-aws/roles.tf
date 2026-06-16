data "aws_iam_policy_document" "eks_cluster_assume_role_policy" {
  statement {
    actions = var.aws_iam_policy_document_eks["statement"]["actions"]
    effect  = var.aws_iam_policy_document_eks["statement"]["effect"]
    principals {
      type        = var.aws_iam_policy_document_eks["statement"]["principals"]["type"]
      identifiers = var.aws_iam_policy_document_eks["statement"]["principals"]["identifiers"]
    }
  }
}

data "aws_iam_policy_document" "eks_node_assume_role_policy" {
  statement {
    actions = var.aws_iam_policy_document_eks_node["statement"]["actions"]
    effect  = var.aws_iam_policy_document_eks_node["statement"]["effect"]
    principals {
      type        = var.aws_iam_policy_document_eks_node["statement"]["principals"]["type"]
      identifiers = var.aws_iam_policy_document_eks_node["statement"]["principals"]["identifiers"]
    }
  }
}

data "aws_iam_policy_document" "eks_autoscaler_policy_doc" {
  statement {
    actions   = var.eks_autoscaler_policy_doc["statement"]["actions"]
    effect    = var.eks_autoscaler_policy_doc["statement"]["effect"]
    resources = var.eks_autoscaler_policy_doc["statement"]["resources"]
  }
}

resource "aws_iam_policy" "eks_autoscaler_policy" {
  name        = var.eks_autoscaler_policy["name"]
  description = var.eks_autoscaler_policy["description"]
  policy      = data.aws_iam_policy_document.eks_autoscaler_policy_doc.json
}

resource "aws_iam_role" "eks_cluster_role" {
  name                = var.eks_cluster_role["name"]
  depends_on          = [data.aws_iam_policy_document.eks_cluster_assume_role_policy]
  managed_policy_arns = var.eks_cluster_role["managed_policy_arns"]
  assume_role_policy  = data.aws_iam_policy_document.eks_cluster_assume_role_policy.json
}

resource "aws_iam_role" "eks_node_role" {
  name                = var.eks_node_role["name"]
  depends_on          = [data.aws_iam_policy_document.eks_cluster_assume_role_policy]
  managed_policy_arns = var.eks_node_role["managed_policy_arns"]
  assume_role_policy  = data.aws_iam_policy_document.eks_node_assume_role_policy.json
}


resource "aws_iam_role" "eks_autoscaler_role" {
  name = "eks-autoscaler-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Federated = "${module.eks.oidc_provider_arn}"
        },
        Action = "sts:AssumeRoleWithWebIdentity",
        Condition = {
          StringEquals = {
            "${module.eks.oidc_provider}:sub" = "system:serviceaccount:${var.namespace_autoscaler}:${var.kubernetes_autoscaler.service_account_name}"
          }
        }
      }
    ]
  })
  managed_policy_arns = [aws_iam_policy.eks_autoscaler_policy.arn]
}

resource "aws_iam_policy" "route53_policy" {
  name        = "route53_policy"
  description = "IAM policy for Route53 permissions"
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : "route53:GetChange",
        "Resource" : "arn:aws:route53:::change/*"
      },
      {
        "Effect" : "Allow",
        "Action" : [
          "route53:ChangeResourceRecordSets",
          "route53:ListResourceRecordSets"
        ],
        "Resource" : "arn:aws:route53:::hostedzone/*"
      },
      {
        "Effect" : "Allow",
        "Action" : ["route53:ListHostedZonesByName",
          "route53:ListHostedZones",
          "route53:GetHostedZone",
          "route53:ChangeResourceRecordSets",
          "route53:ListResourceRecordSets",
          "route53:GetChange",
          "route53:ListTagsForResource"

        ],
        "Resource" : "*"
      }
    ]
  })
}

resource "aws_iam_role" "eks_route53_role" {
  name = "eks_route53_role"

  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : "sts:AssumeRoleWithWebIdentity",
        "Principal" : {
          "Federated" : "${module.eks.oidc_provider_arn}"
        },
        "Condition" : {
          "StringLike" : {
            "${module.eks.oidc_provider}:sub" : [
              "system:serviceaccount:${var.namespace_istio_system}:${var.sa_cert_manager}",
              "system:serviceaccount:${var.namespace_istio_system}:${var.sa_external_dns}",
            ]
          }
        }
      }
    ]
  })
}

# Attach the Route53 policy to the role
resource "aws_iam_role_policy_attachment" "route53_policy_attachment" {
  depends_on = [aws_iam_role.eks_route53_role, aws_iam_policy.route53_policy]
  role       = aws_iam_role.eks_route53_role.name
  policy_arn = aws_iam_policy.route53_policy.arn
}
