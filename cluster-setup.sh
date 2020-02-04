#!/usr/bin/env bash

#Copy kubeconfig to home dir
cp skripsi-cluster-kubeconfig.yaml ~/.kube/config

#Setting up metric server
kubectl apply -f metrics-server/deploy/1.8+/

#Setting up ingress nginx
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.26.2/deploy/static/mandatory.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.26.2/deploy/static/provider/cloud-generic.yaml
