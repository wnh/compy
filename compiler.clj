
(ns compiler
  (:require [instaparse.core :as insta]))

(def parser
  (insta/parser
   "num = ws #'[0-9]+' ws
    <ws> = <#'\\s*'>"))

(parser "0")

(defn bail [token]
  (throw (Exception. "Unexpeccted: " (pr-str token))))

(defn emit [& args]
  (apply printf args)
  (println ""))

(defn node-is? [t node]
  (= (first node) t)) 

(defn emit-num [node]
  ;; [:num "123"]
  (emit "  mov $%d, %%rax" (Integer/parseInt (second node))))


(defn emit-program [ast]
  (emit "  .globl main")
  (emit "main:")
  (emit-num ast)
  (emit "  ret"))

;;;(defn emit-code [ast]
;;;  (emit "  .globl main")
;;;  (emit "main:")
;;;  (emit "  mov $" (first ast) ", %rax")
;;;  (let loop ((ins (rest ast)))
;;;    (cond 
;;;     ((null? ins) #f)
;;;     ((equal? (first ins) '+)
;;;      (begin (emit "  add $" (second ins) ", %rax\n")
;;;	     (loop (cddr ins))))
;;;      ((equal? (first ins) '-)
;;;       (begin (emit "  sub $" (second ins) ", %rax\n")
;;;	      (loop (cddr ins))))))
;;;  (emit "ret\n"))

(defn compiler [src]
  (let [ast (parser src)]
    (emit-program ast)))

;(compiler "0")

(defn main []
  (let [src (first *command-line-args*)]
    (compiler src)))

(main)


