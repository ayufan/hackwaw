(ns hackwaw-clojure-example.endpoint.api
  (:require [compojure.api.sweet :refer :all]
            [clojure.set :refer [rename-keys]]
            [clj-time.jdbc]
            [ring.util.http-response :refer :all]
            [hackwaw-clojure-example.domain.model :as m]
            [hackwaw-clojure-example.domain.state :as s]
            [hackwaw-clojure-example.repository.tweets :as tr]))

(defn api-endpoint [{repo :repo}]
  (api
    {:swagger
     {:ui "/api-docs"
      :spec "/swagger.json"
      :data {:info {:title "Tweet API"
                    :description "Clojure Tweet App"}}}}
    (GET "/latest" []
         :return [m/Tweet]
         (ok
           (tr/find-all repo)))

    (GET "/health" []
         :return m/Health
         (ok @s/health))))
