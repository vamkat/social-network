NAMESPACE=social-network

.PHONY: build-base delete-volumes build-cnpg apply-namespace apply-pvc apply-db build-services deploy-users run-migrations logs-users logs-db all reset

# ==== Docker ====

build-base:
	docker build -t social-network/go-base -f backend/docker/go/base2.Dockerfile .

build-services:
	docker build -f backend/docker/services/api-gateway.Dockerfile -t social-network/api-gateway:dev .
	docker build -f backend/docker/services/chat.Dockerfile -t social-network/chat:dev .
	docker build -f backend/docker/services/live.Dockerfile -t social-network/live:dev .
	docker build -f backend/docker/services/media.Dockerfile -t social-network/media:dev .
	docker build -f backend/docker/services/notifications.Dockerfile -t social-network/notifications:dev .
	docker build -f backend/docker/services/posts.Dockerfile -t social-network/posts:dev .
	docker build -t social-network/users:dev -f backend/docker/services/users.Dockerfile .

docker-up:
	$(MAKE) create-network
	$(MAKE) build-base
	docker-compose up --build

delete-volumes:
	docker compose down
	docker volume rm backend_users-db-data backend_posts-db-data backend_chat-db-data backend_notifications-db-data backend_media-db-data

docker-test:
	$(MAKE) create-network
	docker compose -f docker-test.yml up --build

docker-up-test:
	$(MAKE) docker-up
	$(MAKE) docker-test

api:
	docker-compose up api-gateway --build

# this network is used to let the tester connect to the rest of the containers
# since the tester won't be part of the normal docker-compose, it will have it's own so that someone gotta go out of their way to test
create-network:
	@docker network inspect social-network >nul 2>&1 || docker network create social-network

# ==== K8s ====

# Preliminary
build-cnpg:
	docker buildx bake -f backend/docker/cnpg/bake.hcl postgres16-cloud-native

# 1.
apply-namespace:
	kubectl apply -f backend/k8s/ --recursive --selector stage=namespace

#  2.
apply-configs:
	kubectl apply -R -f backend/k8s/ --recursive --selector stage=config

# 3.
deploy-nginx:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm repo update
	helm upgrade --install nginx-ingress ingress-nginx/ingress-nginx \
		-n ingress-nginx --create-namespace

# 4.
apply-db:
	kubectl apply -f backend/k8s/ --recursive --selector stage=db

# 5.
run-migrations:
	kubectl apply -f backend/k8s/ --recursive --selector stage=migration


# 6.
apply-pvc:
	kubectl apply -f k8s/ --recursive --selector stage=pvc

# 7.
apply-apps:
	kubectl apply -f backend/k8s/ --recursive --selector stage=app

# 8.
apply-ingress:
	kubectl apply -f backend/k8s/nginx/api-gateway-ingress.yaml



build-proto:
	$(MAKE) -f backend/shared/proto/protoMakefile generate

logs-users:
	kubectl logs -l app=users -n users -f

logs-db:
	kubectl logs -l app=users-db -n users -f

reset:
	kubectl delete namespace users --ignore-not-found=true
	kubectl create namespace users


all: 
	build-base 
	build-cnpg 
	build-services 
	apply-namespace
	apply-configs
	deploy-nginx 
	apply-db  
	run-migrations 
	apply-pvc
	apply-apps
	apply-ingress

