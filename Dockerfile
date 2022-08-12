FROM 766589103339.dkr.ecr.us-west-2.amazonaws.com/go:0.1.2
ARG service_name
ADD bin/server  /opt/bin/server
ENV SERVICE_NAME ${service_name}
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/opt/bin/server"]
