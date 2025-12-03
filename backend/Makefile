NAMESPACE=social-network

.PHONY: build-base delete-volumes apply-namespace apply-pvc apply-db build-services deploy-users run-migrations logs-users logs-db all reset

# ==== Docker ====

build-base:
	docker build -t social-network/go-base -f docker/go/base.Dockerfile .

build-services:
	docker build -t social-network/users:dev -f services/users/Dockerfile .
# 	docker build -t social-network/api-gateway:dev -f gateway/Dockerfile .

docker-up:
	$(MAKE) build-base
	docker-compose up --build

delete-volumes:
	docker compose down
	docker volume rm backend_users-db-data backend_posts-db-data backend_chat-db-data backend_notifications-db-data

# ==== K8s ====

# 1.
apply-namespace:
	kubectl apply -f k8s/ --recursive --selector stage=namespace

#  2.
deploy-nginx:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm repo update
	helm upgrade --install nginx-ingress ingress-nginx/ingress-nginx \
		-n ingress-nginx --create-namespace

# 3.
apply-pvc:
	kubectl apply -f k8s/ --recursive --selector stage=pvc

# 4.
apply-db:
	kubectl apply -f k8s/ --recursive --selector stage=db

# 5.
run-migrations:
	kubectl apply -f k8s/ --recursive --selector stage=migration

# 6.
apply-configs:
	kubectl apply -R -f k8s --selector stage=config

# 7.
apply-apps:
	kubectl apply -f k8s/ --recursive --selector stage=app

# 8.
apply-ingress:
	kubectl apply -f k8s/nginx/api-gateway-ingress.yaml


build-proto:
	$(MAKE) -f shared/proto/protoMakefile generate

logs-users:
	kubectl logs -l app=users -n users -f

logs-db:
	kubectl logs -l app=users-db -n users -f

reset:
	kubectl delete namespace users --ignore-not-found=true
	kubectl create namespace users


all: build-base build-services apply-namespace deploy-nginx apply-pvc apply-db  deploy-users run-migrations apply-ingress

