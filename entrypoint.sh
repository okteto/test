#!/bin/sh

set -e

name=$1
namespace=$2
file=$3
deploy=$4
no_cache=$5
variables=$6
timeout=$7
tests=$8
log_level=$9

IFS=$'\t\n '

if [ -n "$OKTETO_CA_CERT" ]; then
   echo "Custom certificate is provided"
   echo "$OKTETO_CA_CERT" > /usr/local/share/ca-certificates/okteto_ca_cert.crt
   update-ca-certificates
fi

params=""

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

if [ "$no_cache" = "true" ]; then
      params="$params --no-cache"
fi

variable_params=""
if [ ! -z "${variables}" ]; then
  for ARG in $(echo "${variables}" | tr ',' '\n'); do
  
    variable_params="${variable_params} --var \"${ARG}\""
  done

  params="${params} ${variable_params}"
fi

github_env_vars=$(env | grep '^GITHUB_')
github_params=""

for VAR in $github_env_vars; do
   VAR_NAME=$(echo $VAR | cut -d= -f1)
   VAR_VALUE=$(echo $VAR | cut -d= -f2 | tr -d ' ')
   github_params="$github_params --var=${VAR_NAME}=\"${VAR_VALUE}\""
done
params="$params $github_params"

if [ -n "$timeout" ]; then
   params="$params --timeout $timeout"
fi

params="$params $tests"

if [ ! -z "$log_level" ]; then
  log_level="--log-level ${log_level}"
fi

# https://docs.github.com/en/actions/monitoring-and-troubleshooting-workflows/enabling-debug-logging
# https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
if [ "${RUNNER_DEBUG}" = "1" ]; then
  log_level="--log-level debug"
fi

echo running: okteto test $log_level "$params"
# shellcheck disable=SC2086
okteto test $log_level $params