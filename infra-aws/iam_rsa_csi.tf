module "ebs_csi_irsa_role" {
  source                = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"
  attach_ebs_csi_policy = var.ebs_csi_irsa_config["attach_ebs_csi_policy"]
  role_policy_arns      = var.ebs_csi_irsa_config["role_policy_arns"]
  role_name             = var.ebs_csi_irsa_config["role_name"]
  oidc_providers = {
    example = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = var.ebs_csi_irsa_config["oidc_providers"]["example"]["namespace_service_accounts"]
    }
  }
}
