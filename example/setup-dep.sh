#!/bin/bash

#Parameters

ClUSTER_NAME="${1:-ksqlcluster-operator}"
NAMESPACE="${2:-demo}"
DO_DEPLOY="${3:-false}"

#Constants
DEP_CLUSTER_NAME="demo"

INCOMING_TOPIC="com.topic.in.order_created"

PRODUCER_SERVICE="producer-service"

get_service_dir() {
    # $1: service name
    echo "$(cd apps/$1 && pwd)"
}

##### KAKFA SETUP
setup_kafka_cluster() {
    echo "Check existence of kafka helm repo"
    KAFKA_REPO_EXISTENCE=$(helm repo list | grep "https://confluentinc.github.io/cp-helm-charts")

    if [ -z "${KAFKA_REPO_EXISTENCE}" ]; then 
        echo "kafka helm chart repo is being added"
        helm repo add confluentinc https://confluentinc.github.io/cp-helm-charts/
    fi

    echo "Kafka cluster is being deployed"
    helm install -f ./config/confluent-value.yaml \
        $DEP_CLUSTER_NAME confluentinc/cp-helm-charts || true
}

wait_for_Kafka() {
    echo "Waiting for kafka cluster to be deployed successfully"
    kubectl wait --for=condition=available --timeout=-1s deployment/$DEP_CLUSTER_NAME-cp-control-center # to prevent topic to be generated before consumer deployed
    declare -a num=("0" "1" "2")
    for i in "${num[@]}"
    do
        kubectl wait --for=condition=ready --timeout=-1s pod/$DEP_CLUSTER_NAME-cp-kafka-${i}
        kubectl wait --for=condition=ready --timeout=-1s pod/$DEP_CLUSTER_NAME-cp-zookeeper-${i}
    done
}

create_topic() {
    TOPIC=$1
    echo "Topic($TOPIC) is being created"
    kubectl exec -c cp-kafka-broker -it $DEP_CLUSTER_NAME-cp-kafka-0 -- /bin/bash /usr/bin/kafka-topics --create --zookeeper $DEP_CLUSTER_NAME-cp-zookeeper:2181 --topic $TOPIC --partitions 3 --replication-factor 1
}

##### SERVICE SETUP
delete() {
    # $1: service name
    SERVICE=$1
    DIR=$(get_service_dir "$1")

    kubectl kustomize $DIR/k8s | kubectl delete -f -
}

build() {
    # $1: service name
    SERVICE=$1
    DIR=$(get_service_dir "$1")
    echo "build $SERVICE"
    docker rm -f ${SERVICE}
    docker rmi $(docker images | grep "${SERVICE}") || true
    docker build ${DIR}/. -t ${SERVICE}
    kind load docker-image --name ${ClUSTER_NAME} ${SERVICE}
}

deploy() {
    # $1: service name, $2: bootstrap, $3: topic name 
    echo "deploy customer service"
    DIR=$(get_service_dir "$1")
    gsed -e "s,VALUE_KAFKA_BOOTSTRAP,$2,g" \
        -e "s,VALUE_SCHEMA_REGISTRY_SERVER,$3,g"  \
        -e "s,VALUE_KAFKA_TOPIC,$4,g" \
        $DIR/k8s/deployment.yml.template > $DIR/k8s/deployment.yml
    
    kubectl kustomize $DIR/k8s | kubectl apply -f -
}


##### Delegation takes place
if [ "$DO_DEPLOY" = true ]; then
    delete "$PRODUCER_SERVICE"
    build "$PRODUCER_SERVICE"

    deploy "$PRODUCER_SERVICE" "$DEP_CLUSTER_NAME-cp-kafka:9092" \
        "http://$DEP_CLUSTER_NAME-cp-schema-registry:8081" "$INCOMING_TOPIC"
else
    helm repo update
    setup_kafka_cluster

    CONSUMER_EXISTENCE="$(kubectl get deployment --no-headers -o custom-columns=":metadata.name" | grep "$CONSUMER_SERVICE")"
    if [ -z "${CONSUMER_EXISTENCE}" ]; then
        wait_for_Kafka
        sleep 10
        create_topic "$INCOMING_TOPIC"
        sleep 2
    fi
fi

echo "finished"