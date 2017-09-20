FROM 1science/sbt

RUN mkdir app
COPY /build.sbt /app/
COPY /project /app/project
COPY /src /app/src

WORKDIR /app

RUN sbt assembly

CMD java -jar target/scala-2.11/hackwaw.jar

