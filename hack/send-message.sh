#!/bin/bash
#: sends message to the bot emulating telegram api.
if [ -z "${1}" ]; then
  #: don't forget to change user identifier, as far as message.json has THE fake one
  _file=pkg/telegram/resources/message.json
else
  _file=$1
fi
curl -k -X POST -d @"${_file}" https://dnq.blacklibrary.ru:8443/mast
