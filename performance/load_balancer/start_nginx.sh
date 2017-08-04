#!/usr/bin/env bash

pkill -f nginx

nginx -c load_balancer/nginx_config.conf -p $PWD