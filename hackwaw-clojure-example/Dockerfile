FROM cstpdk/java8-clojure-leiningen

ADD . /work

WORKDIR /work

ENV LEIN_ROOT true

RUN lein uberjar

CMD java -jar target/uberjar/hackwaw-clojure-example-0.1.0-SNAPSHOT-standalone.jar
