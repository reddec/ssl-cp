#!/usr/bin/env bash
set -e -o pipefail
ROOT=$(dirname "$BASH_SOURCE")
cd "$ROOT"

mkdir -p "/opt/stunnel/{{cert.common_name}}/server"
cp -rv . "/opt/stunnel/{{cert.common_name}}/server/"
rm -rf "/opt/stunnel/{{cert.common_name}}/server/"*.sh
systemctl enable "/opt/stunnel/{{cert.common_name}}/server/stunnel-{{cert.common_name}}-server.service"
systemctl start stunnel-{{cert.common_name}}-server.service