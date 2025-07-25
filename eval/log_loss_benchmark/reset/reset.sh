#!/bin/sh

kubectl delete --force -f ../../experimental_environment_pod.yaml
kubectl apply -f ../../experimental_environment_pod.yaml
sleep 5

kubectl cp ../ls_repeater/ls_repeater default/ubuntu:/
kubectl cp ../ls_repeater/for.sh default/ubuntu:/
