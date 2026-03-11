<div align="center">
    <img src="./images/coderco.jpg" alt="CoderCo" width="300"/>
</div>

# URL Shortener - CoderCo ECS Project v2

A URL shortener with click analytics on AWS. Three services, one cluster. The application code is provided. You build everything else.

## Services

| Service | Language | Port | Description |
|---------|----------|------|-------------|
| **api** | Python | 8080 | Shortens URLs, handles redirects, tracks clicks, publishes events to SQS |
| **worker** | Go | - | Polls SQS for click events, writes analytics to PostgreSQL |
| **dashboard** | Go | 8081 | Analytics API - top URLs, click stats, hourly breakdowns, recent events |

Read the code. Environment variables and endpoints are in the source files.

---

## Your Job

Write the Dockerfiles. Write the Terraform. Write the CI/CD pipeline. Deploy all three services to ECS Fargate on AWS.

### Requirements

- ECS Fargate - three separate services, one cluster
- Application Load Balancer with WAF routing to the correct service
- Database: DynamoDB or RDS PostgreSQL - you choose, you justify
- ElastiCache Redis (caching layer for the API)
- SQS queue (click events from API to worker)
- VPC with private subnets. No NAT gateways.
- GitHub Actions with OIDC. No long-lived AWS credentials.
- Zero-downtime deployments with rollback on failure
- Least-privilege IAM throughout
- Terraform with remote state
- Multi-stage Docker builds

### The Deployment Question

You've deployed the service. Now a developer merges a PR and expects their change live within minutes - safely, with zero downtime.

Design and document the full deployment workflow in your README. Code merge to live traffic.

### Deliverables

- [ ] Dockerfiles (one per service)
- [ ] Terraform for all infrastructure
- [ ] GitHub Actions CI/CD pipeline
- [ ] Deployment workflow documentation
- [ ] Working deployment - all services healthy, end-to-end flow functional
- [ ] README with your decisions, trade-offs and database justification

---

## Local Development

```bash
docker compose up --build
```

---

## Grading

- All three services running and healthy
- End-to-end flow works (shorten -> redirect -> analytics)
- Zero-downtime deployments with auto-rollback on health check failure
- No NAT gateways, no long-lived credentials, no hardcoded secrets
- Deployment workflow section present and coherent
- You can explain every resource you created

**Tear down when done.** ALB + WAF cost money even idle.

[LocalStack](https://docs.localstack.cloud/getting-started/) works for local testing of SQS.

Everything else is on you. Commit small. Good luck.
