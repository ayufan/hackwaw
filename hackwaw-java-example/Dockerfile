FROM java:8

RUN apt-get update && apt-get install -y --no-install-recommends maven

WORKDIR /code

ADD pom.xml /code/pom.xml
RUN ["mvn", "dependency:resolve"]

ADD src /code/src
RUN mvn package -DskipTests

CMD java -Djava.security.egd=file:/dev/./urandom -jar target/disrupt-twitter-slack*.jar

