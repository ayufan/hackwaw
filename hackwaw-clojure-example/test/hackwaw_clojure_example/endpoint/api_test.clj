(ns hackwaw-clojure-example.endpoint.api-test
  (:require [com.stuartsierra.component :as component]
            [duct.component.endpoint :refer [endpoint-component]]
            [duct.component.handler :refer [handler-component]]
            [duct.component.hikaricp :refer [hikaricp]]
            [duct.component.ragtime :refer [ragtime]]
            [ring.component.jetty :refer [jetty-server]]
            [clojure.test :refer :all]
            [peridot.core :refer :all]
            [hackwaw-clojure-example.system :as s]
            [hackwaw-clojure-example.repository.tweets :as tr]
            [hackwaw-clojure-example.db :refer :all]
            [hackwaw-clojure-example.endpoint.api :as api]
            [hackwaw-clojure-example.fixture :as f]
            [clojure.java.io :as io]
            [cheshire.core :as json]))

(def config
  (merge s/base-config
         {:db {:uri "jdbc:h2:mem:hackwaw;MODE=PostgreSQL;DB_CLOSE_ON_EXIT=FALSE"}}))

(def system
  (-> (component/system-map
        :app (handler-component (:app config))
        :db (hikaricp (:db config))
        :ragtime (ragtime (:ragtime config))
        :api (endpoint-component api/api-endpoint)
        :repo (tr/new-jdbc-tweet-repository))
      (component/system-using
        {:app     [:api]
         :ragtime [:db]
         :api     [:repo]
         :repo    [:db]})))

(use-fixtures :each (partial f/system-migrate-fixture #'system))

(defn- status [session]
  (get-in session [:response :status]))

(defn- to-json [response]
  (-> (get-in response [:response :body])
      (io/reader)
      (json/decode-stream true)))

(defn get-latest [handler]
  (-> (session handler)
      (request "/latest")))

(deftest api-test
  (require '[hackwaw-clojure-example.infrastructure.sql :reload :all])
  (let [handler (get-in system [:app :handler])
        repo (get system :repo)]

    (testing "api responds"
      (let [resp (get-latest handler)]
        (is (= 200 (status resp)))))

    (testing "should return added tweet"
      (let [tweet {:date "2016-03-30T15:00:57Z"
                   :id   715194594514698200
                   :body ".@netflix szuka Instagramerów - wcześniejszy"}]
        (tr/insert-tweet repo tweet)
        (let [resp (get-latest handler)
              tweets (to-json resp)]
          (is (= 200 (status resp)))
          (is (= 1 (count tweets)))
          (is (= [715194594514698200] (map :twitter_id tweets))))))))
