FROM java:alpine
MAINTAINER Allan Kerr

EXPOSE 8080
ADD target/api.jar /usr/share/api.jar
ADD Google-Service-Account.json Google-Service-Account.json
ENV GOOGLE_APPLICATION_CREDENTIALS=/Google-Service-Account.json
ENTRYPOINT ["/usr/bin/java", "-jar", "/usr/share/api.jar"]
