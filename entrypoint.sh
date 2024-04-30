#!/bin/sh

set -e

name=$1
namespace=$2
file=$3
deploy=$4
nocache=$5
variables=$6
timeout=$7
tests="${@:8}"

if [ -n "$OKTETO_CA_CERT" ]; then
   echo "Custom certificate is provided"
   echo "$OKTETO_CA_CERT" > /usr/local/share/ca-certificates/okteto_ca_cert.crt
   update-ca-certificates
fi

command="test"

params="--progress plain"

if [ -n "$name" ]; then
   params="$params --name $name"
fi

if [ -n "$namespace" ]; then
   params="$params --namespace $namespace"
fi

if [ -n "$file" ]; then
   params="$params -f $file"
fi

if [ "$deploy" = "true" ]; then
      params="$params --deploy"
fi

if [ "$nocache" = "true" ]; then
      params="$params --no-cache"
fi

variable_params=""
if [ ! -z "${variables}" ]; then
  for ARG in $(echo "${variables}" | tr ',' '\n'); do
    variable_params="${variable_params} --var ${ARG}"
  done

  params="${params} ${variable_params}"
fi


if [ -n "$timeout" ]; then
   params="$params --timeout $timeout"
fi

params="$params $tests"

echo running: okteto "$command" "$params"
# shellcheck disable=SC2086
okteto $command $params