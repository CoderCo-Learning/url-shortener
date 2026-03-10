<div align="center">
    <img src="./images/coderco.jpg" alt="CoderCo" width="300"/>
</div>

# URL Shortener - CoderCo ECS Project v2

The CoderCo ECS Project has levelled up.

For a long time, our community cut their teeth on the original ECS project - a solid rite of passage for containerisation, Terraform and AWS deployment. It's now getting a serious upgrade with production-ready features, better security and cost-efficient patterns you'll see in the real world.

You're not just deploying an app here. You're engineering an end-to-end workload the way a platform team would hand it off in production: private networking, zero-downtime deployments, no long-lived credentials and a CI/CD pipeline that any developer on the team could use to ship code safely.

## What You're Building

A URL shortener service on AWS. The app takes a long URL and returns a short code. Users hit the short link and get redirected. Simple app, serious infrastructure.

```
POST /shorten  { "url": "https://example.com/my/very/long/path" }
→ { "short": "abc123ef", "url": "https://example.com/my/very/long/path" }

GET /abc123ef
→ 302 redirect to https://example.com/my/very/long/path
```

The Python application is provided in `app/`. You build everything else.

## Requirements

### Infrastructure (Terraform)

- **ECS Fargate** service running in **private subnets** (no public IPs on tasks)
- **Application Load Balancer** in public subnets routing to the ECS service
- **VPC Endpoints** for all AWS service access (ECR, DynamoDB, CloudWatch Logs, S3). No NAT gateways.
- **DynamoDB** table with on-demand billing and point-in-time recovery enabled
- **AWS WAF** attached to the ALB with sensible rules
- **S3 + DynamoDB** for Terraform state backend (create this first in `infra/global/backend`)
- Modular Terraform layout split by environment (`infra/modules/`, `infra/envs/dev/`)

### Deployments

- **Blue/green deployments** via AWS CodeDeploy with automatic rollback on failed health checks
- Two target groups (blue + green), health check on `/healthz`
- App container port: **8080**

### CI/CD (GitHub Actions)

- **OIDC authentication** to AWS. No long-lived access keys in GitHub.
- **CI pipeline**: lint, unit tests, build image, push to ECR (on `main` or PR)
- **CD pipeline**: `terraform plan` on PR, `terraform apply` on merge to `main`, trigger CodeDeploy

### Security

- Task IAM role: only `dynamodb:GetItem` and `dynamodb:PutItem` on your table (least privilege)
- Execution role: only ECR pull + CloudWatch Logs write
- No secrets or credentials stored in code or GitHub settings

### App Configuration

The app needs one environment variable: `TABLE_NAME` (your DynamoDB table name).

## The Deployment Question

> You've built the infrastructure and deployed the service. Now a developer on your team wants to push a code change. They merge a PR to `main` and expect their change to be live within minutes - safely and with zero downtime.
>
> **Design and document the full deployment workflow.** From code merge to traffic serving the new version. Cover:
>
> - How the image gets built and tagged
> - How ECS picks up the new image (task definition revision, service update)
> - How CodeDeploy shifts traffic (blue/green strategy, health checks, rollback triggers)
> - How a developer would know their deployment succeeded or failed
> - What happens if the new version is broken
>
> Include this as a section in your README with a diagram or flow description. This is the kind of thing interviewers will ask you to whiteboard.

## Deliverables

1. **Working service** (ALB DNS) with all three endpoints functional
2. **GitHub Actions** workflows for CI and CD
3. **Terraform code** managing all infrastructure
4. **Deployment workflow documentation** (see above)
5. **Evidence** (screenshots or links):
   - OIDC role trust policy
   - CodeDeploy deployment showing blue/green + rollback test
   - WAF associated to ALB
   - VPC Endpoints in use
6. **README** with your decisions, trade-offs and deployment workflow

## Acceptance Criteria

These will be checked:

- No NAT gateways. Tasks still pull images from ECR and write logs.
- CodeDeploy blue/green shifts traffic and auto-rolls back on health check failure.
- Task role is scoped to your DynamoDB table only.
- GitHub workflow uses `id-token: write` and assumes an OIDC role.
- Terraform state is remote (S3 + state locking).
- Deployment workflow section is present and makes sense.
- You can explain every resource you created. Copy-paste without understanding = resubmission.

## Cost Warning

This deploys real AWS resources that cost money. Tear everything down when you're done:

```bash
cd infra/envs/dev && terraform destroy -auto-approve
```

ALB, WAF and DynamoDB storage will cost you even with zero traffic. Don't leave things running.

You can use [LocalStack](https://docs.localstack.cloud/getting-started/) to test locally before deploying to real AWS.

## Bonus (Optional)

- HTTPS via ACM + Route53
- Image scanning (Trivy, Snyk or similar) in CI
- CloudWatch dashboard (p50/p95 latency, 5xx rate, healthy host count)
- Infracost or tfsec in the PR workflow

Everything else is on you. Read the docs, iterate and commit small. Good luck.
