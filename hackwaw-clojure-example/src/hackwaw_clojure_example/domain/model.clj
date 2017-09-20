(ns hackwaw-clojure-example.domain.model
  (:require [schema.core :as s])
  (:import (org.joda.time DateTime)))

(s/defschema Tweet
  {:id         s/Int
   :twitter_id s/Int
   :link       s/Str
   :body       s/Str
   :date       DateTime})

(def Status (s/enum "OPERATIONAL" "DOWN" "SLOW" "ERROR" "UNNECESSARY"))

(s/defschema Health
  {:app      Status
   :database Status
   :slack    Status
   :twitter  Status})
