# Distributed Microservices Integrated with LLM

<img width="4880" height="3851" alt="image" src="https://github.com/user-attachments/assets/7c01dfd7-1284-46c9-b838-e22f378b4733" />


About

This is a robust data streaming system that streams over 240000+ records of the CVE (Common Vulnerabilities and Exposures) data, through 𝗞𝗮𝗳𝗸𝗮 to store it in 𝗣𝗼𝘀𝘁𝗴𝗿𝗲𝗦𝗤𝗟 database. The stored CVE data is then utilized as part of the RAG (𝗥𝗲𝘁𝗿𝗶𝗲𝘃𝗮𝗹-𝗔𝘂𝗴𝗺𝗲𝗻𝘁𝗲𝗱 𝗚𝗲𝗻𝗲𝗿𝗮𝘁𝗶𝗼𝗻) pipeline to power a self-hosted LLM which can answer queries specific to the CVE data. The whole infrastructure is built on top of 𝗞𝘂𝗯𝗲𝗿𝗻𝗲𝘁𝗲𝘀, utilizing 𝗔𝗺𝗮𝘇𝗼𝗻 𝗘𝗞𝗦.

Here are some of the key aspects of our architecture:

    CI/CD - Implemented 𝗝𝗲𝗻𝗸𝗶𝗻𝘀 which monitors multiple repositories to lint the code, automate the 𝗱𝗼𝗰𝗸𝗲𝗿 𝗶𝗺𝗮𝗴𝗲 builds and handles semantic releases of the helm packages.

    Microservices - The setup includes multiple microservices, built using 𝗚𝗼𝗟𝗮𝗻𝗴.

    Scalability - Implemented pod level 𝗮𝘂𝘁𝗼 𝘀𝗰𝗮𝗹𝗲𝗿 to individually scale each microservice in the architecture and implemented 𝗰𝗹𝘂𝘀𝘁𝗲𝗿 𝗹𝗲𝘃𝗲𝗹 auto scaler to scale at kubernetes node level.

    Availability - Placed the pods across different AWS availability zones to ensure 𝗳𝗮𝘂𝗹𝘁 𝘁𝗼𝗹𝗲𝗿𝗮𝗻𝗰𝗲.

    Reliability - Implemented 𝗹𝗶𝘃𝗲𝗻𝗲𝘀𝘀 and 𝗿𝗲𝗮𝗱𝗶𝗻𝗲𝘀𝘀 𝗽𝗿𝗼𝗯𝗲𝘀 to keep the application resilient against failures.

    Security - Followed 𝗣𝗼𝗟𝗣 (Principle of Least Privilege) and used appropriate IAM roles, trust relationships, k8s secrets, service accounts and RBAC. Provided secure access with 𝗦𝗦𝗟 𝗰𝗲𝗿𝘁𝗶𝗳𝗶𝗰𝗮𝘁𝗲𝘀 managed through k8s cert-manager operator.

    Custom Kubernetes operator - Developed a 𝗰𝘂𝘀𝘁𝗼𝗺 𝗸𝟴𝘀 𝗼𝗽𝗲𝗿𝗮𝘁𝗼𝗿 to monitor the latest CVE releases and to keep the cve postgres database updated automatically with the latest changes via k8s CRs.

    Service Mesh - Implemented a service mesh across the whole kubernetes cluster through 𝗜𝘀𝘁𝗶𝗼 to control and monitor the traffic flow between components while securing the internal communication within the infrastrucutre

    Logging - Collected logs across the cluster using 𝗙𝗹𝘂𝗲𝗻𝘁𝗯𝗶𝘁 in JSON and pushed them to 𝗔𝗪𝗦 𝗖𝗹𝗼𝘂𝗱𝗪𝗮𝘁𝗰𝗵

    Monitoring - Used 𝗣𝗿𝗼𝗺𝗲𝘁𝗵𝗲𝘂𝘀 to collect metrics and 𝗚𝗿𝗮𝗳𝗮𝗻𝗮 to display them publicly through dashboards.

    Deployment - Used 𝗛𝗲𝗹𝗺 to manage application deployment on k8s and 𝗧𝗲𝗿𝗿𝗮𝗳𝗼𝗿𝗺 to manage infrastructure deployment through code.

    RAG - Implemented a RAG pipeline integrating 𝗛𝘂𝗴𝗴𝗶𝗻𝗴 𝗳𝗮𝗰𝗲 API and Pinecone Database, along with using a 𝗟𝗹𝗮𝗺𝗮𝟯:𝟴𝗕 𝗺𝗼𝗱𝗲𝗹 𝘀𝗲𝗹𝗳-𝗵𝗼𝘀𝘁𝗲𝗱 in 𝗞𝘂𝗯𝗲𝗿𝗻𝗲𝘁𝗲𝘀 to generate responses based on CVE data.

    Implemented data versioning, while storing the data in PostgreSQL to have different versions of CVE historical data and to track the periodic changes in the data over time. Implemented indexing to improve the querying speed.


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
