FROM ubuntu:16.04

LABEL atom="vnf-api-golang"
LABEL REPO="https://gitlab.forge.orange-labs.fr/lucy/vnf-api-golang"

ENV PATH=$PATH:/opt/vnf-api-golang/bin

WORKDIR /opt/vnf-api-golang/bin

COPY bin/vnf-api-golang.py /opt/vnf-api-golang/bin/
RUN chmod +x /opt/vnf-api-golang/bin/vnf-api-golang.py

CMD ["python", "/opt/vnf-api-golang/bin/vnf-api-golang.py"]
