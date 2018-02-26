FROM alpine:3.6
RUN apk add --no-cache ca-certificates
RUN adduser -D -g "" oluu && mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY /app.o /home/oluu/
WORKDIR /home/oluu
RUN chown -R oluu:oluu /home/oluu && chmod 700 /home/oluu/app.o
USER oluu
CMD ["./app.o"] 
