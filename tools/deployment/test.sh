kubectl apply -f nginx.yaml
sleep 10
kubectl delete -f nginx.yaml
number_of_pods=$(kubectl get pods --output name -n knative-functions | wc -l)
while [ "$number_of_pods" -ne 0 ]; do
    number_of_pods=$(kubectl get pods --output name -n knative-functions| wc -l)
    sleep 10
done
kubectl apply -f nginx-1.yaml
kubectl delete -f nginx-1.yaml




number_of_non_running_pods=$(kubectl get pods -n knative-functions --field-selector=status.phase!=Running --output name | wc -l)
while [ "$number_of_non_running_pods" -ne 0 ]; do
    number_of_non_running_pods=$(kubectl get pods -n knative-functions --field-selector=status.phase!=Running --output name | wc -l)
    sleep 10
done
number_of_running_pods=$(kubectl get pods -n knative-functions --field-selector=status.phase==Running --output name | wc -l)
echo $number_of_running_pods