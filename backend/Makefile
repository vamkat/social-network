NAMESPACE=social-network

.PHONY: build-base apply-namespace apply-db build-services deploy-users run-migrations logs-users logs-db all reset

build-base:
	docker build -t social-network/go-base -f docker/go/base.Dockerfile .

apply-namespace:
	kubectl apply -f k8s/users/namespace-users.yaml
	kubectl apply -f k8s/nginx/namespace-nginx.yaml

deploy-nginx:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm repo update
	helm upgrade --install nginx-ingress ingress-nginx/ingress-nginx \
		-n ingress-nginx --create-namespace

apply-ingress:
	kubectl apply -f k8s/nginx/api-gateway-ingress.yaml

apply-db:
	kubectl apply -f k8s/users-db/pvc.yaml
	kubectl apply -f k8s/users-db/deployment.yaml
	kubectl apply -f k8s/users-db/service.yaml

build-services:
	docker build -t social-network/users:dev -f services/users/Dockerfile .

deploy-users:
	kubectl apply -f k8s/users/deployment.yaml
	kubectl apply -f k8s/users/service.yaml
	kubectl apply -f k8s/users/secret-users.yaml

run-migrations:
	kubectl delete job users-migrate -n users --ignore-not-found=true
	kubectl apply -f k8s/users-db/migration-job.yaml

logs-users:
	kubectl logs -l app=users -n users -f

logs-db:
	kubectl logs -l app=users-db -n users -f

reset:
	kubectl delete namespace users --ignore-not-found=true
	kubectl create namespace users

all: apply-namespace deploy-nginx apply-db build-services deploy-users run-migrations apply-ingress

