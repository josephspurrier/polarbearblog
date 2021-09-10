# This Makefile is an easy way to run common operations.
# Execute commands like this:
# * make
# * make gcp-init
# * make gcp-push
# * make privatekey
# * make mfa
# * make passhash passwordhere
# * make local-init
# * make local-run
# * make kube-init
# * make kube-push
# * make s3-init

# Load the environment variables.
include .env

.PHONY: default
default: gcp-push

################################################################################
# Deploy application
################################################################################

.PHONY: s3-init
aws-init:
	@echo Pushing the initial files to AWS S3
	aws s3 mb s3://${PBB_AWS_BUCKET_NAME} --region ${PBB_AWS_REGION}
	aws s3api put-bucket-versioning --bucket ${PBB_AWS_BUCKET_NAME} --versioning-configuration Status=Enabled
	aws s3 cp storage/site.json s3://${PBB_AWS_BUCKET_NAME}/storage/site.json
	aws s3 cp storage/session.bin s3://${PBB_AWS_BUCKET_NAME}/storage/session.bin

.PHONY: gcp-init
gcp-init:
	@echo Pushing the initial files to Google Cloud Storage.
	gsutil mb -p $(PBB_GCP_PROJECT_ID) -l ${PBB_GCP_REGION} -c Standard gs://${PBB_GCP_BUCKET_NAME}
	gsutil versioning set on gs://${PBB_GCP_BUCKET_NAME}
	gsutil cp testdata/empty.json gs://${PBB_GCP_BUCKET_NAME}/storage/site.json
	gsutil cp testdata/empty.json gs://${PBB_GCP_BUCKET_NAME}/storage/session.json

.PHONY: gcp-push
gcp-push:
	@echo Pushing to Google Cloud Run.
	gcloud builds submit --tag gcr.io/$(PBB_GCP_PROJECT_ID)/${PBB_GCP_IMAGE_NAME}
	gcloud run deploy --image gcr.io/$(PBB_GCP_PROJECT_ID)/${PBB_GCP_IMAGE_NAME} \
		--platform managed \
		--allow-unauthenticated \
		--region ${PBB_GCP_REGION} ${PBB_GCP_CLOUDRUN_NAME} \
		--update-env-vars PBB_USERNAME=${PBB_USERNAME} \
		--update-env-vars PBB_SESSION_KEY=${PBB_SESSION_KEY} \
		--update-env-vars PBB_PASSWORD_HASH=${PBB_PASSWORD_HASH} \
		--update-env-vars PBB_MFA_KEY="${PBB_MFA_KEY}" \
		--update-env-vars PBB_GCP_PROJECT_ID=${PBB_GCP_PROJECT_ID} \
		--update-env-vars PBB_GCP_BUCKET_NAME=${PBB_GCP_BUCKET_NAME} \
		--update-env-vars PBB_ALLOW_HTML=${PBB_ALLOW_HTML}

.PHONY: kube-init
kube-init:
	@echo Installing to Kubernetes.
	helm install polarbearblog ./deployment/helm/polarbearblog \
	    -f ./deployment/helm/polarbearblog/values.yaml \
		--set env.PBB_AWS_REGION=${PBB_AWS_REGION} \
		--set env.PBB_USERNAME=${PBB_USERNAME} \
		--set env.PBB_AWS_BUCKET_NAME=${PBB_AWS_BUCKET_NAME} \
		--set env.PBB_ALLOW_HTML=${PBB_ALLOW_HTML} \
		--set env.PBB_CLOUD_PROVIDER=${PBB_CLOUD_PROVIDER} \
		--set secrets.PBB_SESSION_KEY=${PBB_SESSION_KEY} \
		--set secrets.PBB_PASSWORD_HASH=${PBB_PASSWORD_HASH} \
		--set secrets.AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
		--set secrets.AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}

.PHONY: kube-push
kube-push:
	@echo Pushing to Kubernetes.
	helm upgrade polarbearblog ./deployment/helm/polarbearblog \
	    -f ./deployment/helm/polarbearblog/values.yaml \
		--set env.PBB_AWS_REGION=${PBB_AWS_REGION} \
		--set env.PBB_USERNAME=${PBB_USERNAME} \
		--set env.PBB_AWS_BUCKET_NAME=${PBB_AWS_BUCKET_NAME} \
		--set env.PBB_ALLOW_HTML=${PBB_ALLOW_HTML} \
		--set env.PBB_CLOUD_PROVIDER=${PBB_CLOUD_PROVIDER} \
		--set env.PBB_SESSION_KEY=${PBB_SESSION_KEY} \
		--set env.PBB_PASSWORD_HASH=${PBB_PASSWORD_HASH} \
		--set env.AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
		--set env.AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}

.PHONY: lambda-init
lambda-init:
	@echo Pushing the initial files to AWS Lambda.
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build main.go
	zip function.zip main
	aws lambda create-function --function-name polarbearblog \
		--runtime go1.x \
		--role arn:aws:iam::${PBB_AWS_ACCOUNT_ID}:role/lambda_basic_execution \
		--handler main \
		--zip-file fileb://function.zip \
		--timeout 300 \
		--memory-size 128
			
.PHONY: lambda-push
lambda-push:
	@echo Pushing to AWS Lambda.
	aws lambda update-function-code --function-name polarbearblog --zip-file fileb://function.zip

.PHONY: privatekey
privatekey:
	@echo Generating private key for encrypting sessions.
	@echo You can paste private key this into your .env file:
	@go run cmd/privatekey/main.go

.PHONY: mfa
mfa:
	@echo Generating MFA for user.
	@echo You can paste private key this into your .env file:
	@go run cmd/mfa/main.go

# Save the ARGS.
# https://stackoverflow.com/a/14061796
ifeq (passhash,$(firstword $(MAKECMDGOALS)))
  ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(ARGS):;@:)
endif

.PHONY: passhash
passhash:
	@echo Generating password hash.
	@echo You can paste private key this into your .env file:
	@go run cmd/passhash/main.go ${ARGS}

.PHONY: local-init
local-init:
	@echo Creating session and site storage files locally.
	cp storage/initial/session.bin storage/session.bin
	cp storage/initial/site.json storage/site.json

.PHONY: local-run
local-run:
	@echo Starting local server.
	LOCALDEV=true go run main.go
