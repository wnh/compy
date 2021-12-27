
(ns compiler
  (:require [instaparse.core :as insta]
            [clojure.test :as test :refer [deftest testing is]]))

(def parser
  (insta/parser
   "<expr> = sum-expr

    <sum-expr> = add | sub | fact-expr
    add  = fact-expr <'+'> sum-expr
    sub  = fact-expr <'-'> sum-expr

    <fact-expr> = mul | num
    mul = num <'*'> fact-expr

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
      :add (emit "  add %%rbx, %%rax")
      :mul (emit "  mul %%rbx"))
    (emit "  push %%rax")))

(defn emit-expr [node]
  (case (first node); type
    :num (emit-num node)
    :add (emit-binop node)
    :sub (emit-binop node)
    :mul (emit-binop node)))
                 

(defn emit-program [ast]
  (emit "  .globl main")
  (emit "main:")
  (emit-expr ast)
  (emit "  pop %%rax")
  (emit "  ret"))

(defn compiler [src]
  (let [ast (first (parser src))]
    (emit-program ast)))

(defn -main []
  (let [src (first *command-line-args*)]
    (compiler src)))


(defn compile-and-run [src] 
  (let [asm (with-out-str (compiler src))]
    (spit "tmp.s" asm)
    (let [compile-out (clojure.java.shell/sh "gcc" "-static" "-o" "tmp" "tmp.s")]
      (if (< 0 (:exit compile-out))
        (throw (Exception. (:err compile-out)))
        (:exit (clojure.java.shell/sh "./tmp"))))))


(deftest  test-compiler-works
  (testing "Constants"
    (is (=  0 (compile-and-run "0")))
    (is (= 42 (compile-and-run "42"))))
  (testing "addition and subtraction"
    (is (= 21 (compile-and-run "5+20-4")))
    (is (= 41 (compile-and-run  " 12 + 34 - 5 "))))
  (testing "multiplication"
    (is (=  3 (compile-and-run  " 5*3 -12")))
    (is (= 47 (compile-and-run  "5+6*7")))))

;(is (= 15 (compile-and-run "5*(9-6)")))
;(is (=  4 (compile-and-run "(3+5)/2")))
