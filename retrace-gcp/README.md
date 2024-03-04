# retrace-gcp

This is a bit experimental because I'm not sure how I want to structure things to avoid dependency hell.

I think this should "just work" if you're on GCP with ambient credentials that have the `Cloud Trace Agent` role.

If you want to test this locally, though...

```console
PROJECT_ID=$(gcloud config get project)
USER_EMAIL=$(gcloud config get account)

# Create our service account.
gcloud iam service-accounts create retracer 

# Allow it to upload traces.
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:retracer@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/cloudtrace.agent"

# Give myself access to the retracer service account.
gcloud iam service-accounts add-iam-policy-binding \
    "retracer@${PROJECT_ID}.iam.gserviceaccount.com" \
    --member="user:${USER_EMAIL}" \
    --role="roles/iam.serviceAccountUser"

# Give myself the ability to impersonate service accounts.
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="user:${USER_EMAIL}" \
    --role="roles/iam.serviceAccountTokenCreator"

# Set retracer as our ADC.
gcloud auth application-default login \
    --project "${PROJECT_ID}" \
    --impersonate-service-account "retracer@${PROJECT_ID}.iam.gserviceaccount.com"

# This is janky and there is probably a better way to do this, but the exproter breaks if we don't.
# If you don't have `sponge` or `jq`, this is just adding setting "project_id": "$PROJECT_ID" as a top-level field of your application default credentials file.
ADC="$HOME/.config/gcloud/application_default_credentials.json"
jq --arg project "${PROJECT_ID}" '.project_id = $project' ${ADC} | sponge ${ADC}
```
