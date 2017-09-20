(ns hackwaw-clojure-example.infrastructure.sql
  (:require [hugsql.core :as hugsql]))

(hugsql/def-db-fns "hackwaw_clojure_example/sql/tweets.sql")
