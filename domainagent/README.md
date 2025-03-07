https://claude.ai/chat/d20801f4-e057-4cc9-b762-09d63af611f1

domain project agents

debug thread: https://claude.ai/project/7d6447a5-0f0f-46c9-b452-3980e1be51eb

Jira
https://jqls.atlassian.net/jira/software/projects/AG/boards/38?selectedIssue=AG-21

kubectl get pods -n ingress-nginx
kubectl get svc -n ingress-nginx
sudo iptables -A INPUT -p tcp --dport 30080 -j ACCEPT
sudo iptables -A OUTPUT -p tcp --sport 30080 -j ACCEPT
gcloud compute firewall-rules create allow-k8s-ingress --allow tcp:30080 --target-tags=aakubecontroller
gcloud projects list
curl localhost:30080 -H "Host: workdomain.co.uk"
curl 35.214.74.66:30080 -H "Host: workdomain.co.uk" -v
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller
kubectl get ingress
kubectl describe ingress workdomain-ingress
kubectl describe svc ingress-nginx-controller -n ingress-nginx

kubectl -n monitoring rollout restart deployment prometheus
kubectl -n monitoring rollout restart deployment grafana

OVH
nginx
proxy_pass http://35.214.74.66:30080;
sudo nano /etc/nginx/sites-enabled/felines.conf
sudo nginx -t
sudo systemctl reload nginx
