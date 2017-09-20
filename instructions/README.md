# Intro

Witaj na Hackwaw Disrupt. Ta instrukcja pomoże uczestnikom Hackwaw Disrupt odnaleźć się w tym co będziemy robili podczas tego Hackatonu.

## Cel hackatonu

Naszym zadaniem będzie zbudowanie aplikacji, która w swoim działaniu będzie kuloodporna.
Nie zniszczą jej niedziałające serwisy zewnętrzne, nieprzyjemne środowisko itp.

## Aplikacja

Aplikacja realizuje bardzo prostą funkcjonalność (być może trochę mało przydatną w codziennym życiu),
polegającą na pobieraniu twitów z twittera, przechowywaniu ich, oraz wysyłaniu do [slacka](https://hackwawdisrupt.slack.com/). 

Na podstawie tego, aplikacja wystawia API do:

* Pobierania ostatnich, zaciągniętych twitów,
* Wyświetlania stanu aplikacji.

Ponadto, w wewnętrznych procesach, aplikacja:

* Cyklicznie pobiera twity z Twittera,
* Wysyła nowe twity na Slacka,
* Ponawia wysyłanie twitów na Slacka jeśli wystąpił błąd w komunikacji z serwerem,
* Nie wysyła duplikatów na Slacka jeśli serwer zwróci status 200,
* Zapisuje dane na dysku w katalogu `/storage`,
* Dane zapisane na dysku są dostępne dla wielu jej procesów,
* Jest w stanie wykryć problemy z zewnętrznymi usługami i określić ich przyczynę (np. długi czas odpowiedzi).

Aplikacja powinna implementować `/health` oraz `/latest`.

### /health

Pierwszy z endpointów służy do sprawdzenia stanu aplikacji:

```
{
    "app": "OPERATIONAL",
    "twitter": "OPERATIONAL",
    "slack": "OPERATIONAL"
}
```

Dostępne są następujące statusy dla `slacka` i `twittera`:

* `OPERATIONAL` - API Twittera i Slacka odpowiada poprawnie,
* `ERROR` - API Twittera i Slacka odpowiada kodem błędu 5xx,
* `SLOW` - API Twittera i Slacka odpowiada poprawnie, ale długo, powyżej 5s,
* `DOWN` - API Twittera i Slacka nie odpowiedziało w ciągu co najmniej 15s,

### /latest

Drugi z endpointów zwraca listę ostatnio pobranych twitów (może być wszystkich).

### Starter pack

Żeby nie tracić czasu na rzeczy mało interesujące, przygotowaliśmy szablony aplikacji,
realizującą (lub próbującą realizować) taką funkcjonalności, w różnych językach programowania:

* [Java](https://hackwaw.githost.io/hackwaw-disrupt/hackwaw-java-example)
* [PHP](https://hackwaw.githost.io/hackwaw-disrupt/hackwaw-php-example)
* [Scala](https://hackwaw.githost.io/hackwaw-disrupt/hackwaw-scala-example)
* [JS](https://hackwaw.githost.io/hackwaw-disrupt/hackwaw-js-example)
* [Python](https://hackwaw.githost.io/hackwaw-disrupt/hackwaw-python-example)

Na każdej z tych aplikacji, niezależnie od języka powinniśmy móc dostać się do API:

* Ostatnie twity `curl "http://localhost:8080/latest?page=0"`
* Status aplikacji `curl "http://localhost:8080/health"`
[Szczegółowe API Aplikacji - Swagger](/app-swagger.yaml)

### Docker

Żeby uwspólnić proces uruchamiania aplikacji, każda z naszych aplikacji jest uruchamiana na Dockerze.
W każdym z tych przykładów, znajdziecie `Dockerfile` odpowiedzialny za start aplikacji.
Nasze aplikacje, pobierają ze zmiennych środowiskowych URL do Twittera i Slacka.
Są to zmienne środowiskowe `SLACK_URL` oraz `TWITTER_URL`.

### GitLab CI

W każdym z tych repozytoriów znajdziecie plik `.gitlab-ci.yml`.
Jest on niezbędny do tego, żeby aplikacja była przetestowana i zdobywała punkty.

### Proxy

Nie będziemy korzystali z prawdziwego API Twittera i Slacka.
Przygotowaliśmy proxy tych serwisów, z bardzo prostym API.
Dzięki temu, możemy symulować niespodziewane działanie tych serwisów.
Ponadto, rozwijane przez Was aplikacje są prostsze.
Nie tracimy też czasu na autoryzację, do tych serwisów.

#### Twitter Proxy

Proxy twitterowe posiada jedną operację (GET) do pobierania listy twitów z danego przedziału czasowego. W odpowiedzi dostajemy listę twitów z tego okresu czasu.
[Szczegółowe API Proxy Twittera - Swagger](/twitter-proxy-swagger.yaml)

W celach podglądowych dostępne jest proxy na heroku:
* Przykład curl `curl "https://hackwaw-proxy-app.herokuapp.com/tweets?from=2016-03-01T18:54:00.706Z&to=2016-04-25T18:54:00.710Z"`
* Przykładowa odpowiedź:
```
[  
   {  
      "id":715208188246518688,
      "body":"Hello World",
      "date":"2016-03-30T16:04:54Z"
   }
]
```
* Format daty: RFC3339

#### Slack Proxy

Proxy twitterowe posiada jedną operację (POST) do wysyłania nowej wiadomości do Slacka.
[Szczegółowe API Proxy Slacka - Swagger](/slack-proxy-swagger.yaml)

W celach podglądowych dostępne jest proxy na heroku:
`curl -X POST --header "Content-Type: application/json" -d '{"date": "2016-03-21T16:12:12Z", "icon_url": "http://www.veryicon.com/icon/ico/System/Arcade%20Daze/Mario.ico", "team": "Best team ever", "text": "Testing #HacWaw", "tweetId": 1231231132123}' "https://hackwaw-proxy-app.herokuapp.com/push"`

###  Slack

Utworzyliśmy specjalnego Slacka, do którego publikuje nasza aplikacja. Możecie dołączyć do tego Slacka pod [tym adresem](https://hackwawdisrupt.slack.com/).

### Magiczne królestwo

Przygotowane przez nas przykłady działają dobrze w *Magicznym królestwie*.
Oznacza to, że nie ma problemów z siecią, z miejscem na dysku, z pamięcią,
a API zależnych serwisów działa bez zarzutów.
Ponadto, nikt niespodziewanie nie restartuje naszej aplikacji.
Jeżeli znaleźlibyśmy się w takim "Magicznym królestwie",
przygotowane przez nas przykłady działały by do końca świata, bez żadnych przeszkód.

## Start aplikacji

Nasz aplikacja powinna wstać na dockerze w mniej nić minutę. Jeżeli to się nie stanie, zostanie zabita i nie zdobędzie punktów. Żaden test się nie uruchomi.
Dodatkowo są przyznawane punkty za start w 3, 10 i 30 sekund.

## Nie żyjemy w magicznym królestwie

Niestety w naszych warunkach nasz aplikacja może natknąć się na następujęce problemy:
* API do Twittera/Slacka nie działa/działa wolno/działa źle
* Skończyło się miejsce na dysku/pamięć
* Dostajemy dużo requestów (DDOS)
* Aplikacja zostaje zrestartowana (czy zapmięta ostatnie twity)
* Nasz serwer ma tylko 50MB pamięci RAM (JVM'owcy mogą mieć z tym problem, najwyżej nie będzie z tego punktów)

## Baza danych

Żeby zdobyć punkty związane z testem, który restartuje aplikacje będzie wymagana jakaś forma bazy danych, którą trzeba sobie zrobić.

### Ograniczania

Prawa do zapisu mamy tylko w katalogach `/tmp` i `/storage`.
Nie ma możliwości zapisu poza tymi katalogami.

# Grywalizacja

Naszym celem będzie napisanie takiej aplikacji, która jak najlepiej poradzi sobie z tymi sytuacjami.
Na serwerze ciągłej integracji, będą chodziły testy które odpalają tę aplikację i ją oceniają.
Wygrywa ten zespół, który zdobędzie jak najwięcej punktów.

## Sędzia

Na GitLab, aplikacja będzie przechodziła serię testów, za które będzie otrzymywała punkty.
Pełen podgląd testów, możemy zobaczyć w każdym w widoku swojego projektu oraz na naszym [dashbordzie](http://hackwaw-stats.touk.pl/specs?id=962704)

## Testy podstawowe

W każdej przykładowej aplikacji, powinny przechodzić testy:

* Zaczynające się od: `Verify normal` - sprawdzające podstawową funkcjonalność aplikacji
* Sprawdzające w ile startuje aplikacja `Check startup time IT should start in X seconds`.
Które weryfikują w ile czasu wstaje aplikacja.

## Testowanie

Zalecane jest, żeby odpalać testy lokalnie, przed wypchnięciem. Jak to zrobić? Należy:

* mieć Dockera, najlepiej na maszynie Linuxowej / OSX-owej
* pobrać projekt [go-judge](https://hackwaw.githost.io/hackwaw-disrupt/go-judge)
* uruchomić skrypt bashowy `run-locally` ze wskazaniem na projekt np `bash ./run-locally ../apps/hackwaw-java-example/`

Inna możliwość to wykonać polecenia z katalogu twojej aplikacji:

```
docker run -it --rm --privileged -v "/var/run/docker.sock:/var/run/docker.sock" \
    -v "$(pwd):/app:ro" quay.io/hackwawdisrupt/go-judge
```

Lub przekazując ścieżkę bezpośrednio do kontenera:

```
docker run -it --rm --privileged -v "/var/run/docker.sock:/var/run/docker.sock" \
    -v "/home/path/to/my/application:/app:ro" quay.io/hackwawdisrupt/go-judge
```

Docker pobierze obraz z aplikacją sędziego, zbuduje i uruchomi aplikację z katalogu `/home/path/to/my/application`.

W trakcie Hackatonu możemy was prosić o zaktualizowanie obrazu:

```
docker pull quay.io/hackwawdisrupt/go-judge
```

## Aplikacja referencyjna

Udostępniamy aplikację referencyjną [go-reference](https://gitlab.com/hackwaw-disrupt/go-reference).
Jest to aplikacja, która przechodzi zdecydowaną większość testów, ale nie wszystkie.
Pomoże wam określić kierunek w którym należy poprawić aplikacje w waszych językach.

## Ile mam punktów?

Sprawdź w dashboardzie:
* [Dashboard - statystyki](http://hackwaw-stats.touk.pl/)
* [Dashboard - wykresy](http://hackwaw-grafana.touk.pl/) - login hackwaw/hackwaw

## Tips and tricks

* Zadbaj o to, żeby twoja aplikacja wstawała szybko - zyskasz punkty za szybkie wstanie
* Zadbaj o to, żeby twój obraz szybko się budował.
Dockerfile powinien używać warstw pośrednich do zależności (np. maven, npm itp),
a na samym końcu kod palikacji.
[Przykładowy Dockerfile, korzystający z warstw pośrednich dla Javy.](https://hackwaw.githost.io/hackwaw-disrupt/hackwaw-java-example/blob/master/Dockerfile)
Dzięki temu będzie szybszy feedback, przez co będzie można uruchomić więcej testów. 

### Uruchamiaj testy lokalnie
* Sklonuj projekt [Judge](git@hackwaw.githost.io:hackwaw-disrupt/go-judge.git)
* Uruchom skrypt `run-locally`, wskazując adres projektu z kodem swojej aplikacji np. `./run-locally ../apps/hackwaw-java-example/`

## Jak wziąć udział w zabawie?

* Zrób forka projektu,
* Zmień coś w tej aplikacji,
* Wypchnij zmiany do repozytorium,
* Zmiany będą testowane, a wynik będzie wyświetlany na Dashboardzie.

## Zbiór linków

* [Slack aplikacji](https://hackwawdisrupt.slack.com/)
* [Katalog główny repozytorium](https://hackwaw.githost.io/groups/hackwaw-disrupt)
* [Edytor/Podgląd Swagger API](http://editor.swagger.io/#/)
