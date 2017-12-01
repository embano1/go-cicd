# Minimalistic CI/CD with Go and Kubernetes

## Preparation

Clean Docker image cache

```
docker system prune -f
docker rmi -f `docker images|grep embano1/mynginx|awk '{print $3}'`
# Delete "mynginx" repo on https://hub.docker.com/
```

Build v1 image

```
cd examples/mynginx
./build_v1.sh
```

Deploy v1 to Kubernetes

```
kubectl run mynginx --image=embano1/mynginx:v1 --port 80 --image-pull-policy='Always'
kubectl expose deploy mynginx --type NodePort
kubectl scale deploy mynginx --replicas=10
minikube service mynginx
```

## Start CI/CD pipeline with "go-cicd"

```
cd examples/mynginx
go-cicd -e ./deploy_mynginx.sh -f deploy_mynginx.sh

   _____  ____     _____ _____       _____ _____
  / ____|/ __ \   / ____|_   _|     / ____|  __ \
 | |  __| |  | | | |      | |______| |    | |  | |
 | | |_ | |  | | | |      | |______| |    | |  | |
 | |__| | |__| | | |____ _| |_     | |____| |__| |
  \_____|\____/   \_____|_____|     \_____|_____/

        The World's most basic CI/CD Tool

[go-cicd] 23:05:43 Starting CI/CD process (executable: "./deploy_mynginx.sh", watching: "deploy_mynginx.sh" (isDir: false))
```

## Trigger Pipeline

Modify `${SITE}` variable in `deploy_mynginx.sh`.