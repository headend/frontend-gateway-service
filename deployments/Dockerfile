FROM centos:7.2.1511

RUN mkdir /app
WORKDIR /app
ADD backend-gateway /app/backend-gateway
RUN rm -f etc/localtime; ln -sf /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime
EXPOSE 8888
# Run
CMD ["/app/backend-gateway"]