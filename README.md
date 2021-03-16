# Polar Bear Blog 🐻‍❄️

Lightweight blogging system for a single author. Written in Go and deploys to your own GCP project with a few commands. It's a derivative of the beautifully simple [Bear Blog 🐻](https://github.com/HermanMartinus/bearblog/). The data storage and session storage are stored in Google Cloud Storage as objects. Depending on the traffic and blog size, it should (not guaranteed) cost less than $1 USD per month (compute and storage) to host this blog because it will be deployed to [Google Cloud Run](https://cloud.google.com/run/pricing) which bills to the nearest 100 millisecond. You can also [map your own domain name](https://cloud.google.com/run/docs/mapping-custom-domains) and Google will provide a free SSL certificate. This project uses `make` to simplify the deployment process.

Out of all the ways to get a Go web application online, this is one of the quickest, cheapest, and simplest I found.

You can see a real website using this stack [here](https://www.josephspurrier.com/).

## Quickstart

By following these instructions, you can get a blog public pretty easily:

- Create a Google GCP project
- Generate a .env file with this content and then fill each variable in:

```bash
# App Configuration
## Email to use to login to the platform at: https://example.run.app/login/admin
SS_USERNAME=
## GCP project ID.
GCP_PROJECT_ID=
## GCP bucket name (this can be one that doesn't exist yet).
GCP_BUCKET_NAME=
## Session key to encrypt the cookie store. Generate with: make privatekey
SS_SESSION_KEY=
## Password hash that is base64 encoded. Generate with: make passhash passwordhere
SS_PASSWORD_HASH=
## MFA (TOTP) that works with apps like Google Authenticator. Generate with: make mfa
SS_MFA_KEY=
## Optional: set the timezone from here:
## https://golang.org/src/time/zoneinfo_abbrs_windows.go
# SS_TIMEZONE=America/New_York

# GCP Deployment
## Name of the docker image that will be created and stored in GCP Repository.
GCP_IMAGE_NAME=
## Name of the Cloud Run service to create.
GCP_CLOUDRUN_NAME=
## Region (not zone) where the Cloud Run service will be created:
## https://cloud.google.com/compute/docs/regions-zones#available
GCP_REGION=us-central1

# MFA Configuration
## Friendly identifier when you generate the MFA string.
SS_ISSUER=www.example.com

# Local Development
## Set this to any value to allow you to do testing locally without GCP access.
## See 'Local Development Flag' section below for more information.
SS_LOCAL=true
```
- Run this command to initialize the store by creating the GCP bucket, enabling versioning, and then copying 2 blank files to the bucket: `make gcpstore`. You will need to have the [Google Cloud SDK installed](https://cloud.google.com/sdk/docs/install). You will also need a [service account key](https://console.cloud.google.com/apis/credentials/serviceaccountkey) downloaded on your system with an environment variables set to the JSON file like this: `GOOGLE_APPLICATION_CREDENTIALS=~/gcp-cloud-key.json`.
- Run this command to build the docker image, push to the Google repository, and then create a Cloud Run job: `make`

Once the process completes in a few minutes, you should get a URL to access the website. The login page is located at (replace with your real URL): https://example.run.app/login/admin.

To login, you'll need:

- the username from `SS_USERNAME`
- the password from which the `SS_PASSWORD_HASH` was derived
- the 6 digit MFA code generated from an app like Google Authenticator from the `SS_MFA_KEY`

Once you are logged in, you should see a new menu option call `Dashboard`. From this screen, you'll be able to make changes to the site as we as the home page. To add new posts, click on `Posts` and add the posts or pages from there.

## Development

If you would like to make changes to the code, I recommend these tools to help streamline your workflow.

```bash
# Install air to allow hot reloading so you can make changes quickly.
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

# Install direnv and hook into your shell. This allows you to manage 
# https://direnv.net/docs/installation.html
```

Once you have `direnv` installed, create .envrc file. Update the `GOOGLE_APPLICATION_CREDENTIALS` variable to the correct location on your hard drive of the app credentials. You can generate and download a service account key from: https://console.cloud.google.com/apis/credentials/serviceaccountkey.

```bash
# Load the shared environment variables (shared with Makefile).
# Export the vars in .env into the shell.
export $(egrep -v '^#' .env | xargs)

export PATH=$PATH:$(pwd)/bin
export GOOGLE_APPLICATION_CREDENTIALS=~/gcp-cloud-key.json
```

You can then use this commands to test and then to deploy.

```bash
# Start hot reload. The web application should be available at: http://localhost:8080
air

# Upload new version of the application to Google Cloud Run.
make
```

### Local Development Flag

When `SS_LOCAL` is set, the following things will happen:

- data storage will be local instead of in GCP storage
- redirects will not happen if not accessing through the primary URL
- MFA will accept any number instead of validating
- Google Analytics will be disabled
- Disqus will be disabled