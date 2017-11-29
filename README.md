# VNF HTTP rest API

VNF HTTP rest API with basic operations. This VNF is meant for the OAI artefact.

## Operations

The basic operations defined are below
 ```console
GET   /status
POST  /configure
    Headers: "Content-Type: application/json"
    Body: '{
            "param": "value"
          }'
POST  /start
DELETE  /stop/{vnfID}
GET   /status/{vnfID}
 ```

More operations and configuration parameters may come later

## Information

More about authors, contribution, and license in the docs

![api](./docs/api.png)
