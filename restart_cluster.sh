#!/bin/bash

minikube stop && minikube delete
minikube start
minikube kubectl -- apply -f infra/
minikube kubectl -- get deploy,rs,pods,services,configmaps