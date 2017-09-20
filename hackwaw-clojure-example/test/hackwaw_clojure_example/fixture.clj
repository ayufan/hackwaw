(ns hackwaw-clojure-example.fixture
  (:require [com.stuartsierra.component :as component]
            [clojure.test :refer [is]]
            [hackwaw-clojure-example.db :refer [migrate rollback]]))

(defn system-migrate-fixture [system-var f]
  (alter-var-root system-var component/start)
  (let [system (var-get system-var)]
    (is (not (nil? (:db system))))
    (is (not (nil? (:ragtime system))))
    (migrate system)
    (f)
    (rollback system))
  (alter-var-root system-var component/stop))
