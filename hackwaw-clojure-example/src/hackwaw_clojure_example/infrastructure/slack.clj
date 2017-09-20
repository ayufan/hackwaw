(ns hackwaw-clojure-example.infrastructure.slack
  (:require [clj-http.client :as http]
            [schema.core :as s]
            [taoensso.timbre :as l]
            [hackwaw-clojure-example.infrastructure.twitter :as t]))

(s/defschema SlackMessage
  {:tweetId  s/Int
   :date     s/Str
   :icon_url s/Str
   :text     s/Str
   :team     s/Str})

(s/defn ->slack-message
  [tweet :- t/Tweet] :- SlackMessage
  {:tweetId  (:id tweet)
   :date     (:date tweet)
   :icon_url "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-03-01/23711067665_1081343b8ffaa157a175_132.png"
   :text     (:body tweet)
   :team     "Team A"})

(defprotocol SlackSink
  (push [this tweet]))

(defrecord HttpSlackSink [config]
  SlackSink
  (push [this tweet]
    (let [slack-message (->slack-message tweet)]
      (l/info "Pushing slack message" slack-message)
      (http/post (str (:uri config) "/push")
                 {:form-params  slack-message
                  :content-type :json}))))

(defn new-http-slack-sink [config]
  (->HttpSlackSink config))
