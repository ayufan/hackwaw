(ns hackwaw-clojure-example.config
  (:require [environ.core :refer [env]]))

(def defaults
  ^:displace {:http {:port 3000}})

(def environ
  {:http {:port (some-> env :port Integer.)}
   :db   {:uri  (some-> env :database-url String.)}
   :twitter {:uri (some-> env :twitter-url String.)}
   :slack {:uri (some-> env :slack-url String.)}})
