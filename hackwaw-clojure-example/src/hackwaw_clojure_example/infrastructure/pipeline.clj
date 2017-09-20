(ns hackwaw-clojure-example.infrastructure.pipeline
  (:require [com.stuartsierra.component :as component]
            [taoensso.timbre :as l]
            [chime :refer [chime-ch]]
            [clj-time.core :as t]
            [clj-time.periodic :refer [periodic-seq]]
            [clojure.core.async :as a :refer [<! go-loop]]
            [hackwaw-clojure-example.repository.tweets :as tr]
            [hackwaw-clojure-example.infrastructure.twitter :as ts]
            [hackwaw-clojure-example.infrastructure.slack :as ss]))

(defn process [repo tweet-source slack]
  (l/info "Processing tweets")
  (try
    (let [tweets (->> (or (tr/get-last-tweet-date repo) (t/minus (t/now) (t/seconds 300)))
                      (ts/get-tweets-since tweet-source))]
      (doseq [tweet tweets]
        (l/info "Got tweet" tweet)
        (tr/insert-tweet repo tweet)
        (ss/push slack tweet)))
    (catch Exception ex
      (l/error "Exception in pipeline" ex))))

(defrecord TweetPipeline [repo tweet-source slack-sink]
  component/Lifecycle
  (start [this]
    (l/info "Starting Tweet pipeline")
    (let [ch (chime-ch (periodic-seq (t/now) (t/seconds 5)))]
      (go-loop []
        (when (<! ch)
          (process repo tweet-source slack-sink)
          (recur)))
      (assoc this :chime-ch ch)))
  (stop [this]
    (l/info "Stopping Tweet pipeline")
    (if-let [ch (:chime-ch this)]
      (do (a/close! ch)
          (dissoc this :chime-ch))
      this)))

(defn new-tweet-pipeline []
  (map->TweetPipeline {}))
