#!/bin/bash

# Usage: 
# run all scenarios from single feature file: ./deployments/k8s_bdd_test.bash features/bob/add_a_class_member.feature
# run all scenarios from all services: ./deployments/k8s_bdd_test.bash
# run all scenarios from single service: ./deployments/k8s_bdd_test.bash user

paths=$1
working_paths=features/

if [[ $paths == "" ]]
    then
        paths=.
else
    if [[ $paths == $working_paths* ]]
        then paths=${paths#$working_paths}
    fi
fi

cd $working_paths
exec go run main.go -godog.paths $paths