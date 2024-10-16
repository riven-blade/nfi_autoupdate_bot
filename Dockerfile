FROM alpine:3.14
WORKDIR /app
COPY check_nfi_upate ./check_nfi_upate
RUN chmod +x ./check_nfi_upate
CMD ["./check_nfi_upate"]
