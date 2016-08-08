FROM scratch

COPY spider /

ENTRYPOINT ["/spider"]
