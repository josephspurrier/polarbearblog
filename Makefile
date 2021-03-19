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

# Load the environment variables.
include .env

.PHONY: default
default: gcp-push

################################################################################
# Deploy application
################################################################################

.PHONY: gcp-init
gcp-init:
	@echo Pushing to GCP storage initial files.
	gsutil mb -p $(GCP_PROJECT_ID) -l ${GCP_REGION} -c Standard gs://${GCP_BUCKET_NAME}
	gsutil versioning set on gs://${GCP_BUCKET_NAME}
	gsutil cp testdata/empty.json gs://${GCP_BUCKET_NAME}/storage/site.json
	gsutil cp testdata/empty.json gs://${GCP_BUCKET_NAME}/storage/session.json

.PHONY: gcp-push
gcp-push:
	@echo Pushing to GCP.
	gcloud builds submit --tag gcr.io/$(GCP_PROJECT_ID)/${GCP_IMAGE_NAME}
	gcloud run deploy --image gcr.io/$(GCP_PROJECT_ID)/${GCP_IMAGE_NAME} \
		--platform managed \
		--allow-unauthenticated \
		--region ${GCP_REGION} ${GCP_CLOUDRUN_NAME} \
		--update-env-vars SS_USERNAME=${SS_USERNAME} \
		--update-env-vars SS_SESSION_KEY=${SS_SESSION_KEY} \
		--update-env-vars SS_PASSWORD_HASH=${SS_PASSWORD_HASH} \
		--update-env-vars SS_MFA_KEY="${SS_MFA_KEY}" \
		--update-env-vars GCP_PROJECT_ID=${GCP_PROJECT_ID} \
		--update-env-vars GCP_BUCKET_NAME=${GCP_BUCKET_NAME} \
		--update-env-vars SS_ALLOW_HTML=${SS_ALLOW_HTML}

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