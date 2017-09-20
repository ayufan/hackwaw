FROM java:8

RUN mkdir workplace
COPY /pom.xml /workplace/
COPY /src /workplace/src/

RUN apt-get update && apt-get install -y --no-install-recommends maven

WORKDIR /workplace

RUN mvn clean install
CMD java -jar target/disrupt-twitter-proxy*.jar

