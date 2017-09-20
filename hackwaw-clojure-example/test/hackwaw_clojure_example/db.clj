(ns hackwaw-clojure-example.db
  (:require [duct.component.ragtime :as ragtime]))

(defn migrate [system]
  (-> system :ragtime ragtime/reload ragtime/migrate))

(defn rollback
  ([system]  (rollback system 1))
  ([system x] (-> system :ragtime ragtime/reload (ragtime/rollback x))))

(defn get-connection [system]
  (get-in system [:db :spec]))