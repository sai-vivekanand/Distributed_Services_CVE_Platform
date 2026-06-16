# Distributed Microservices Integrated with LLM

<img width="4880" height="3851" alt="image" src="https://github.com/user-attachments/assets/7c01dfd7-1284-46c9-b838-e22f378b4733" />


Cloud_Adv is a multi-component cloud platform repository that combines AWS infrastructure provisioning, Kubernetes operators, Helm charts, and an LLM-powered web application.

## Repository structure

- `/infra-aws`  
  Terraform code to provision AWS networking and EKS infrastructure, with Kubernetes/Helm resources for add-ons.
- `/ami-jenkins`  
  Packer + shell scripts to build a Jenkins-ready AMI, plus Jenkins job definitions and CI workflows.
- `/cve-operator`  
  Go-based Kubernetes operator (Kubebuilder) with CRDs for monitoring GitHub releases and creating processing jobs.
- `/helm-cve-operator`  
  Helm chart to deploy the CVE operator and its supporting Kubernetes resources.
- `/webapp-llm`  
  Python (Flask + Streamlit + LangChain) application that performs retrieval-augmented responses using Pinecone and Ollama.
- `/helm-webapp-llm`  
  Helm chart to deploy the LLM web application on Kubernetes.

## High-level architecture

1. AWS infrastructure is created via Terraform in `infra-aws` (VPC, IAM, EKS, and supporting resources).
2. Jenkins automation is prepared via `ami-jenkins` and used to run build/deploy pipelines.
3. `cve-operator` monitors GitHub release feeds and creates Kubernetes Jobs for matching release assets.
4. Helm charts (`helm-cve-operator`, `helm-webapp-llm`) package and deploy workloads to Kubernetes.
5. `webapp-llm` exposes API/UI endpoints for LLM-backed question answering using vector retrieval.

## Core components

### CVE Operator (`/cve-operator`)

- Defines two CRDs:
  - `GithubReleaseMonitor`: polls a GitHub releases API URL and tracks processing status.
  - `GithubRelease`: represents a release asset to process and triggers a Kubernetes Job.
- Includes controller logic to:
  - fetch release metadata from GitHub API,
  - filter release assets (e.g., `_delta_` artifacts),
  - create/recreate processing resources,
  - track job lifecycle in CR status.

### LLM Web App (`/webapp-llm`)

- Flask routes include:
  - `POST /generate` for prompt-based response generation,
  - `GET /healthz` for health checks.
- Uses:
  - Pinecone as vector store,
  - Hugging Face embeddings (`all-MiniLM-L6-v2`),
  - Ollama as the LLM backend,
  - Streamlit for UI.
- Runtime is containerized via Docker (`Dockerfile`, `startup.sh`).

## Getting started

Because this repository is split into multiple deployable modules, start in the component you want:

- Infrastructure: `/infra-aws/README.md`
- Jenkins AMI: `/ami-jenkins/README.md`
- Operator development/deployment: `/cve-operator/README.md`
- Web app runtime: `/webapp-llm/README.md`

For Kubernetes deployment, use Helm charts in:

- `/helm-cve-operator`
- `/helm-webapp-llm`

## CI/CD and automation

- Jenkinsfiles exist in multiple modules for pipeline execution.
- `.github/workflows` in `ami-jenkins` contains packer validation/build workflows.

## Notes

- Existing module READMEs are concise/incomplete in places; this root README provides a cross-repository overview.
- See each module directory for exact environment variables, credentials, and deployment specifics.
