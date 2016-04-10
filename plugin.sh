#!/bin/bash

echo '
{
  "api.4game.com": {
    "/": [
      {
        "staging": "stage",
        "address":"127.0.0.1",
        "port":8081
      },
      {
        "staging": "live",
        "address":"127.0.0.1",
        "port":8082
      }
    ],
    "/rumble/": [
      {
        "address":"127.0.0.1",
        "port":8083
      }
    ]
  },
  "api.kidzite.qa": {
    "/": [
      {
        "address":"127.0.0.1",
        "port":8082
      }
    ]
  }
}
'
