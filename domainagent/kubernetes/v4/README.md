https://claude.ai/chat/9abe7d81-24bc-43e0-84c9-331a18e6a5dd

kubernetes set up for workdomain.co.uk and felines.uk

mount felines-config-map.yaml as well.
kubectl run -i -t --rm debug --image=curlimages/curl --restart=Never -- curl http://feline-imagery-service

gcloud config set project aakubecontroller
Updated property [core/project].
(base) ant@aalenovo:~$ gcloud config get-value project
aakubecontroller

gcloud compute firewall-rules create allow-k8s-ingress --allow tcp:30080 --target-tags=aak8s-master 

