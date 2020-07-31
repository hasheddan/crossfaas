# crossfaas

`crossfaas` brings together [Crossplane](https://crossplane.io/),
[OpenFaaS](https://www.openfaas.com/), and [ArgoCD](https://argoproj.github.io/)
to provide a GitOps-driven serverless platform on
[Kubernetes](https://kubernetes.io/).

## Setup

Install ArgoCD.
```
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

Expose ArgoCD server.
```
kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'
```

Install OpenFaaS using Arkade.
```
curl -SLsf https://dl.get-arkade.dev/ | sudo sh

arkade install openfaas --load-balancer
```

Install Crossplane and any providers.
```
kubectl create namespace crossplane-system

helm repo add crossplane-alpha https://charts.crossplane.io/alpha
helm install crossplane --namespace crossplane-system crossplane-alpha/crossplane --set clusterPackages.gcp.deploy=true --set clusterPackages.gcp.version=master --devel
```

Login to ArgoCD.
```
ARGO_URL=$(kubectl get svc -n argocd argocd-server | awk 'FNR == 2 {print $4}')
ARGO_PASSWORD=$(kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d'/' -f 2)

argocd login $ARGO_URL
```

Login to OpenFaaS.
```
export OPENFAAS_URL=$(kubectl get svc -o wide gateway-external -n openfaas | awk 'FNR == 2 {print $4}'):8080

OPENFAAS_PASSWORD=$(kubectl get secret -n openfaas basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 --decode; echo)

echo -n $OPENFAAS_PASSWORD | faas-cli login --username admin --password-stdin
```

Build functions.
```
faas-cli build -f crossfaas.yaml --build-arg GO111MODULE=on --parallel 3
faas-cli push -f crossfaas.yaml
```