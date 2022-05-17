# Google Chat Webhook bot for Cloud Run

This is sample code to create chat bot to post something in your Google Chat room via a Webhook URL.

## Manual Deploy

### Create Webhook URL on the Google Chat room


### (Optional) Create Artifact Registry repo
```
export REPO_NAME=<Your Repo Name>

gcloud artifacts repositories create ${REPO_NAME} --repository-format=docker --location=asia-northeast1`
```

### Build & push the docker image to the repo
```
export PROJECT_ID=<Your Project ID>
export REPO_NAME=<Your Repo Name>
export APP_NAME=<Your App Name>

docker build -t asia-northeast1-docker.pkg.dev/${PROJECT_ID}/${REPO_NAME}/${APP_NAME} .
docker push asia-northeast1-docker.pkg.dev/${PROJECT_ID}/${REPO_NAME}/${APP_NAME}
```

### Create a Service Account to invoke the Cloud Run
```
gcloud iam service-accounts create ${APP_NAME}-invoker
```

### Deploy Cloud Run service (Copy the Service URL after the deployment)
```
export WEBHOOK_URL="<Webhook URL you created>"

gcloud run deploy ${APP_NAME} --image=asia-northeast1-docker.pkg.dev/${PROJECT_ID}/${REPO_NAME}/${APP_NAME} --no-allow-unauthenticated --set-env-vars=WEBHOOK_URL=${WEBHOOK_URL} --region=asia-northeast1

gcloud run services add-iam-policy-binding ${APP_NAME} --region='asia-northeast1' --member=serviceAccount:${APP_NAME}-invoker@${PROJECT_ID}.iam.gserviceaccount.com --role='roles/run.invoker'
```

https://cloud.google.com/sdk/gcloud/reference/run/deploy

### Create a Cloud Scheduler job to post the message automatically
```
export SERVICE_URL=<Service URL>
export SERVICE_ACCOUNT=<Service Account Email>

gcloud scheduler jobs create http chat-bot-every-wednesday-morning \
  --schedule "30 08 * * 3" --http-method=GET --uri=${SERVICE_URL} \
  --location=asia-northeast1 --time-zone=Asia/Tokyo \
  --oidc-service-account-email=${APP_NAME}-invoker@${PROJECT_ID}.iam.gserviceaccount.com --oidc-token-audience=${SERVICE_URL}
```

https://cloud.google.com/sdk/gcloud/reference/scheduler/jobs/create/http

## CI/CD with Cloud Build

TODO