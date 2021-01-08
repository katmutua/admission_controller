### Local Cluster

We can use minikube to bring up a local kubernetes cluster

```
brew install minikube
brew cask remove minikube
brew link minikube


minikube start
```
To confirm that you are pointing to the correct cluster run

````
  kubectl config get-contexts # this should point to minikube
```
