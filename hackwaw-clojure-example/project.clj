(defproject hackwaw-clojure-example "0.1.0-SNAPSHOT"
  :description "FIXME: write description"
  :url "http://example.com/FIXME"
  :min-lein-version "2.0.0"
  :dependencies [[org.clojure/clojure "1.8.0"]
                 [com.stuartsierra/component "0.3.1"]
                 [compojure "1.5.0"]
                 [metosin/compojure-api "1.0.2"]
                 [duct "0.5.10"]
                 [environ "1.0.2"]
                 [meta-merge "0.1.1"]
                 [ring "1.4.0"]
                 [ring/ring-defaults "0.2.0"]
                 [ring-jetty-component "0.3.1"]
                 [duct/hikaricp-component "0.1.0"]
                 [com.h2database/h2 "1.4.191"]
                 [com.layerware/hugsql "0.4.6"]
                 [clj-time "0.11.0"]
                 [duct/ragtime-component "0.1.3"]
                 [jarohen/chime "0.1.9"]
                 [clj-http "2.1.0"]
                 [com.taoensso/timbre "4.3.1"]]
  :plugins [[lein-environ "1.0.2"]
            [lein-gen "0.2.2"]]
  :generators [[duct/generators "0.5.10"]]
  :duct {:ns-prefix hackwaw-clojure-example}
  :main ^:skip-aot hackwaw-clojure-example.main
  :target-path "target/%s/"
  :aliases {"gen"   ["generate"]
            "setup" ["do" ["generate" "locals"]]}
  :profiles
  {:dev  [:project/dev  :profiles/dev]
   :test [:project/test :profiles/test]
   :uberjar {:aot :all
             :env {:port "8080"
                   :database-url "jdbc:h2:/tmp/hackwaw;MODE=PostgreSQL;DB_CLOSE_ON_EXIT=FALSE"}}
   :profiles/dev  {}
   :profiles/test {}
   :project/dev   {:dependencies [[reloaded.repl "0.2.1"]
                                  [org.clojure/tools.namespace "0.2.11"]
                                  [org.clojure/tools.nrepl "0.2.12"]
                                  [eftest "0.1.1"]
                                  [kerodon "0.7.0"]]
                   :source-paths ["dev"]
                   :repl-options {:init-ns user}
                   :env {:port "3000"
                         :database-url "jdbc:h2:./hackwaw;MODE=PostgreSQL;DB_CLOSE_ON_EXIT=FALSE"
                         :twitter-url "https://hackwaw-twitter-proxy.herokuapp.com"
                         :slack-url "https://hackwaw-slack-proxy.herokuapp.com"}}
   :project/test  {}})
