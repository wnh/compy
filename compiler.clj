
(ns compiler
  (:require [instaparse.core :as insta]
            [clojure.test :as test :refer [deftest testing is]]))

(def parser
  (insta/parser
   "
    block = ws <'{'> ws (stmt ws)+ ws <'}'> ws
    <stmt> = stmt-single ws <';'>
           | block
    <stmt-single> = return | assign | expr
    assign = ident <'='> expr
    return = ws <'return'> ws expr
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

    <factor-expr> = num | uneg | uplus | var-ref
                  | <'('> expr <')'>
    var-ref = ident
    uneg    = ws <'-'> ws factor-expr
    uplus   = ws <'+'> ws factor-expr
    num     = ws #'[0-9]+' ws
    ident   = ws !keyword #'[a-zA-Z][a-zA-Z0-9]*' ws
    <keyword> = <'return'>
    <ws> = <#'\\s*'>"))


(def ^:dynamic *debugging* false)
(defn debug [& args]
  (when *debugging*
    (binding [*out* *err*]
      (apply prn args))))

(comment
  (compile-and-run "{x=2;y=3; return 9;}")
  (parser "{x=2;y=3; return 0;}")
  (parser " { return 12 + 34 - 5 ; }")
  (take 2 (insta/parses parser "{ {1; {2;} return 3;} }"))
  (parser "{ {1; {2;} return 3;} }")
  ;; assert 3 '{ {1; {2;} return 3;} }'
  )

(meta
 (get-in (parser "5+3*12;") [1 2 2]))

(defn bail [token]
  (throw (Exception. "Unexpeccted: " (pr-str token))))

(defn emit [env & args]
  (apply printf args)
  (println ""))

(declare emit-expr)

(defn emit-num [env node]
  (debug :emit-num)
  ;; [:num "123"]
  (emit env "  push $%d" (Integer/parseInt (second node))))

(defn emit-unary [env node]
  (debug :emit-unary)
  ;; [:uneg  "123"]
  ;; [:uplus "123"]
  (case (first node)
    :uneg (do (emit-expr env (second node))
              (emit env "  pop %%rax")
              (emit env "  neg %%rax")
              (emit env "  push %%rax"))
    :uplus (do (emit-expr env (second node)))))

(defn emit-rel [env node]
  (debug :emit-rel)
  (emit-expr env (node 1))
  (emit-expr env (node 2))
  (emit env "  pop %%rax")
  (emit env "  pop %%rbx")
  (case (node 0)
    :eq (do (emit env "  cmp %%rbx, %%rax")
            (emit env "  sete %%al"))
    :neq (do (emit env "  cmp %%rbx, %%rax")
             (emit env "  setne %%al"))
    :gt (do (emit env "  cmp %%rbx, %%rax")
            (emit env "  setl %%al"))
    :gte (do (emit env "  cmp %%rbx, %%rax")
             (emit env "  setle %%al"))
    :lt (do (emit env "  cmp %%rax, %%rbx")
            (emit env "  setl %%al"))
    :lte (do (emit env "  cmp %%rax, %%rbx")
             (emit env "  setle %%al")))
  (emit env "  push %%rax"))

(defn emit-binop [env node]
  (let [[nodetype left right] node]
    (emit-expr env left)
    (emit-expr env right)
    (emit env "  pop %%rbx")
    (emit env "  pop %%rax")
    (case nodetype
      :sub (emit env "  sub %%rbx, %%rax")
      :add (emit env "  add %%rbx, %%rax")
      :mul (emit env "  imul %%rbx, %%rax")
      :div (do (emit env "  xor %%rdx, %%rdx")
               (emit env "  div %%rbx")))
    (emit env "  push %%rax")))

(defn emit-assignment [env node]
  (debug :emit-assignment env)
  ;;[:assign [:ident "a"] [:num "2"]]
  (let [ident (second node)
        offsets (:local-offsets env)
        offset (get offsets ident)]
    (emit-expr env (nth node 2))
    (emit env "  pop %%rax")
    (emit env "  lea %d(%%rbp), %%rdi" (- offset))
    (emit env "  mov %%rax, (%%rdi)")
    (emit env "  push $0")))

(defn emit-var-ref [env node]
  (debug :emit-var-ref)
  ;;[:var-ref [:ident "a"]]
  #_(prn :var-ref env node)
  (let [ident (second node)
        offsets (:local-offsets env)
        offset (get offsets ident)]
    (emit env "  lea %d(%%rbp), %%rdi" (- offset))
    (emit env "  mov (%%rdi), %%rax")
    (emit env "  push %%rax")))

(defn emit-return [env node]
  (debug :emit-return)
  ;;[:return expr]
  (emit-expr env (second node))
  (emit env "  pop %%rax")
  (emit env "  jmp .L.return"))

(comment
  (parser "x=2;return x;13;")
  (compile-and-run "x=5;x+1;13;")
  (:local-offsets (build-env (parser "x=2;return x;13;")))
  (compiler "a=2;b=3;b;")
  (compile-and-run "a=2;b=3;a+b;")
  (compiler "a=2;b=3;c=10;a+b+c-3;")
  (parser "a=2;b=3;a+b;")
  ;;
  )

(defn emit-expr [env node]
  (case (first node); type
    :num (emit-num env node)

    :add (emit-binop env node)
    :sub (emit-binop env node)
    :mul (emit-binop env node)
    :div (emit-binop env node)

    :uneg  (emit-unary env node)
    :uplus (emit-unary env node)

    :eq   (emit-rel env node)
    :neq  (emit-rel env node)
    :gt   (emit-rel env node)
    :gte  (emit-rel env node)
    :lt   (emit-rel env node)
    :lte  (emit-rel env node)
    :assign  (emit-assignment env node)
    :var-ref  (emit-var-ref env node)
    :return  (emit-return env node)
    :block  (emit-block env node)
    ))

(defn emit-block [env node]
  ;; [:block expr ...]
  (doseq [expr (rest node)]
    (emit-expr env expr)
    (emit env "  pop %%rax")))


(defn emit-program [env]
  (emit env "  .globl main")
  (emit env "main:")
  (emit env "  push %%rbp")
  (emit env "  mov %%rsp, %%rbp")
  (emit env "  sub  $%d, %%rsp" (:stack-size env))
  (emit-block env (:block env))

  (emit env "  mov $0, %%rax") ; ensure no fallthough makes our tests work
  (emit env ".L.return:")
  (emit env "  mov %%rbp, %%rsp")
  (emit env "  pop %%rbp")
  (emit env "  ret"))


(defn find-locals [ast]
  (case (first ast) 
    :block (apply clojure.set/union
                  (for [s (rest ast)]
                    (find-locals s)))
    :assign #{(-> ast second)}
    #{}))


(defn build-env [stmts]
  (let [locals (find-locals stmts)
        stack-size (* 8 (count locals))]
    {:block stmts
     :stack-size stack-size
     :local-offsets (into {}
                        (map vector
                             locals
                             (range 0 stack-size 8)))}))


(defn compiler [src]
  (let [ast (parser src)]
    (if (insta/failure? ast)
      (throw (Exception. (pr-str (insta/get-failure ast))))
      (let [env (build-env ast)]
        (emit-program env)))))

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
    (is (=  0 (compile-and-run "{return 0;}")))
    (is (= 42 (compile-and-run "{return 42;}"))))
  (testing "addition and subtraction"
    (is (= 21 (compile-and-run "{return 5+20-4;}")))
    (is (= 41 (compile-and-run  " { return 12 + 34 - 5 ; }"))))
  (testing "multiplication"
    (is (=  3 (compile-and-run  " {return 5*3 -12;  }")))
    (is (= 47 (compile-and-run  " {return 5+6*7;}"))))
  (testing "brackets"
    (is (= 15 (compile-and-run "{return 5*(9-6);}")))
    (is (= 15 (compile-and-run "{return 5*(6*2-9);}")))
    (is (= 51 (compile-and-run "{return 5* 6*2-9 ;}"))))
  (testing "division"
    (is (=  4 (compile-and-run "{return (3+5)/2;}")))
    (is (=  4 (compile-and-run "{return (3+10)/3;}"))))
  (testing "unary minus"
    (is (= 10 (compile-and-run "{return -10+20;}")))
    (is (= 10 (compile-and-run "{return - -10;}")))
    (is (= 10 (compile-and-run "{return - - +10;}"))))
  (testing "equality"
    (is (= 0 (compile-and-run "{return 0 ==1;}")))
    (is (= 1 (compile-and-run "{return 42==42;}")))
    (is (= 1 (compile-and-run "{return 0!=1;}")))
    (is (= 0 (compile-and-run "{return 42!=42;}"))))
  (testing "greater than"
    (is (= 1 (compile-and-run "{return 10>5;}")))
    (is (= 0 (compile-and-run "{return 10>50;}")))
    (is (= 0 (compile-and-run "{return 10>=50;}")))
    (is (= 1 (compile-and-run "{return 10>=10;}"))))
  (testing "less than"
    (is (= 1 (compile-and-run "{return 0<1;}")))
    (is (= 0 (compile-and-run "{return 1<1;}")))
    (is (= 0 (compile-and-run "{return 2<1;}")))
    (is (= 1 (compile-and-run "{return 0<=1;}")))
    (is (= 1 (compile-and-run "{return 1<=1;}")))
    (is (= 0 (compile-and-run "{return 2<=1;}"))))
  (testing "multiple expressions"
    (is (= 1 (compile-and-run "{3;2; return 1;}"))))
  (testing "variable assignment"
    (is (= 3 (compile-and-run "{x=2; return 3;}")))
    (is (= 7 (compile-and-run "{x=2;x=3; return 7;}")))
    (is (= 0 (compile-and-run "{x=2;y=3; return 0;}"))))
  (testing "variable references"
    (is (= 3  (compile-and-run "{foo=3; return foo;}")))
    (is (= 8  (compile-and-run "{foo123=3; bar=5; return foo123+bar;}")))
    (is (= 2  (compile-and-run "{x=2; return x;}")))
    (is (= 3  (compile-and-run "{a=2;b=3; return b;}")))
    (is (= 5  (compile-and-run "{a=2;b=3; return a+b;}")))
    (is (= 12 (compile-and-run "{a=2;b=3;c=10; return a+b+c-3;}")))
    (is (= 6  (compile-and-run "{a=2;b=3; return a*b;}"))))
  (testing "early return"
    (is (= 1 (compile-and-run "{ return 1; 2; }")))
    (is (= 7 (compile-and-run "{ a=12; return 5+2; a*2;}"))))
  (testing "blocks"
    (is (= 3 (compile-and-run "{ {1; {2;} return 3;} }")))
    (is (= 2 (compile-and-run "{ {1; {return 2;} return 3;} }")))))
