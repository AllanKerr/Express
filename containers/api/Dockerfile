FROM java:alpine
MAINTAINER Allan Kerr

EXPOSE 8080
ADD target/api.jar /usr/share/api.jar
ENTRYPOINT ["/usr/bin/java", "-jar", "/usr/share/api.jar"]
