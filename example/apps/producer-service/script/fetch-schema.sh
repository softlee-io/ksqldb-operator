#!/bin/bash

schema_registry_url_default="http://localhost:8081"

schema="$1"
schema_registry_url="${2:-$schema_registry_url_default}"
api_key="$3"
api_secret="$4"

validate() {
    # $1: name, $2: value, $3: default value
    if [ -z "$2" ]; then 
        echo "[Error] $1 is null"
        exit 1
    fi

    if [ -n "$3" ]; then
        echo "[Warn] $2 default value ($2) "
    fi
}

validate "schema_registry_url" "$schema_registry_url" schema_registry_url_default
validate "schema" "$schema"

if [ -z "$api_key" ]; then 
    curl -X GET $schema_registry_url/subjects/Kafka-value/versions/1
else 
    curl -X GET --user <schema-registry-api-key>:<schema-registry-api-secret> \
        $schema_registry_url/subjects/Kafka-value/versions/1
fi
