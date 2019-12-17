# Dotscience Operator 

Dotscience Operator deploys dotscience deployers that can in turn deploy TF, sklearn models and configure services & ingresses. 

## Local deployment to OpenShift

1. Generate csv from the manifests (operator, role_binding, role, service_account and crds/):

  ```
  make gen-csv
  ```

2. Manually edit the last one to add `displayName: Dotscience Operator` in the spec.customresourcedefinitions.owned (TODO: somehow figure it out how to do this automatically?)
3. Probably just remove `replaces` from the clusterserviceversion.yaml file (it doesn't like it when it's there but there's no published previous version)
3. Fill env variables in scripts/operator_courier_env.sh
4. Push operator bundle to your quay app repo:

  ```
  make push-app
  ```

5. Create operator source:

  ```
  kubectl apply -f scripts/operator_source.yaml
  ```

6. Check the source:


  ```
  kubectl get operatorsource test-operators -n openshift-marketplace
  ```

  This should return something like:

  ```
  NAME             TYPE          ENDPOINT              REGISTRY   DISPLAYNAME   PUBLISHER   STATUS      MESSAGE                                       AGE
test-operators   appregistry   https://quay.io/cnr   rusenask                             Succeeded   The object has been successfully reconciled   6m17s

  ```


  Additionally, CatalogSource will be created:

  ```
  kubectl get catalogsource -n openshift-marketplace 
NAME                  DISPLAY               TYPE   PUBLISHER   AGE
certified-operators   Certified Operators   grpc   Red Hat     24d
community-operators   Community Operators   grpc   Red Hat     24d
redhat-operators      Red Hat Operators     grpc   Red Hat     24d
test-operators                              grpc               79s
  ```

7. View available operators:

  ```
  kubectl get opsrc test-operators -o=custom-columns=NAME:.metadata.name,PACKAGES:.status.packages -n openshift-marketplace 
NAME             PACKAGES
test-operators   dotscience-operator
  ```

8. Create operator group. An OperatorGroup is used to denote which namespaces your Operator should be watching. It must exist in the namespace where your operator should be deployed:

  ```
  kubectl apply -f scripts/operator_group.yaml
  ```

9. The last piece ties together all of the previous steps. A Subscription is created to the operator:

  ```
  kubectl apply -f scripts/operator_subscription.yaml
  ```

10. View operator health

  ```
  kubectl get clusterserviceversion -n openshift-marketplace
  ```

11. Once you have it, you can iterate rapidly by

  1. uninstalling the operator from operatorhub
  2. rerunning the operator-courier push, just with an incremented release number (`export PACKAGE_VERSION=0.2.0`)
  3. deleting the pod in the openshift-marketplace namespace
  4. installing the new version when it shows up


To delete operator you will need to remove group, source, subscription and the actual deployment. Otherwise OpenShift doesn't appear to be picking up a new version:

  ```
  make delete-operator
  ```

To install a new version:

  ```
  make install-operator
  ```

Useful links:

- More info on developing operators: https://github.com/operator-framework/community-operators/blob/master/docs/testing-operators.md
- Operator SDK CLI ref: https://github.com/operator-framework/operator-sdk/blob/master/doc/sdk-cli-reference.md
- crc: https://code-ready.github.io/crc/
- Running crc on Ubuntu: https://labs.consol.de/devops/linux/2019/11/29/codeready-containers-on-ubuntu.html
- Operator framework repo: https://github.com/operator-framework/operator-sdk