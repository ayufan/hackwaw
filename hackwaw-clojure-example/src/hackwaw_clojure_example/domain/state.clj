(ns hackwaw-clojure-example.domain.state
  (:require [schema.core :as s]
            [hackwaw-clojure-example.domain.model :as m]))

(def health (atom {:app "OPERATIONAL"
                   :database "UNNECESSARY"
                   :slack "SLOW"
                   :twitter "OPERATIONAL"}))

(s/defn mutate!
  [entry :- s/Keyword
   value :- m/Status]
  (swap! health assoc entry value))
