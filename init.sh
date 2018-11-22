#! /bin/bash

if [ ! -e ./vendor/github.com ]; then
	$GOPATH/bin/govendor sync
fi

if [ ! -e ./.env ]; then
    mv ./.env.example ./.env
    echo 'Please modify env file.'
fi

echo 'Initialization finished.'
