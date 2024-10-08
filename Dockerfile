FROM scratch

EXPOSE 8080 8085

ENTRYPOINT ["/usr/bin/terrapi"]
CMD ["server"]


COPY bin/terrapi /usr/bin/terrapi