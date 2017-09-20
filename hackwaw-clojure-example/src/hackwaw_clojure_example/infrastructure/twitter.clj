(ns hackwaw-clojure-example.infrastructure.twitter
  (:require [schema.core :as s]
            [clj-http.client :as http]
            [clj-time.core :as t]
            [taoensso.timbre :as l]))

(s/defschema Tweet
  {:id   s/Int
   :body s/Str
   :date s/Str})

(defprotocol TweetSource
  (get-tweets-since [this from-date]))

(defrecord HttpTweetSource [config]
  TweetSource
  (get-tweets-since [this from-date]
    (l/info "Getting tweets since" from-date)
    (-> (http/get (str (:uri config) "/tweets")
                  {:query-params {"from" (str from-date)
                                  "to"   (str (t/now))}
                   :as           :json})
        :body)))

(defn new-http-tweet-source [config]
  (->HttpTweetSource config))