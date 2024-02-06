# commands

minikube start

minikube status

docker ps -a

minikube addons list

minikube addons enable dashboard

minikube dashboard

minikube kubectl -- get deploy,rs,pods,svc -n kube-system

# defining alias
alias="minikube kubectl --"

# code and deploy
minikube kubectl -- apply -f infra/

minikube kubectl -- get deploy,rs,pods,svc,configmaps

minikube kubectl -- describe deploy post-service-deployment

minikube service post-app-service --url
> enter the service, interact with the application

minikube kubectl -- exec -it <pod-name-here> -- /bin/sh
> export
> ps ax

minikube kubectl -- delete pod <pod-name-here>

minikube kubectl -- get pods

minikube kubectl -- describe pod <new-pod-name-here>

## update

minikube kubectl -- apply -f infra/deployments.yaml

minikube kubectl -- describe deploy post-service-deployment


## breaking and rollback

minikube kubectl -- apply -f infra/deployments.yaml

minikube kubectl -- describe deploy post-service-deployment

minikube kubectl -- undo deploy post-service-deployment --to-revision=1

## finishing

minikube stop && minikube delete
