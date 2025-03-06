#! /bin/bash

# This script updates the HTML content in the Kubernetes ConfigMap (for the admin website workdomain.co.uk)

HTML_CONTENT=$(cat project-tracker.html | sed 's/"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')

# Create a temporary file for the updated ConfigMap
cat > configmap-temp.yaml << EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: project-tracker.html
data:
  index.html: |
$(sed 's/^/   /' project-tracker.html)
EOF

kubectl apply -f configmap-temp.yaml

rm configmap-temp.yaml

# restart the pods
kubectl rollout restart deployment project-tracker
