(ns hackwaw-clojure-example.system
  (:require [com.stuartsierra.component :as component]
            [duct.component.endpoint :refer [endpoint-component]]
            [duct.component.handler :refer [handler-component]]
            [duct.component.hikaricp :refer [hikaricp]]
            [duct.component.ragtime :refer [ragtime]]
            [duct.middleware.not-found :refer [wrap-not-found]]
            [meta-merge.core :refer [meta-merge]]
            [ring.component.jetty :refer [jetty-server]]
            [ring.middleware.defaults :refer [wrap-defaults api-defaults]]
            [hackwaw-clojure-example.endpoint.api :refer [api-endpoint]]
            [hackwaw-clojure-example.repository.tweets :as tr]
            [hackwaw-clojure-example.infrastructure.twitter :as ts]
            [hackwaw-clojure-example.infrastructure.slack :as ss]
            [hackwaw-clojure-example.infrastructure.pipeline :as p]))

(def base-config
  {:app {:middleware [[wrap-not-found :not-found]
                      [wrap-defaults :defaults]]
         :not-found  "Resource Not Found"
         :defaults   (meta-merge api-defaults {})}
   :ragtime {:resource-path "hackwaw_clojure_example/migrations"}})

(defn new-system [config]
  (let [config (meta-merge base-config config)]
    (-> (component/system-map
          :app (handler-component (:app config))
          :http (jetty-server (:http config))
          :db (hikaricp (:db config))
          :ragtime (ragtime (:ragtime config))
          :api (endpoint-component api-endpoint)
          :repo (tr/new-jdbc-tweet-repository)
          :pipeline (p/new-tweet-pipeline)
          :tweet-source (ts/new-http-tweet-source (:twitter config))
          :slack-sink (ss/new-http-slack-sink (:slack config)))
        (component/system-using
          {:http    [:app]
           :app     [:api]
           :ragtime [:db]
           :api     [:repo]
           :repo    [:db]
           :pipeline [:repo :tweet-source :slack-sink]}))))
