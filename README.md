<div align="center">
    <img src="./images/coderco.jpg" alt="CoderCo" width="300"/>
</div>

# URL Shortener - CoderCo ECS Project v2

A URL shortener on AWS. The app is provided. You build everything else.

```
POST /shorten  { "url": "https://example.com/long/path" }  →  { "short": "abc123ef" }
GET /abc123ef  →  302 redirect
GET /healthz   →  { "status": "ok" }
```

## Requirements

- **ECS Fargate** in private subnets, behind an **ALB** with **WAF**
- **VPC Endpoints** for AWS service access. No NAT gateways.
- **Database**: DynamoDB or RDS PostgreSQL - you choose, you justify
- **Blue/green deployments** via CodeDeploy with automatic rollback
- **GitHub Actions** CI/CD using **OIDC** (no long-lived AWS keys)
- **Terraform** with remote state, modular layout
- **Least-privilege IAM** throughout. No hardcoded credentials.

### App Config

- `TABLE_NAME` env var for DynamoDB
- `DATABASE_URL` env var for PostgreSQL
- Container port: **8080**

## The Deployment Question

You've deployed the service. Now a developer merges a PR and expects their change live within minutes - safely, with zero downtime.

**Design and document the full deployment workflow** in your README. Code merge to live traffic. Cover image builds, task definition updates, traffic shifting, rollback and observability.

## Deliverables

1. Working service (ALB DNS) with all endpoints functional
2. GitHub Actions workflows (CI + CD)
3. Terraform code for all infrastructure
4. Deployment workflow documentation
5. Evidence: OIDC trust policy, CodeDeploy blue/green, WAF on ALB, VPC Endpoints
6. README with decisions, trade-offs and database justification

## Acceptance Criteria

- No NAT gateways. Tasks pull images and write logs via VPC endpoints.
- Blue/green with auto-rollback on health check failure.
- Least-privilege IAM for your chosen database.
- OIDC auth in GitHub Actions.
- Remote Terraform state with locking.
- Deployment workflow section present and coherent.
- You can explain every resource. Copy-paste without understanding = resubmission.

## Cost Warning

Tear down when done. ALB + WAF cost money even idle.

```bash
cd infra/envs/dev && terraform destroy -auto-approve
```

[LocalStack](https://docs.localstack.cloud/getting-started/) works for local testing.

Everything else is on you. Commit small. Good luck.
