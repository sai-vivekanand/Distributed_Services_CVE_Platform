provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster_data.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster_data.certificate_authority.0.data)
    token                  = data.aws_eks_cluster_auth.cluster_auth.token
  }
  alias = "helm-eks"
}

resource "helm_release" "fluentbit" {
  name       = "fluentbit"
  depends_on = [aws_iam_role.eks_node_role, kubernetes_namespace.amazon_cloudwatch, module.eks]
  repository = "./helm-charts/"
  provider   = helm.helm-eks
  chart      = "fluentbit"
  namespace  = kubernetes_namespace.amazon_cloudwatch.metadata[0].name
}

resource "null_resource" "download_asset" {
  provisioner "local-exec" {
    command = <<EOT
#!/bin/bash

# Variables passed from Terraform
GITHUB_TOKEN="${var.github_pat}"
ASSET_URL="${var.autoscaler_config.asset_url}"
OUTPUT_FILE="${path.module}/asset.tgz"

# Download the asset using the asset URL
curl -L \
  -H "Accept: application/octet-stream" \
  -H "Authorization: Bearer $GITHUB_TOKEN" \
  -o $OUTPUT_FILE \
  $ASSET_URL

echo "Downloaded the asset to $OUTPUT_FILE"
EOT
  }
}

resource "helm_release" "kubernetes-autoscaler" {
  name       = "kubernetes-autoscaler"
  provider   = helm.helm-eks
  depends_on = [aws_iam_role.eks_autoscaler_role, kubernetes_namespace.namespace_autoscaler, module.eks, kubernetes_secret.dockerhub_secret, null_resource.download_asset]
  chart      = "${path.module}/asset.tgz"
  namespace  = var.namespace_autoscaler
  values     = [file("manifests/autoscaler-values.yaml")]
  set {
    name  = "autoDiscovery.clusterName"
    value = var.kubernetes_autoscaler.cluster_name
  }

  set {
    name  = "awsRegion"
    value = var.kubernetes_autoscaler.aws_region
  }

  set {
    name  = "rbac.serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
    value = aws_iam_role.eks_autoscaler_role.arn
  }

  set {
    name  = "rbac.serviceAccount.name"
    value = var.kubernetes_autoscaler.service_account_name
  }

  set {
    name  = "image.repository"
    value = var.autoscaler_config.image.repository
  }
  set {
    name  = "image.tag"
    value = var.autoscaler_config.image.tag
  }

  set {
    name  = "image.pullSecrets[0]"
    value = kubernetes_secret.dockerhub_secret.metadata[0].name
  }
}

resource "helm_release" "istio_base" {
  name       = "istio-base"
  repository = "https://istio-release.storage.googleapis.com/charts"
  provider   = helm.helm-eks
  chart      = "base"
  namespace  = kubernetes_namespace.namespace_istio.metadata[0].name
  depends_on = [
    kubernetes_namespace.namespace1, kubernetes_namespace.namespace2, kubernetes_namespace.namespace3,
    helm_release.kubernetes-autoscaler
  ]
  set {
    name  = "global.istioNamespace"
    value = kubernetes_namespace.namespace_istio.metadata[0].name
  }
}

resource "helm_release" "istio_daemon" {
  name       = "istio-daemon"
  repository = "https://istio-release.storage.googleapis.com/charts"
  provider   = helm.helm-eks
  chart      = "istiod"
  namespace  = kubernetes_namespace.namespace_istio.metadata[0].name
  depends_on = [
    kubernetes_namespace.namespace1, kubernetes_namespace.namespace2, kubernetes_namespace.namespace3,
    helm_release.kubernetes-autoscaler, helm_release.istio_base
  ]
  set {
    name  = "meshConfig.ingressService"
    value = "istio-gateway"
  }

  set {
    name  = "meshConfig.ingressSelector.app"
    value = "istio-gateway"
  }

  # set{
  #   name = "meshConfig.enablePrometheusMerge"
  #   value = false   
  # }

  # set{
  #   name = "meshConfig.defaultConfig.holdApplicationUntilProxyStarts"
  #   value = true
  # }

  set {
    name  = "global.proxy.holdApplicationUntilProxyStarts"
    value = "true"
  }

  set {
    name  = "global.logAsJson"
    value = "true"
  }
}

