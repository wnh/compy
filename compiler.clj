
(ns compiler
  (:require [instaparse.core :as insta]
            [clojure.test :as test :refer [deftest testing is]]))

(def parser
  (insta/parser
   "<expr> = add | sub | num
    add  = num <'+'> expr
    sub  = num <'-'> expr
    num = ws #'[0-9]+' ws
    <ws> = <#'\\s*'>"))


(defn bail [token]
  (throw (Exception. "Unexpeccted: " (pr-str token))))

(defn emit [& args]
  (apply printf args)
  (println ""))

(defn node-is? [t node]
  (= (first node) t)) 

(defn emit-num [node]
  ;; [:num "123"]
  (emit "  push $%d" (Integer/parseInt (second node))))

(declare emit-expr)

(defn emit-binop [node]
  (let [[nodetype left right] node]
    (emit-expr left)
    (emit-expr right)
    (emit "  pop %%rbx")
    (emit "  pop %%rax")
    (case nodetype
      :sub (emit "  sub %%rbx, %%rax")
      :add (emit "  add %%rbx, %%rax"))
    (emit "  push %%rax")))

(defn emit-expr [node]
  (case (first node); type
    :num (emit-num node)
    :add (emit-binop node)
    :sub (emit-binop node)))
                 

(defn emit-program [ast]
  (emit "  .globl main")
  (emit "main:")
  (emit-expr ast)
  (emit "  pop %%rax")
  (emit "  ret"))

(comment
  (compiler "5+20-4")
  (compile-and-run "5+20-4")
  ;;
  )

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
  (let [ast (first (parser src))]
    (emit-program ast)))

;(compiler "0")

(defn -main []
  (let [src (first *command-line-args*)]
    (compiler src)))


(defn compile-and-run [src] 
  (let [asm (with-out-str (compiler src))]
    (spit "tmp.s" asm)
    (prn (clojure.java.shell/sh "gcc" "-static" "-o" "tmp" "tmp.s"))
    (:exit (clojure.java.shell/sh "./tmp"))))


(deftest  test-compiler-works
  (testing "Constants"
    (is (=  0 (compile-and-run "0")))
    (is (= 42 (compile-and-run "42")))
    (is (= 21 (compile-and-run "5+20-4")))
    (is (= 41 (compile-and-run  " 12 + 34 - 5 ")))
    ;(is (= 47 (compile-and-run  "5+6*7")))
;assert 47 '5+6*7'
))

;(test-compiler-works)

;assert 42 42
;assert 21 '5+20-4'
;assert 41 ' 12 + 34 - 5 '
;assert 47 '5+6*7'
;assert 15 '5*(9-6)'
;assert 4 '(3+5)/2'
