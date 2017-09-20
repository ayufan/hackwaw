(ns hackwaw-clojure-example.repository.tweets
  (:require [schema.core :as s]
            [hackwaw-clojure-example.infrastructure.sql :as sql]
            [hackwaw-clojure-example.infrastructure.twitter :as ts]
            [hackwaw-clojure-example.domain.model :as m]))

(defprotocol TweetRepository
  (insert-tweet [this tweet])
  (get-last-tweet-date [this])
  (find-all [this]))

(s/defn ->db-tweet
  [tweet :- ts/Tweet] :- m/Tweet
  {:body (:body tweet)
   :twitter_id (:id tweet)
   :link (str "twitter.com/anyuser/status/" (:id tweet))
   :date (:date tweet)})

(defn- get-connection [db]
  (:spec db))

(defrecord JdbcTweetRepository [db]
  TweetRepository
  (insert-tweet [this tweet]
    (sql/save-tweet! (get-connection db) (->db-tweet tweet)))
  (get-last-tweet-date [this]
    (-> (sql/get-last-tweet-date (get-connection db))
        :max))
  (find-all [this]
    (sql/get-all-tweets (get-connection db))))

(defn new-jdbc-tweet-repository []
  (->JdbcTweetRepository {}))