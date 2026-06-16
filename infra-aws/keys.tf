data "aws_caller_identity" "current" {}

# resource "aws_kms_key" "key_for_eks_cluster" {
#   description             = "Symmetric KMS key for EKS cluster"
#   depends_on              = [aws_iam_role.eks_cluster_role]
#   enable_key_rotation     = true
#   key_usage               = "ENCRYPT_DECRYPT"
#   deletion_window_in_days = 20
#   policy = jsonencode({
#     Version = "2012-10-17"
#     Statement = [
#       {
#         Sid    = "Enable IAM User Permissions"
#         Effect = "Allow"
#         Principal = {
#           AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
#         },
#         Action   = "kms:*"
#         Resource = "*"
#       },
#       {
#         Sid    = "Key Administrators"
#         Effect = "Allow"
#         Principal = {
#           AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
#         },
#         Action = [
#           "kms:Update*",
#           "kms:UntagResource",
#           "kms:TagResource",
#           "kms:ScheduleKeyDeletion",
#           "kms:Revoke*",
#           "kms:ReplicateKey",
#           "kms:Put*",
#           "kms:List*",
#           "kms:ImportKeyMaterial",
#           "kms:Get*",
#           "kms:Enable*",
#           "kms:Disable*",
#           "kms:Describe*",
#           "kms:Delete*",
#           "kms:Create*",
#           "kms:CancelKeyDeletion"

#         ],
#         Resource = "*"
#       },
#       {
#         Sid    = "KeyUsage"
#         Effect = "Allow"
#         Principal = {
#           AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${aws_iam_role.eks_cluster_role.name}"
#         },
#         Action = [
#           "kms:ReEncrypt*",
#           "kms:GenerateDataKey*",
#           "kms:Encrypt",
#           "kms:DescribeKey",
#           "kms:Decrypt"
#         ],
#         Resource = "*"
#       }
#     ]
#   })
# }

# resource "aws_kms_alias" "cluster_key_alias" {
#   name          = "alias/eks-cluster-key"
#   target_key_id = aws_kms_key.key_for_eks_cluster.id

# }


resource "aws_kms_key" "key_for_ebs_volume" {
  description             = var.key_for_ebs_volume.description
  depends_on              = [module.ebs_csi_irsa_role]
  enable_key_rotation     = var.key_for_ebs_volume.enable_key_rotation
  key_usage               = var.key_for_ebs_volume.key_usage
  deletion_window_in_days = var.key_for_ebs_volume.deletion_window_in_days
  policy = jsonencode({
    Version = var.key_for_ebs_volume.policy.Version
    Statement = [
      {
        Sid    = var.key_for_ebs_volume.policy.Statement[0].Sid
        Effect = var.key_for_ebs_volume.policy.Statement[0].Effect
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
        },
        Action   = var.key_for_ebs_volume.policy.Statement[0].Action
        Resource = var.key_for_ebs_volume.policy.Statement[0].Resource
      },
      {
        Sid    = var.key_for_ebs_volume.policy.Statement[1].Sid
        Effect = var.key_for_ebs_volume.policy.Statement[1].Effect
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
        },
        Action   = var.key_for_ebs_volume.policy.Statement[1].Action,
        Resource = var.key_for_ebs_volume.policy.Statement[1].Resource
      },
      {
        Sid    = var.key_for_ebs_volume.policy.Statement[2].Sid
        Effect = var.key_for_ebs_volume.policy.Statement[2].Effect
        Principal = {
          # AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${module.ebs_csi_irsa_role.iam_role_name}",

          AWS = [
            "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${aws_iam_role.eks_node_role.name}",
            "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${module.ebs_csi_irsa_role.iam_role_name}"
          ]
        },
        Action   = var.key_for_ebs_volume.policy.Statement[2].Action,
        Resource = var.key_for_ebs_volume.policy.Statement[2].Resource
      }
    ]
  })
}

resource "aws_kms_alias" "cluster_node_alias" {
  name          = "alias/ebs-volume-key"
  target_key_id = aws_kms_key.key_for_ebs_volume.id
}

output "key_id_ebs_volume_arn" {
  value = aws_kms_key.key_for_ebs_volume.arn

}
