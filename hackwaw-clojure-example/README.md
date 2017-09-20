# hackwaw-clojure-example

Wzór aplikacji Clojure wykorzystywanej na Hackwaw Disrupt - Warszawa 2 kwietnia

## Development

Do pracy z aplikacją używane jest narzędzie [Leiningen](http://leiningen.org/).
Kod edytujemy w ulubionym edytorze lub w Idei w wtyczką [Cursive](https://cursive-ide.com/).

Najlepiej lokalny development oprzeć na pracy w REPL:

```sh
lein repl
```

Po wystartowaniu uruchamiamy aplikację komendą `go`:

```clojure
user=> (go)
:started
```

Przy pierwszym uruchomieniu musimy wykonać migrację bazy danych by utworzyć schemat używany przez aplikację:

```clojure
user=> (migrate)
Applying 001-tweets
=> nil

```

Po starcie odpalić się serwer pod <http://localhost:3000>. Dokumentacja Swaggera dostępna jest pod <http://localhost:3000/api-docs>

Gdy robimy jakieś zmiany możemy ewaluować funkcje bezpośrednio w REPL. Aby przeładować aplikację z aktualnym kodem używamy `reset`:

```clojure
user=> (reset)
:reloading (...)
:resumed
```

### Testowanie 

Testy w REPL uruchamiamy komendą `test`:

```clojure
user=> (test)
...
```

Lub z shell'a:

```sh
lein test
```

### Uruchomienie w Dockerze

* Zbuduj obraz: `docker build -t disrupt-clojure-app .`
* Uruchom obraz: `docker run -d disrupt-clojure-app`

# Zmienne środowiskowe
Obraz powinien pobierać adresu URI do twittera i slacka ze zmiennej środowiskowej. Jest to wymagane w celach testowania obrazów. Domyślnie jest ustawione na heroku. 
`docker run -d -e TWITTER_URL=https://hackwaw-twitter-proxy.herokuapp.com -e SLACK_URL=https://hackwaw-slack-proxy.herokuapp.com disrupt-clojure-app`