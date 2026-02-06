# Social Network

## Overview

Social-sphere is an typical forum that has posts, comments, real time chat and notifications. It is powered by a cloud-native distributed platform built in Go and Next.js, designed with scalability, observability, availability and maintainability in mind. 

The platform was developed as part of a **Zone01 Athens** team project with a focus on robust service design and future scalability.

## Team
- Aleksis Gioldaseas
- Katerina Vamvasaki
- Ypatios Chaniotakos
- Vagelis Stefanopoulos
- Magnus Edvall

//**TODO link to hosted website (may be down)**

---

## Main Design Elements

* **Distributed architecture**
  * Core business services communicate via gRPC, NATS and Kafka
  * Optimization of communication using Redis caching
  * TLS termination on platforms surface
  * Fully scalable stateless services separated based on business domains
  * HTTP API + SSR for responsible User Experience
  * File Hosting with Minio, emulating CDN-like asset delivery

* **Scalable Databases**
  * PostgreSQL (Cloud Native PG operator) for core business CRUD, with replication (and sharding?)
  * Redis with Sentinel for High Availability caching/rate limiting

* **Messaging / Event Streaming**
  * NATS for lightweight pub/sub communication for live events
  * KRaft Kafka for High Availability and reliable inter-service communication of business critical events

* **Observability & Monitoring**
  * OpenTelemetry standard
  * Full VictoriaMetrics stack for metrics, traces and logs
  * Grafana Alloy for routing telemetry
  * Grafana dashboards for visualizations, monitoring and alerts

* **Security**
  * JWT-based authentication with expiration
  * Encrypted IDs exposed to front-end
  * Rate limiting and paginated requests
  * Authentication for access to databases

* **Cloud-Native Deployment**
  * Fully containerized using Docker
  * Deployed and orchestrated with Kubernetes

* **Go Services Code Architecture**
  * Shared business logic promoted to generic library-like packages
  * Wrappers around most 3rd party services, for imposing business rules through types and decoupling using interfaces
  * Composition based code organization with Dependency Injection
  * Clean and convinient telemetry functionality using Otel, with autoinstrumentaion of communication layers
  * Model based validation and encryption system
  * Advanced custom error handling system that also integrates gRpc and HTTP error codes
  * Custom built http middleware system

* **Developer Environments**
  * Run go and node.js services locally and 3rd party services with Docker, for the most lightweight DX
  * Run the full platform locally using K8s on Minikube
  * Run full platform on AKS (TODO)

* **CI/CD**
  //EXPLAIN HOW WE DEPLOY TO AKS AND UPDATE THE PLATFORM AFTER MAKING CHANGES

* **Testing**
  * End-to-end testing using containerized service that tests the api gateway, the services directly, and generates load for benchmarking.
  * Integration testing for complex libraries and inter-service features
  * Fuck unit tests!

---

## Architecture
![Page 1](./graph.png)

//**TODO indepth breakdown of the platform atlas, moved into another readme**
---

## Getting Started

### Prerequisites

* Docker & Docker Compose (or Colima)
* Kubernetes (Minikube / kind / cloud provider)


### Deployment Steps
//**TODO details like these should be moved and integrated into the makefile**
1. Install cloud native operator
2. Build docker base for golang and cloudnative image for Postgress
3. Build Docker images for services
4. Load images to K8s cluster (not needed on Colima)
5. Run deployment in this order using make commands.
	- make apply-kafka (need to first create name space if doing manual run. Name space creation is included in make commands)
	- make apply-namespace
	- make apply-pvc
	- make apply-configs
	- make apply-monitoring
	- make apply-db1  
    - make apply-db2
        - **Wait for db pods and all replicas to run**
    - make run-migrations 
	- make apply-apps   
        - **Wait for storage pod to run**
    - make apply-cors
    - make port-forward

    *Make commands run `kubectl apply -f` recursively on K8s dir for each stage selector*
    
    #### Ports exposed
    - Frontend -> localhost:3000 
    - Grafana -> localhost:3001
        - user: social
        - password: wearecool
    - victoria logs -> localhost:9428

---

## Observability

* **Metrics**: Latency, throughput, error rates (VictoriaMetrics)
* **Tracing**: gRPC call traces across services
* **Logging**: Structured logs captured per service

This allows **data-driven improvements** and performance optimization.

---

## Security Considerations

* JWT authentication with expiration
* Input validation for all endpoints
* Rate limiting for API requests
* Secrets are stored securely in Kubernetes secrets //TODO kubernetes secrets are not secure, aks vault is probably what we'll mention here
* Services insulated via gRPC endpoints

---

## Future Improvements

* CI/CD pipelines for automatic builds and deployments
* Enhanced JWT rotation and refresh tokens
* Multi-region deployments for high availability
* Load testing and performance tuning using pprof



---

## License

[MIT License](LICENSE)


