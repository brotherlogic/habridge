kubectl delete secret hatoken -n hatoken
kubectl create secret generic hatoken --from-literal hatoken=$1 -n habridge
