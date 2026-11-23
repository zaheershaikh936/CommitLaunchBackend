# Static Site Deployment Service — Phase 1 Requirements

## 1. Goal

Build a backend service that takes a public GitHub repository containing a ready-to-serve static site, publishes it to Amazon S3, and returns a browser URL. This phase covers HTML, CSS, browser JavaScript, and public assets such as images and fonts.

## 2. Scope

- One trusted operator deploys repositories they control.
- Deployments are triggered manually through the API.
- All deployments use one preconfigured S3 bucket with static website hosting enabled and `index.html` configured as the index document.
- Each deployment gets its own URL. The S3 website endpoint uses HTTP, and published files are publicly readable. See [AWS S3 website endpoint documentation](https://docs.aws.amazon.com/AmazonS3/latest/userguide/WebsiteEndpoints.html).
- The backend manages deployment jobs and status. S3 serves the site, and the browser executes its JavaScript.

## 3. Functional requirements

| ID | Requirement | Expected behavior |
|---|---|---|
| FR-01 | Submit a deployment | Accept a public GitHub repository URL and a branch. Validate the input, create a deployment record, and return its ID with status `queued`. |
| FR-02 | Fetch the source | A background worker resolves the branch to a commit SHA, records it, and fetches that exact revision. |
| FR-03 | Validate the site | Require `index.html` at the repository root. Publish only files intended for public serving. Exclude Git metadata, environment files, and other non-site files. Preserve asset paths. Site authors should use relative asset URLs so resources resolve under the deployment prefix. |
| FR-04 | Upload to S3 | Upload the validated files under `sites/<project-slug>/<deployment-id>/`. Set each object's content type. Upload assets before the root `index.html`. |
| FR-05 | Verify publication | After all uploads succeed, fetch the deployment URL through the S3 website endpoint. Mark the deployment `ready` and return its URL after a successful response. |
| FR-06 | Report progress | A deployment lookup returns the ID, repository, branch, commit SHA once resolved, status, timestamps, URL when ready, and a clear failure reason when applicable. |
| FR-07 | Handle failures | If source retrieval, validation, upload, or verification fails, record the failed stage and mark the deployment `failed`. Keep the URL out of the API response for failed deployments. Previous deployments retain their own files and URLs. |
| FR-08 | Retry safely | If the worker retries a job after a transient error, reuse its recorded commit, deployment ID, and S3 prefix. Prevent duplicate processing from changing a completed deployment. |

## 4. Deployment flow

1. The operator sends a repository URL and branch to the API.
2. The API records the request and returns a deployment ID.
3. A worker fetches the selected revision and validates the site.
4. The worker uploads the assets and index document to the deployment's prefix.
5. The worker verifies the public URL and updates the deployment record.
6. The operator checks that deployment ID to retrieve the URL or failure reason.

Suggested status progression is `queued → fetching → validating → uploading → verifying → ready`. Any failed processing stage can end in `failed`.

## 5. API and deployment record

| Operation | Purpose |
|---|---|
| `POST /deployments` | Accept a repository URL and branch. Return HTTP 202 with the deployment ID and initial status. |
| `GET /deployments/{id}` | Return the deployment status and metadata. Return HTTP 404 for an unknown ID. |

Persist the deployment ID, project slug, repository URL, branch, resolved commit SHA, status, S3 prefix, public URL, timestamps, and failure details. The public URL has the form `http://<bucket-website-endpoint>/sites/<project-slug>/<deployment-id>/`.

## 6. Acceptance criteria

- Deploy a sample portfolio and open its returned URL in a browser. The index page displays, the stylesheet loads, and a JavaScript interaction works.
- Confirm the deployment record shows the exact commit SHA and reaches `ready` after verification.
- Submit a commit missing `index.html`. Confirm the deployment reaches `failed` with a clear reason and the previous URL still works.
- Simulate an S3 upload failure. Confirm the deployment does not return a ready URL and a previous successful deployment remains accessible.
- Redeploy a valid commit as a new deployment. Confirm it receives a separate URL and both successful URLs work independently.

## 7. Later phases

Build commands and framework builds, Docker and server-side runtimes, private repositories, push-triggered deployments, stable project URLs, custom domains, HTTPS, CloudFront, and rollback controls.
