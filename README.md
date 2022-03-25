# 12-factor-app
Course demo for 12-factor, Docker, Kubernetes concepts

The sample Golang application provides three functions. No Golang programming is required
- /hello  # Echo back hello
- /headers   # Display the HTTP headers used
- /dowork # Wastes time by generating millions of random numbers so we can play with Kubernetes scaling

## Setup
It's assumed you are running in GCP Cloudshell as it provides builtin Go, Docker, and Kubernetes support.  
Your cloudshell should be associated with a fully authorized project and located in your preferred region.  
The following commands are provided for your conveniance.
```
gcloud projects create YOUR_PROJECT_ID
gcloud config set project YOUR_PROJECT_ID
gcloud services list --available
gcloud services enable container.googleapis.com
gcloud config set compute/zone us-west1-a
```
[Docker cheat sheet](https://dockerlabs.collabnix.com/docker/cheatsheet/)  
[Kubenetes cheat sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)  

## Lesson #1
Review Docker Lab's Github: [12 Factor Application](https://github.com/docker/labs/tree/master/12factor)

***Extra Credit***
Does Docker help drive all 12 factors?
Which factor is Docker most/less helpful?

## Lesson #2
Experiment with a builder style Docker deployment.  

The build
```
git clone https://github.com/kengraf/12-factor-app.git
cd 12-factor-app
docker build -t http .
```

Validation: this is a good time to experiment with commands from the Docker cheat sheet.  
```
docker images # 3 images: goland pulled from DockerHub, a unnamed developer image, and the scratch based http image
docker images --all # show additional image layers created during the build
docker ps -a # No containers until you run something
```

Run
```
# The built application, /hello and /header return immediately, /dowork should take a minute or two
docker run -d -p 80:8090 http
curl http://localhost/hello
curl http://localhost/headers
curl http://localhost/dowork
```
***Extra Credit***  
Why not seperate building the **developer** and **deployer** images into different Dockerfiles?
Can you retrieve the "app" and run it locally?  
Is the "app" located in the same directory for both the **developer** and **deployer** images?  Can you prove it?  
Can you explain the significant size difference in the images?  
Which image can you debug by getting shell access? docker run -it IMAGE SHELL  
Is there any way to debug the "http" container?  
What happens to the images and containers when the cloudshell session temrinates?  

## Lesson #3
Now that we have a deploy focused Docker image, it is time to learn about scaling it with Kubernete.  
A more complete deployment back & front endis provided in Google's Kubenetes examples: [GCP demo](https://cloud.google.com/kubernetes-engine/docs/tutorials/guestbook)  

Push image to repo
You will a docker hub account http://hub.docker.com.  
```
docker login --username YOUR_NAME
docker tag http YOUR_NAME/http:v1
docker push YOUR_NAME/http:v1

```

Create a new cluster for the deployment
```
gcloud container clusters create http --num-nodes=4
gcloud container clusters list
gcloud container clusters describe http

```

Deploy frontend
Note: This assumes the http image generated above is still available locally
```
kubectl apply -f deployment.yaml
kubectl get pods -l app=http -l tier=frontend
kubectl apply -f service.yaml
  or
kubectl expose deployment http --name=http --type=LoadBalancer --port 80 --target-port 8090

kubectl apply -f scaling.yaml
kubectl autoscale deployment http --cpu-percent=50 --min=1 --max=10
kubectl get service http
kubectl scale deployment http --replicas=5
kubectl get pods
```

Commands to check status
```
kubectl get deployments
kubectl get services
kubectl rollout status deployment/http
kubectl get rs
kubectl get pods --show-labels
docker run --rm -it -v ~/.kube/config:/root/.kube/config quay.io/derailed/k9s
```

```
# Run this in a separate terminal
# so that the load generation continues and you can carry on with the rest of the steps
kubectl run -i --tty load-generator --rm --image=busybox --restart=Never -- /bin/sh -c "while sleep 1; do wget -b -q -O- http://EXTERNAL_LOADBALANCER_IP/hello; done"
```


Clean up Kubernetes
```
kubectl delete service http
kubectl delete deployment http
gcloud compute forwarding-rules list
gcloud container clusters delete http

```

## Clean up Docker
```
docker rm -vf $(docker ps -aq)  
docker rmi -f $(docker images -aq)

```

