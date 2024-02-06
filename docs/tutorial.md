# Deploying a simple Kubernetes cluster locally
This guide was mostly based on the course [Linux Foundation - Introduction to Kubernetes](https://www.edx.org/learn/kubernetes/the-linux-foundation-introduction-to-kubernetes) and [this excellent blog of on the subject](https://blog.andreaskrahl.de/tag/introduction-to-kubernetes/). The objective is to deploy a simple application cluster using `minikube`, exploring the concepts of Pods, ReplicaSets, Deployments, Services and ConfigMaps.

Also, we will explore the `kubectl` CLI to monitor and manage our simple cluster.

The following sections assume you already have a working `docker` and `minikube` installation. They also assume basic knowledge about the inner workings of a Kubernetes cluster. If you want to learn about this, I strongly recommend the aforementioned links.

## What we will cover?

- Cluster bootstrap with `minikube`
- Deploy the application: frontent, backend and database
- Deploy a ConfigMap
- Exposing Pods with Services
- Remove Pods and check ReplicaSet behavior
- Rollout changes to the cluster
- Rollback changes to the cluster
- Destroy everything

### Cluster bootstrap with `minikube`
`minikube start` will set-up everything we need to get a working Kubernetes cluster. By default, the `docker` isolation driver will be used to isolate the `minikube` environment in your machine. However, you can select which driver to use with the argument `--driver`.
Regardless of your choice, `minikube` will download the images required to boot the cluster. Then, it will download a version of `kubectl` for cluster management. You use a standalone CLI version to manage the cluster, if you prefer.

`minikube start` outputs something along those lines:
```sh
$ minikube start
😄  minikube v1.32.0 on Ubuntu 20.04
✨  Using the docker driver based on existing profile
💨  For improved Docker performance, enable the overlay Linux kernel module using 'modprobe overlay'
👍  Starting control plane node minikube in cluster minikube
🚜  Pulling base image ...
🤷  docker "minikube" container is missing, will recreate.
🔥  Creating docker container (CPUs=2, Memory=3736MB) ...
🐳  Preparing Kubernetes v1.28.3 on Docker 24.0.7 ...
    ▪ Generating certificates and keys ...
    ▪ Booting up control plane ...
    ▪ Configuring RBAC rules ...
🔗  Configuring bridge CNI (Container Networking Interface) ...
🔎  Verifying Kubernetes components...
    ▪ Using image gcr.io/k8s-minikube/storage-provisioner:v5
🌟  Enabled addons: storage-provisioner, default-storageclass
💡  kubectl not found. If you need it, try: 'minikube kubectl -- get pods -A'
🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default

```

First of all, let's define an alias to `kubectl`. This will make our commands much more concise:
```sh
alias kubectl="minikube kubectl --" # on .bashrc or similar
```

After sourcing your `.bashrc` file (`source ~/.bashrc`), the alias is available. Let's check which version of `kubectl` we are using:
```sh
$ kubectl version
Client Version: v1.28.3
Kustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3
Server Version: v1.28.3
```

We have a working environment now. First, let's take a look at the `~/kube/conf` file. This file configures the access to the Kubernetes cluster, and is automatically setup by `minikube`:
```sh
$ cat ~/.kube/config
apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/luiz.melo/.minikube/ca.crt
    extensions:
    - extension:
        last-update: Tue, 06 Feb 2024 16:08:20 -03
        provider: minikube.sigs.k8s.io
        version: v1.32.0
      name: cluster_info
    server: https://127.0.0.1:44521
  name: minikube
contexts:
- context:
    cluster: minikube
    extensions:
    - extension:
        last-update: Tue, 06 Feb 2024 16:08:20 -03
        provider: minikube.sigs.k8s.io
        version: v1.32.0
      name: context_info
    namespace: default
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/luiz.melo/.minikube/profiles/minikube/client.crt
    client-key: /home/luiz.melo/.minikube/profiles/minikube/client.key
```

We can see that `minikube` configured certificates for both the cluster and default user, and see the address of the API Server. We can see this info with `kubectl` too: `kubectl config view`.

We can check that the control plane node is running and the CoreDNS addon is enabled with:
```sh
$ kubectl cluster-info
Kubernetes control plane is running at https://127.0.0.1:44521
CoreDNS is running at https://127.0.0.1:44521/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
```

minikube also reports that everything is fine:
```sh
$ minikube status
minikube
type: Control Plane
host: Running
kubelet: Running
apiserver: Running
kubeconfig: Configured
```

### Deploy the application: frontent, backend and database
Application containers are deployed on Pods by Kubernetes. However, working directly with Pods can be cumbersome, especially if we need replication or scaling to happen smoothly.

As such, to deploy our containers, we use a workload object called Deployment. The deployment will be responsible to manage a ReplicaSet object. A ReplicaSet guarantees that the desired number of Pod replicas is available on our environment at all times. This ensures high availability for our apps. ReplicaSet will manage the number of pods, so if there are less Pods than desired, ReplicaSet will spin more of it.

Workload objects are declared with manifest files. As such, let's declare a manifest creating a mongodb deployment named `deployments.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mongodb-deployment
  labels:
    app: mongodb
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mongodb
  template:
    metadata:
      labels:
        app: mongodb
    spec:
      containers:
      - name: mongodb
        image: mongo:6.0.13-jammy
        ports:
        - containerPort: 27017
```

There's some things to unpack here. First of all, manifest files are mostly composed by four sections: apiVersion, kind, metadata and spec.

apiVersion defines the API version for the managed object. Different objects use different API versions.

kind defines the type of object being managed. In this case, the type is Deployment.

metadata defines identifiable data used to recognize objects as part of some application. This includes, name, namespace (in this case, the default namespace is being used) and labels. Labels are used to select and group Kubernetes objects.

spec defines the desired configuration for the given object. In this example, we specify that we want a single replica (single Pod) for the Deployment. Inside the spec, we define the selector which will be used to identify which Pods are to be managed by this ReplicaSet. The template field defines metadata and the spec for the Pods. Note that the labels defined internally for the Pods match the labels used as selector on the ReplicaSet.

To deploy the manifest, execute:
```sh
$ kubectl apply -f deployment.yaml
deployment.apps/mongodb-deployment created
```

We can check the status of the Pods, Deployments and ReplicaSets with:
```sh
$ kubectl get pods,deploy,rs
NAME                                      READY   STATUS    RESTARTS   AGE
pod/mongodb-deployment-66c66c656b-vmk8h   1/1     Running   0          97s

NAME                                 READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/mongodb-deployment   1/1     1            1           97s

NAME                                            DESIRED   CURRENT   READY   AGE
replicaset.apps/mongodb-deployment-66c66c656b   1         1         1       97s
```

We can also describe a Deployment to view more details:
```sh
$ kubectl describe deploy mongodb-deployment
Name:                   mongodb-deployment
Namespace:              default
CreationTimestamp:      Tue, 06 Feb 2024 16:42:46 -0300
Labels:                 app=mongodb
Annotations:            deployment.kubernetes.io/revision: 1
Selector:               app=mongodb
Replicas:               1 desired | 1 updated | 1 total | 1 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25% max unavailable, 25% max surge
Pod Template:
  Labels:  app=mongodb
  Containers:
   mongodb:
    Image:        mongo:6.0.13-jammy
    Port:         27017/TCP
    Host Port:    0/TCP
    Environment:  <none>
    Mounts:       <none>
  Volumes:        <none>
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      True    MinimumReplicasAvailable
  Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  <none>
NewReplicaSet:   mongodb-deployment-66c66c656b (1/1 replicas created)
Events:
  Type    Reason             Age    From                   Message
  ----    ------             ----   ----                   -------
  Normal  ScalingReplicaSet  3m20s  deployment-controller  Scaled up replica set mongodb-deployment-66c66c656b to 1
```

The Pod is available, and the ReplicaSet reports success. Note the Labels section inside the Pod template, reflecting our user-defined labels. We can also check the revision number and the events section. These will be useful for rollback purposes.

All is good, but the application needs frontend and backend. First, let's add the backend manifest. We will reference a prebuilt container image.
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: post-api-deployment
  labels:
    app: post-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: post-api
  template:
    metadata:
      labels:
        app: post-api
    spec:
      containers:
      - name: post-api
        image: lfmtsml/post-service:latest
        ports:
        - containerPort: 8090
        env:
          - name: MONGODB_URI
            value: "mongodb://mongodb-service.default.svc.cluster.local:27017"
          - name: MONGODB_DATABASE
            value: "post-database"
```

You can add the previous content to the same manifest or create another file. Now we want 2 Pods for the application, and inside the container definition, we specify a env section including environment variables to be available for our container. `MONGODB_URI` refers to a cluster-defined domain-name. This name will be available when we create the Services for internal Pod-to-Pod communication.

Now the frontend application will have a similar configuration:
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: post-app-deployment
  labels:
    app: post-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: post-app
  template:
    metadata:
      labels:
        app: post-app
    spec:
      containers:
      - name: post-app
        image: lfmtsml/post-app:latest
        ports:
        - containerPort: 80
        env:
          - name: POST_API_SVC_HOST
            value: "http://post-api-service.default.svc.cluster.local:8090"
```

Using a single manifest file, let's update our cluster:
```sh
$ kubectl apply -f deployments.yaml
deployment.apps/mongodb-deployment unchanged
deployment.apps/post-api-deployment created
deployment.apps/post-app-deployment created
```

Now we have 2 Pods for post-api-deployment and 1 Pod for post-app-deployment:
```sh
$ kubectl get pods,deploy,rs                
NAME                                      READY   STATUS    RESTARTS   AGE
pod/mongodb-deployment-66c66c656b-vmk8h   1/1     Running   0          21m
pod/post-api-deployment-86ff6bb47-6nf5v   1/1     Running   0          20s
pod/post-api-deployment-86ff6bb47-nk7jk   1/1     Running   0          20s
pod/post-app-deployment-549bc854-hjqjm    1/1     Running   0          20s

NAME                                  READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/mongodb-deployment    1/1     1            1           21m
deployment.apps/post-api-deployment   2/2     2            2           20s
deployment.apps/post-app-deployment   1/1     1            1           20s

NAME                                            DESIRED   CURRENT   READY   AGE
replicaset.apps/mongodb-deployment-66c66c656b   1         1         1       21m
replicaset.apps/post-api-deployment-86ff6bb47   2         2         2       20s
replicaset.apps/post-app-deployment-549bc854    1         1         1       20s
```

Let's take a look at our post-api-deployment:
```sh
$ kubectl describe deploy post-api-deployment
Name:                   post-api-deployment
Namespace:              default
CreationTimestamp:      Tue, 06 Feb 2024 17:03:43 -0300
Labels:                 app=post-api
Annotations:            deployment.kubernetes.io/revision: 1
Selector:               app=post-api
Replicas:               2 desired | 2 updated | 2 total | 2 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25% max unavailable, 25% max surge
Pod Template:
  Labels:  app=post-api
  Containers:
   post-api:
    Image:      lfmtsml/post-service:latest
    Port:       8090/TCP
    Host Port:  0/TCP
    Environment:
      MONGODB_URI:       mongodb://mongodb-service.default.svc.cluster.local:27017
      MONGODB_DATABASE:  post-database
    Mounts:              <none>
  Volumes:               <none>
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      True    MinimumReplicasAvailable
  Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  <none>
NewReplicaSet:   post-api-deployment-86ff6bb47 (2/2 replicas created)
Events:
  Type    Reason             Age   From                   Message
  ----    ------             ----  ----                   -------
  Normal  ScalingReplicaSet  11m   deployment-controller  Scaled up replica set post-api-deployment-86ff6bb47 to 2
```

The envvars we defined can be seen. But what if we want to centralize our configuration variables?

### Deploy a ConfigMap
Our Pods are running, but is a good idea to separate environment variables from the Deployments definition. For that, we can create a new object called ConfigMap. The ConfigMap stores key-value pairs which can be used by our Deployments to retrieve environment variables.

Please note that ConfigMaps should not be used to store Secrets. There's a specific Secrets object for that purpose.

Let's create another manifest file called `configs.yaml` with the content:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: post-app-config
data:
  mongodb_uri: "mongodb://mongodb-service.default.svc.cluster.local:27017"
  mongodb_database: "post-database"
  post_api_svc_host: "http://post-api-service.default.svc.cluster.local:8090"
```

We define a ConfigMap called post-app-config with the same environment variables we used earlier. Now we need to reference this ConfigMap on our deployments manifest. As such, update the env property of the Pod template declaration for post-api-deployment and post-app-deployment:
```yaml
<...>
env:
          - name: MONGODB_URI
            valueFrom:
              configMapKeyRef:
                name: post-app-config
                key: mongodb_uri
          - name: MONGODB_DATABASE
            valueFrom:
              configMapKeyRef:
                name: post-app-config
                key: mongodb_database
<...>
env:
          - name: POST_API_SVC_HOST
            valueFrom:
              configMapKeyRef:
                name: post-app-config
                key: post_api_svc_host
```

Let's deploy the new manifest and update the configuration of the old one:
```sh
$ kubectl apply -f configs.yaml && kubectl apply -f deployments.yaml
configmap/post-app-config created
deployment.apps/mongodb-deployment unchanged
deployment.apps/post-api-deployment configured
deployment.apps/post-app-deployment configured
```

We now have a new resource, which can be seen with:
```sh
$ kubectl get pod,rs,deploy,configmaps
NAME                                       READY   STATUS    RESTARTS   AGE
pod/mongodb-deployment-66c66c656b-vmk8h    1/1     Running   0          36m
pod/post-api-deployment-dccf699b9-krvcx    1/1     Running   0          61s
pod/post-api-deployment-dccf699b9-wfnmv    1/1     Running   0          65s
pod/post-app-deployment-548564b749-pvvk4   1/1     Running   0          65s

NAME                                             DESIRED   CURRENT   READY   AGE
replicaset.apps/mongodb-deployment-66c66c656b    1         1         1       36m
replicaset.apps/post-api-deployment-86ff6bb47    0         0         0       15m
replicaset.apps/post-api-deployment-dccf699b9    2         2         2       65s
replicaset.apps/post-app-deployment-548564b749   1         1         1       65s
replicaset.apps/post-app-deployment-549bc854     0         0         0       15m

NAME                                  READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/mongodb-deployment    1/1     1            1           36m
deployment.apps/post-api-deployment   2/2     2            2           15m
deployment.apps/post-app-deployment   1/1     1            1           15m

NAME                         DATA   AGE
configmap/kube-root-ca.crt   1      70m
configmap/post-app-config    3      65s
```

The old ReplicaSets were scaled down to 0, and new ReplicaSets for post-api-deployment and post-app-deployment were created.

Let's take a look at our post-api-deployment:
```sh
$ kubectl describe deploy post-api-deployment
Name:                   post-api-deployment
Namespace:              default
CreationTimestamp:      Tue, 06 Feb 2024 17:03:43 -0300
Labels:                 app=post-api
Annotations:            deployment.kubernetes.io/revision: 2
Selector:               app=post-api
Replicas:               2 desired | 2 updated | 2 total | 2 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25% max unavailable, 25% max surge
Pod Template:
  Labels:  app=post-api
  Containers:
   post-api:
    Image:      lfmtsml/post-service:latest
    Port:       8090/TCP
    Host Port:  0/TCP
    Environment:
      MONGODB_URI:       <set to the key 'mongodb_uri' of config map 'post-app-config'>       Optional: false
      MONGODB_DATABASE:  <set to the key 'mongodb_database' of config map 'post-app-config'>  Optional: false
    Mounts:              <none>
  Volumes:               <none>
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      True    MinimumReplicasAvailable
  Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  post-api-deployment-86ff6bb47 (0/0 replicas created)
NewReplicaSet:   post-api-deployment-dccf699b9 (2/2 replicas created)
Events:
  Type    Reason             Age    From                   Message
  ----    ------             ----   ----                   -------
  Normal  ScalingReplicaSet  17m    deployment-controller  Scaled up replica set post-api-deployment-86ff6bb47 to 2
  Normal  ScalingReplicaSet  2m57s  deployment-controller  Scaled up replica set post-api-deployment-dccf699b9 to 1
  Normal  ScalingReplicaSet  2m53s  deployment-controller  Scaled down replica set post-api-deployment-86ff6bb47 to 1 from 2
  Normal  ScalingReplicaSet  2m53s  deployment-controller  Scaled up replica set post-api-deployment-dccf699b9 to 2 from 1
  Normal  ScalingReplicaSet  2m50s  deployment-controller  Scaled down replica set post-api-deployment-86ff6bb47 to 0 from 1
```

The envvars from the ConfigMap are in use, and the scaling of replicas can be noted in the Events section.

Though our cluster is up and running, we cannot access our frontend application. Also, the internal DNS are not defined through services. Let's configure this.

### Deploy Services for internal communication between the Pods and for external-to-Pod communication (for frontend app)
To establish internal communication (Pod-to-Pod) using internal DNS addresses like the ones we defined as environment variables, we need to define Service objects.

A Service object exposes Pods through internal DNS or to the external world through the control plane node IP. The Service exposes a port that forwards traffic to a container destination port. We will explore two types of services:

- `ClusterIP`: a fixed internal IP address is defined and is accessible only inside the cluster. We will define this type of Service for our post-api-deployment and mongodb-deployment, as they do not need to be accessible from outside the cluster;
- `NodePort`: a fixed port is allocated on the control plane node, forwarding requests directly to the backend Pods. We will define this type of Service for the post-app-deployment, as we need to access the frontend application through the browser

As we are using Deployments, our Services will point to the Deployment objects through selectors to forward traffic to the application Pods.

Let's create a new manifest file called `services.yaml` with the following content:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: post-app-service
spec:
  type: NodePort
  selector:
    app: post-app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
      nodePort: 30008
---
apiVersion: v1
kind: Service
metadata:
  name: post-api-service
spec:
  type: ClusterIP
  selector:
    app: post-api
  ports:
    - protocol: TCP
      port: 8090
      targetPort: 8090
---
apiVersion: v1
kind: Service
metadata:
  name: mongodb-service
spec:
  type: ClusterIP
  selector:
    app: mongodb
  ports:
    - protocol: TCP
      port: 27017
      targetPort: 27017
---
```

For simplicity, we define ports and targetPorts with the same value for each Service. port indicates the port on which the Service is receiving connections, and targetPort is the destination port where the Pod is listening for requests.
For the post-app-service, we define both type: NodePort and nodePort: 30008. This value is not arbitrary; by default Kubernetes allocate a port in the range 30000-32767. If the selected port is not available, Kubernetes will error out, and we can select another one.

Let's deploy the manifest:
```sh
$ kubectl apply -f services.yaml
service/post-app-service created
service/post-api-service created
service/mongodb-service created
```

Our services can be seen as follows. Note that a kubernetes service already existed:
```sh
$ kubectl get svc
NAME               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
kubernetes         ClusterIP   10.96.0.1        <none>        443/TCP        92m
mongodb-service    ClusterIP   10.111.142.151   <none>        27017/TCP      26s
post-api-service   ClusterIP   10.104.19.99     <none>        8090/TCP       26s
post-app-service   NodePort    10.105.148.171   <none>        80:30008/TCP   27s
```

Our post-app-service can be accessed through a minikube tunnel. Execute the following command to open the tunnel:
```sh
$ minikube service post-app-service --url
http://127.0.0.1:38461
❗  Because you are using a Docker driver on linux, the terminal needs to be open to run it.
```


We can interact with our application.

### Remove Pods and check ReplicaSet behavior
Our application is working as intended. Let's explore more Kubernetes functionality. For example, what happens if we delete a Pod?

First, we need to retrieve the Pod name for the post-app-deployment:
```sh
$ kubectl get pod
NAME                                   READY   STATUS    RESTARTS   AGE
mongodb-deployment-66c66c656b-vmk8h    1/1     Running   0          63m
post-api-deployment-dccf699b9-krvcx    1/1     Running   0          28m
post-api-deployment-dccf699b9-wfnmv    1/1     Running   0          28m
post-app-deployment-548564b749-pvvk4   1/1     Running   0          28m
```

Let's delete the Pod, then verify that the Pod is gone:
```sh
$ kubectl delete pod post-app-deployment-548564b749-pvvk4  && kubectl get pod
pod "post-app-deployment-548564b749-pvvk4" deleted
NAME                                   READY   STATUS    RESTARTS   AGE
mongodb-deployment-66c66c656b-vmk8h    1/1     Running   0          66m
post-api-deployment-dccf699b9-krvcx    1/1     Running   0          31m
post-api-deployment-dccf699b9-wfnmv    1/1     Running   0          31m
post-app-deployment-548564b749-cshnp   1/1     Running   0          2s
```

Please note that a new Pod appeared (started 2s ago). Now let's take a closer look at our ReplicaSet:
```sh
$ kubectl get rs          
NAME                             DESIRED   CURRENT   READY   AGE
mongodb-deployment-66c66c656b    1         1         1       68m
post-api-deployment-86ff6bb47    0         0         0       47m
post-api-deployment-dccf699b9    2         2         2       32m
post-app-deployment-548564b749   1         1         1       32m
post-app-deployment-549bc854     0         0         0       47m

$ kubectl describe rs post-app-deployment-548564b749
Name:           post-app-deployment-548564b749
Namespace:      default
Selector:       app=post-app,pod-template-hash=548564b749
Labels:         app=post-app
                pod-template-hash=548564b749
Annotations:    deployment.kubernetes.io/desired-replicas: 1
                deployment.kubernetes.io/max-replicas: 2
                deployment.kubernetes.io/revision: 2
Controlled By:  Deployment/post-app-deployment
Replicas:       1 current / 1 desired
Pods Status:    1 Running / 0 Waiting / 0 Succeeded / 0 Failed
Pod Template:
  Labels:  app=post-app
           pod-template-hash=548564b749
  Containers:
   post-app:
    Image:      lfmtsml/post-app:latest
    Port:       80/TCP
    Host Port:  0/TCP
    Environment:
      POST_API_SVC_HOST:  <set to the key 'post_api_svc_host' of config map 'post-app-config'>  Optional: false
    Mounts:               <none>
  Volumes:                <none>
Events:
  Type    Reason            Age    From                   Message
  ----    ------            ----   ----                   -------
  Normal  SuccessfulCreate  33m    replicaset-controller  Created pod: post-app-deployment-548564b749-pvvk4
  Normal  SuccessfulCreate  2m16s  replicaset-controller  Created pod: post-app-deployment-548564b749-cshnp
```

As we can see, the ReplicaSet automatically recognized that the old Pod gone downhill and created another one.

### Rollout changes to the cluster
We've made some changes to our deployments. Let's take a look at the rollout history for the post-app-deployment:
```sh
$ kubectl rollout history deploy post-app-deployment
deployment.apps/post-app-deployment 
REVISION  CHANGE-CAUSE
1         <none>
2         <none>
```

As we can see, we have two revisions. Let's confirm which revision is in use:
```sh
$ kubectl describe deploy post-app-deployment | grep revision -A 2 -B 2
CreationTimestamp:      Tue, 06 Feb 2024 17:03:43 -0300
Labels:                 app=post-app
Annotations:            deployment.kubernetes.io/revision: 2
Selector:               app=post-app
Replicas:               1 desired | 1 updated | 1 total | 1 available | 0 unavailable
```

So, revision 2 is currently available and its working normally. Let's make a change to the deployment and rollout. We can change the container image tag for the Pod template. We will define a non-existent tag just to break our application.
```sh
$ kubectl set image deploy post-app-deployment post-app=lfmtsml/post-app:1.0.0
deployment.apps/post-app-deployment image updated
$ kubectl get pod
NAME                                   READY   STATUS             RESTARTS   AGE
mongodb-deployment-66c66c656b-vmk8h    1/1     Running            0          83m
post-api-deployment-dccf699b9-krvcx    1/1     Running            0          47m
post-api-deployment-dccf699b9-wfnmv    1/1     Running            0          47m
post-app-deployment-548564b749-cshnp   1/1     Running            0          16m
post-app-deployment-68546bb485-9ncdq   0/1     ImagePullBackOff   0          87s
```

The new Pod cannot spin, because the tag does not exists.

### Rollback changes to the cluster
Our application is down, because our last change introduced a problem. We can revert our changes to a previous revision. As such, we can restore our cluster to a valid configuration.
```sh
$ kubectl rollout undo deployment post-app-deployment --to-revision=1
deployment.apps/post-app-deployment rolled back
$ kubectl rollout history deploy post-app-deployment
deployment.apps/post-app-deployment 
REVISION  CHANGE-CAUSE
2         <none>
3         <none>
4         <none>
$ kubectl get pod
mongodb-deployment-66c66c656b-vmk8h   1/1     Running   0          88m
post-api-deployment-dccf699b9-krvcx   1/1     Running   0          52m
post-api-deployment-dccf699b9-wfnmv   1/1     Running   0          53m
post-app-deployment-549bc854-xjs2c    1/1     Running   0          13s

```

Our Pod is available again, and the application is working as intended.

### Destroy everything
Because minikube guarantees the isolation of the cluster environment inside our machine, deleting all traces of it is extremely easy. We can stop minikube and delete the cluster with:
```sh
minikube stop && minikube delete
```
