(ns hackwaw-clojure-example.repository.repo-test
  (:require [clojure.test :refer :all]
            [clj-time.core :as t]
            [com.stuartsierra.component :as component]
            [duct.component.hikaricp :refer [hikaricp]]
            [duct.component.ragtime :refer [ragtime]]
            [hackwaw-clojure-example.system :as s]
            [hackwaw-clojure-example.fixture :as f]
            [hackwaw-clojure-example.repository.tweets :as tr]))

(def config
  (merge s/base-config
         {:db {:uri "jdbc:h2:mem:hackwaw;MODE=PostgreSQL;DB_CLOSE_ON_EXIT=FALSE"}}))

(def system
  (-> (component/system-map
        :db (hikaricp (:db config))
        :ragtime (ragtime (:ragtime config))
        :repo (tr/new-jdbc-tweet-repository))
      (component/system-using
        {:ragtime [:db]
         :repo    [:db]})))

(use-fixtures :each (partial f/system-migrate-fixture #'system))

(deftest repository-test
  (require '[hackwaw-clojure-example.infrastructure.sql :reload :all])
  (let [repo (:repo system)]
    (testing "should store tweet"
      (let [tweet {:date "2016-03-30T15:00:57Z"
                   :id   715194615490486300
                   :body ".@netflix szuka Instagramerów - wcześniejszy"}]
        (tr/insert-tweet repo tweet)
        (let [all-tweets (tr/find-all repo)]
          (is (= 1 (count all-tweets)))
          (is (= [715194615490486300] (map :twitter_id all-tweets))))
        (is (t/before? (tr/get-last-tweet-date repo) (t/now)))))))
