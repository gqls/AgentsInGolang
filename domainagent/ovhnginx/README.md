ssh -o IdentitiesOnly=yes -i ~/.ssh/ovh ubuntu@51.89.148.216
SrvOVH0114123!


server {
listen 80;
server_name felines.co.uk;

    location / {
        proxy_pass 35.214.74.66:30080;
	proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
listen 80;
server_name workdomain.co.uk;

    location / {
        proxy_pass 35.214.74.66:30080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	proxy_set_header X-Forwarded-Proto $scheme;
    }
}

/etc/nginx/sites-available
sudo ln ./workdomain.conf ../sites-enabled/
sudo chmod 777 workdomain.conf
sudo systemctl reload nginx
https://www.ovh.com/manager/#/dedicated/vps/vps-e68b309c.vps.ovh.net/dashboard
ssh ssh -o IdentitiesOnly=yes -i ~/.ssh/ovh ubuntu@51.89.148.216
https://jqls.atlassian.net/jira/software/projects/AG/boards/38?selectedIssue=AG-17

https://jqls.atlassian.net/jira/software/projects/AG/boards/38?selectedIssue=AG-19
google GCP
https://console.cloud.google.com/welcome?inv=1&invt=AbqtiA&project=aakubecontroller

dashboard:
https://console.cloud.google.com/home/dashboard?inv=1&invt=AbqtiA&project=aakubecontroller

compute:
https://console.cloud.google.com/compute/instances?inv=1&invt=AbqtiA&project=aakubecontroller

static ip address reserved in Google and pointed to master: 35.214.74.66
https://console.cloud.google.com/networking/addresses/list?inv=1&invt=AbqtiA&project=aakubecontroller