resource "helm_release" "istio_gateway" {
  name = "istio-gateway" # should be istio-ingressgateway for defaults of istiod 
  # name should be the same what you for kubernetes_service.istio_gateway_service_data 
  repository = "https://istio-release.storage.googleapis.com/charts"
  provider   = helm.helm-eks
  chart      = "gateway"
  namespace  = kubernetes_namespace.namespace_istio.metadata[0].name
  depends_on = [
    kubernetes_namespace.namespace1, kubernetes_namespace.namespace2, kubernetes_namespace.namespace3,
    helm_release.kubernetes-autoscaler, helm_release.istio_base, helm_release.istio_daemon
  ]

  values = ["${file("manifests/helm-istio-gateway-values.yaml")}"]
  set {
    name  = "labels.app"
    value = "istio-gateway"
  }

}


resource "helm_release" "postgresql-ha-release" {
  name       = "postgresql-ha-release"
  repository = "https://charts.bitnami.com/bitnami"
  provider   = helm.helm-eks
  chart      = "postgresql-ha"
  depends_on = [kubernetes_namespace.namespace1, kubernetes_namespace.namespace2, kubernetes_namespace.namespace3,
    helm_release.kubernetes-autoscaler, kubernetes_storage_class.pvc_sc
  ]
  namespace = var.postgres_ha.namespace
  values = [
    "${file("manifests/postgres-values.yaml")}"
  ]
}

resource "helm_release" "kafka" {
  name       = "kafka"
  repository = "https://charts.bitnami.com/bitnami"
  provider   = helm.helm-eks
  chart      = "kafka"
  depends_on = [kubernetes_namespace.namespace1, kubernetes_namespace.namespace2, kubernetes_namespace.namespace3,
    helm_release.kubernetes-autoscaler, helm_release.istio_base, helm_release.istio_daemon, helm_release.istio_gateway
  ]
  namespace = var.kafka_config.namespace
  values = [
    "${file("manifests/kafka-values.yaml")}"
  ]
}

resource "helm_release" "prometheus_grafana" {
  name       = "grafana-prometheus"
  repository = "https://prometheus-community.github.io/helm-charts"
  provider   = helm.helm-eks
  chart      = "kube-prometheus-stack"
  namespace  = kubernetes_namespace.prometheus_grafana_ns.metadata[0].name
  depends_on = [
    kubernetes_namespace.namespace1, kubernetes_namespace.namespace2, kubernetes_namespace.namespace3,
    helm_release.kubernetes-autoscaler, helm_release.istio_base, helm_release.istio_daemon, helm_release.istio_gateway
  ]
  values = [
    "${file("manifests/prometheus-values.yaml")}"
  ]
}

resource "helm_release" "cert_manager" {
  name       = "cert-manager"
  depends_on = [kubernetes_namespace.cert_manager, module.eks]
  repository = "https://charts.jetstack.io"
  provider   = helm.helm-eks
  chart      = "cert-manager"
  # namespace  = kubernetes_namespace.prometheus_grafana_ns.metadata[0].name
  namespace = kubernetes_namespace.namespace_istio.metadata[0].name
  values = [
    "${file("manifests/cert-manager-values.yaml")}"
  ]
}

resource "helm_release" "external_dns" {
  name       = "external-dns"
  depends_on = [kubernetes_namespace.prometheus_grafana_ns, module.eks]
  repository = "https://kubernetes-sigs.github.io/external-dns/"
  provider   = helm.helm-eks
  chart      = "external-dns"
  # version    = "1.14.5"
  namespace = kubernetes_namespace.namespace_istio.metadata[0].name
  values = [
    "${file("manifests/external-dns-values.yaml")}"
  ]

}

resource "helm_release" "ollam_release" {
  name       = "ollam-release"
  repository = "https://otwld.github.io/ollama-helm/"
  provider   = helm.helm-eks
  chart      = "ollama"
  namespace  = kubernetes_namespace.ollama.metadata[0].name
  depends_on = [
    kubernetes_namespace.ollama, helm_release.kubernetes-autoscaler, helm_release.istio_base
  ]
  values = [
    "${file("manifests/ollama-helm.yaml")}"
  ]
}
