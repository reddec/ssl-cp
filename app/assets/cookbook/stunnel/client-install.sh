#!/usr/bin/env bash
set -e -o pipefail
ROOT=$(dirname "$BASH_SOURCE")
cd "$ROOT"

mkdir -p "/opt/stunnel/{{cert.common_name}}/client"
cp -rv . "/opt/stunnel/{{cert.common_name}}/client/"
rm -rf "/opt/stunnel/{{cert.common_name}}/client/"*.sh
systemctl enable "/opt/stunnel/{{cert.common_name}}/client/stunnel-{{cert.common_name}}-client.service"
systemctl start stunnel-{{cert.common_name}}-client.service