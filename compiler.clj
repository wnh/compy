
(ns compiler
  (:require [instaparse.core :as insta]
            [clojure.test :as test :refer [deftest testing is]]))

(def parser
  (insta/parser
   "
    stmts = (expr <';'>)+
    <expr> = eq-expr

    <eq-expr> = eq | neq | rel-expr
    eq  = rel-expr <'=='> rel-expr
    neq = rel-expr <'!='> rel-expr

    <rel-expr> = gt | gte | lt | lte | sum-expr
    gt  = sum-expr ws <'>'> ws rel-expr
    gte = sum-expr ws <'>='> ws rel-expr
    lt  = sum-expr ws <'<'> ws rel-expr
    lte = sum-expr ws <'<='> ws rel-expr

    <sum-expr> = add | sub | term-expr
    add  = term-expr <'+'> sum-expr
    sub  = term-expr <'-'> sum-expr

    <term-expr> = mul | div | factor-expr
    mul = factor-expr <'*'> term-expr
    div = factor-expr <'/'> term-expr

    <factor-expr> = num | uneg | uplus
                  | <'('> expr <')'>
    uneg  = ws <'-'> ws factor-expr
    uplus = ws <'+'> ws factor-expr
    num = ws #'[0-9]+' ws
    <ws> = <#'\\s*'>"))

(take 2 (insta/parses parser "- - +10"))

(defn bail [token]
  (throw (Exception. "Unexpeccted: " (pr-str token))))

(defn emit [& args]
  (apply printf args)
  (println ""))

(declare emit-expr)

(defn emit-num [node]
  ;; [:num "123"]
  (emit "  push $%d" (Integer/parseInt (second node))))

(defn emit-unary [node]
  ;; [:uneg  "123"]
  ;; [:uplus "123"]
  (case (first node)
    :uneg (do (emit-expr (second node))
              (emit "  pop %%rax")
              (emit "  neg %%rax")
              (emit "  push %%rax"))
    :uplus (do (emit-expr (second node)))))

(defn emit-rel [node]
  (emit-expr (node 1))
  (emit-expr (node 2))
  (emit "  pop %%rax")
  (emit "  pop %%rbx")
  (case (node 0)
    :eq (do (emit "  cmp %%rbx, %%rax")
            (emit "  sete %%al"))
    :neq (do (emit "  cmp %%rbx, %%rax")
             (emit "  setne %%al"))
    :gt (do (emit "  cmp %%rbx, %%rax")
            (emit "  setl %%al"))
    :gte (do (emit "  cmp %%rbx, %%rax")
             (emit "  setle %%al"))
    :lt (do (emit "  cmp %%rax, %%rbx")
            (emit "  setl %%al"))
    :lte (do (emit "  cmp %%rax, %%rbx")
             (emit "  setle %%al")))
  (emit "  push %%rax"))

(defn emit-binop [node]
  (let [[nodetype left right] node]
    (emit-expr left)
    (emit-expr right)
    (emit "  pop %%rbx")
    (emit "  pop %%rax")
    (case nodetype
      :sub (emit "  sub %%rbx, %%rax")
      :add (emit "  add %%rbx, %%rax")
      :mul (emit "  imul %%rbx, %%rax")
      :div (do (emit "  xor %%rdx, %%rdx")
               (emit "  div %%rbx")))
    (emit "  push %%rax")))

(defn emit-expr [node]
  (case (first node); type
    :num (emit-num node)

    :add (emit-binop node)
    :sub (emit-binop node)
    :mul (emit-binop node)
    :div (emit-binop node)

    :uneg  (emit-unary node)
    :uplus (emit-unary node)

    :eq   (emit-rel node)
    :neq  (emit-rel node)
    :gt   (emit-rel node)
    :gte  (emit-rel node)
    :lt   (emit-rel node)
    :lte  (emit-rel node)))

(defn emit-stmts [node]
  ;; [:stmts expr ...]
  (doseq [expr (rest node)]
    (emit-expr expr)
    (emit "  pop %%rax")))


(defn emit-program [ast]
  (emit "  .globl main")
  (emit "main:")
  (emit-stmts ast)
  (emit "  ret"))

(defn compiler [src]
  (let [ast (parser src)]
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
        (let [out (clojure.java.shell/sh "./tmp")]
          (prn out)
          (:exit out))))))


(deftest  test-compiler-works
  (testing "Constants"
    (is (=  0 (compile-and-run "0;")))
    (is (= 42 (compile-and-run "42;"))))
  (testing "addition and subtraction"
    (is (= 21 (compile-and-run "5+20-4;")))
    (is (= 41 (compile-and-run  " 12 + 34 - 5 ;"))))
  (testing "multiplication"
    (is (=  3 (compile-and-run  " 5*3 -12;")))
    (is (= 47 (compile-and-run  "5+6*7;"))))
  (testing "brackets"
    (is (= 15 (compile-and-run "5*(9-6);")))
    (is (= 15 (compile-and-run "5*(6*2-9);")))
    (is (= 51 (compile-and-run "5* 6*2-9 ;"))))
  (testing "division"
    (is (=  4 (compile-and-run "(3+5)/2;")))
    (is (=  4 (compile-and-run "(3+10)/3;"))))
  (testing "unary minus"
    (is (= 10 (compile-and-run "-10+20;")))
    (is (= 10 (compile-and-run "- -10;")))
    (is (= 10 (compile-and-run "- - +10;"))))
  (testing "equality"
    (is (= 0 (compile-and-run "0 ==1;")))
    (is (= 1 (compile-and-run "42==42;")))
    (is (= 1 (compile-and-run "0!=1;")))
    (is (= 0 (compile-and-run "42!=42;"))))
  (testing "greater than"
    (is (= 1 (compile-and-run "10>5;")))
    (is (= 0 (compile-and-run "10>50;")))
    (is (= 0 (compile-and-run "10>=50;")))
    (is (= 1 (compile-and-run "10>=10;"))))
  (testing "less than"
    (is (= 1 (compile-and-run "0<1;")))
    (is (= 0 (compile-and-run "1<1;")))
    (is (= 0 (compile-and-run "2<1;")))
    (is (= 1 (compile-and-run "0<=1;")))
    (is (= 1 (compile-and-run "1<=1;")))
    (is (= 0 (compile-and-run "2<=1;"))))
  (testing "multiple expressions"
    (is (= 1 (compile-and-run "3;2;1;")))))
