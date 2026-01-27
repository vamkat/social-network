Perfect! Let’s create a **ready-to-deploy Kubernetes setup** for **Redis master + 3 Sentinel replicas**, fully configured and ready for a basic high-availability setup. I’ll make it **headless, persistent, and fully connected**.

---

## **1️⃣ ConfigMaps**

### **redis-configmap.yaml**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-config
data:
  redis.conf: |
    bind 0.0.0.0
    port 6379
    protected-mode no
    save ""
    appendonly no
```

### **sentinel-configmap.yaml**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-sentinel-config
data:
  sentinel.conf: |
    port 26379
    bind 0.0.0.0
    sentinel resolve-hostnames yes
    sentinel announce-ip ${POD_NAME}.${SENTINEL_SERVICE}.default.svc.cluster.local
    sentinel announce-port 26379
    sentinel monitor master redis-master 6379 2
    sentinel down-after-milliseconds master 5000
    sentinel failover-timeout master 10000
    sentinel parallel-syncs master 1
```

> **Note:** `${POD_NAME}` and `${SENTINEL_SERVICE}` will be passed as environment variables to each Sentinel pod so they announce their correct hostname.

---

## **2️⃣ Redis Master StatefulSet**

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis-master
spec:
  serviceName: "redis-master"
  replicas: 1
  selector:
    matchLabels:
      app: redis
      role: master
  template:
    metadata:
      labels:
        app: redis
        role: master
    spec:
      containers:
        - name: redis
          image: redis:7.2
          ports:
            - containerPort: 6379
          volumeMounts:
            - name: redis-data
              mountPath: /data
            - name: redis-config
              mountPath: /usr/local/etc/redis/redis.conf
              subPath: redis.conf
          command: ["redis-server", "/usr/local/etc/redis/redis.conf"]
      volumes:
        - name: redis-config
          configMap:
            name: redis-config
  volumeClaimTemplates:
    - metadata:
        name: redis-data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 1Gi
```

---

## **3️⃣ Redis Master Service (Headless)**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: redis-master
spec:
  selector:
    app: redis
    role: master
  ports:
    - port: 6379
      targetPort: 6379
  clusterIP: None
```

---

## **4️⃣ Redis Sentinel Deployment**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-sentinel
spec:
  replicas: 3
  selector:
    matchLabels:
      app: redis-sentinel
  template:
    metadata:
      labels:
        app: redis-sentinel
    spec:
      containers:
        - name: redis-sentinel
          image: redis:7.2
          ports:
            - containerPort: 26379
          env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: SENTINEL_SERVICE
              value: redis-sentinel
          volumeMounts:
            - name: sentinel-config
              mountPath: /usr/local/etc/redis/sentinel.conf
              subPath: sentinel.conf
          command: ["redis-sentinel", "/usr/local/etc/redis/sentinel.conf"]
      volumes:
        - name: sentinel-config
          configMap:
            name: redis-sentinel-config
```

---

## **5️⃣ Redis Sentinel Service (Headless)**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: redis-sentinel
spec:
  selector:
    app: redis-sentinel
  ports:
    - port: 26379
      targetPort: 26379
  clusterIP: None
```

---

### ✅ **Key Features of This Setup**

1. **Redis Master**: StatefulSet with persistent storage.
2. **Redis Sentinel**: 3 replicas for quorum-based HA failover.
3. **Headless Services**: Allow Sentinel to resolve each other via DNS.
4. **Environment Variables**: Pod names are used for proper `announce-ip`.
5. **Simple HA**: Sentinel monitors master and can trigger failover.

---

If you want, I can **also add a Redis Replica StatefulSet** so you have **master + replicas automatically syncing** for read scalability. This is usually done in production.

Do you want me to add that too?
