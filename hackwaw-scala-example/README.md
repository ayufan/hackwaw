# Aplikacja - Scala

Wzór aplikacji w Scali (Akka HTTP) wykorzystywanej na Hackwaw Disrutp - Warszawa 2 kwietnia.

## Uruchamianie

### SBT

Do uruchomienia potrzebne jest SBT (i Java):

```
$ sbt
> ~re-start
```

Przed uruchomieniem upewnij się czy masz wyeksportowane odpowiednie zmienne środowiskowe:

```
export TWITTER_URL=https://hackwaw-twitter-proxy.herokuapp.com
export SLACK_URL=https://hackwaw-slack-proxy.herokuapp.com
```

### Docker

Zbuduj i uruchom obraz:

```
$ docker build -t hackwaw-scala-example .
$ docker run -p 8080:8080 -e TWITTER_URL=https://hackwaw-twitter-proxy.herokuapp.com -e SLACK_URL=https://hackwaw-slack-proxy.herokuapp.com hackwaw-scala-example
```

### Pytania?

W razie pytań lub wątpliwości szukaj Łukasza Sowy - <contact@luksow.com>