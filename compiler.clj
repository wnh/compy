
(ns compiler
  (:require [instaparse.core :as insta]
            [clojure.test :as test :refer [deftest testing is]]))

(def ^:dynamic *debugging* false)

(def parser
  (insta/parser
   "<program> = block

    block =   <'{'>   stmt (  stmt)*   <'}'>  

    <stmt> = simple-stmt <semi>
           | return <semi>
           | block
           | if
           | for
           | while
    if = <kw-if>   expr   stmt   (<kw-else>   stmt)?

    for = <kw-for>   <'('> simple-stmt  <semi> simple-stmt <semi> simple-stmt <')'> stmt

    while = <kw-while>   <'('> expr <')'>   stmt

    <simple-stmt> = assign | expr | empty-stmt

    empty-stmt =   ''  

    assign = ident <'='> expr
    return =   <kw-ret>   expr
    <expr> = eq-expr

    <eq-expr> = eq | neq | rel-expr
    eq  = rel-expr <'=='> rel-expr
    neq = rel-expr <'!='> rel-expr

    <rel-expr> = gt | gte | lt | lte | sum-expr
    gt  = sum-expr   <'>'>   rel-expr
    gte = sum-expr   <'>='>   rel-expr
    lt  = sum-expr   <'<'>   rel-expr
    lte = sum-expr   <'<='>   rel-expr

    <sum-expr> = add | sub | term-expr
    add  = term-expr <'+'> sum-expr
    sub  = term-expr <'-'> sum-expr

    <term-expr> = mul | div | factor-expr
    mul = factor-expr <'*'> term-expr
    div = factor-expr <'/'> term-expr

    <factor-expr> = num | unary | var-ref
                  | <'('> expr <')'>
    <unary> = uneg | addr-of | deref | uplus
    var-ref = ident
    uneg      =   <'-'>   factor-expr
    addr-of   =   <'&'>   factor-expr
    deref     =   <'*'>   factor-expr
    uplus    =   <'+'>   factor-expr
    num      =   #'[0-9]+'  
    ident    =   !keyword #'[a-zA-Z][a-zA-Z0-9]*'  
    <keyword> = kw-ret | kw-if
    <kw-else>  = 'else'
    <kw-for>  = 'for'
    <kw-while>  = 'while'
    <kw-if>  = 'if'
    <kw-ret> = 'return'
    <semi> =   <';'>"
   :auto-whitespace :standard
   :output-format :enlive))


(defn debug [& args]
  (when *debugging*
    (binding [*out* *err*]
      (apply prn args))))

(defn prn-ast [node]
  (case (:tag node)
    :num   [:num   (-> node :content first) ]
    :ident [:ident (-> node :content first) ]
    (into [(:tag node)]
          (map prn-ast (:content node)))))

(comment
  (binding [*debugging* true])
  (parser "if (123/2) {return 123;}" :start :if)
  (compiler "{return 42;}")
  (parser "{return 42;}")
  (insta/parses parser "{if 6}" :partial true)
  (find-locals (first (parser "{x=2;y=3; return 0;}")))
  (compiler "{x=2;y=3; return 0;}")
  (parser " { return 12 + 34 - 5 ; }")
  (take 2 (insta/parses parser "{ {1; {2;} return 3;} }"))
  (parser "{ {1; {2;} return 3;} }")
  (tokens "{ {1; {2;} return 3;} }")
  ;; assert 3 '{ {1; {2;} return 3;} }'
  )

(defn bail [token]
  (throw (Exception. "Unexpeccted: " (pr-str token))))

(defn emit [env & args]
  (let [asm (:asm env)
        line (apply format args)]
    (swap! asm (fn [s] (str s line "\n")))))


(declare emit-expr)

(defn emit-num [env node]
  (debug :emit-num)
  ;; [:num "123"]
  (emit env "  mov $%d, %%rax" (Integer/parseInt (-> node :content first))))

(defn emit-unary [env node]
  (debug :emit-unary)
  ;; [:uneg  "123"]
  ;; [:uplus "123"]
  (case (:tag node)
    :uneg (do (emit-expr env (-> node :content first))
              (emit env "  neg %%rax"))
    :uplus (do (emit-expr env (-> node :content first)))))

(defn emit-rel [env node]
  (debug :emit-rel)
  (let [[left right] (:content node)]
    (emit-expr env left)
    (emit env "  push %%rax")
    (emit-expr env right)
    (emit env "  pop %%rbx"))
  (case (:tag node)
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
             (emit env "  setle %%al"))))

(defn emit-binop [env node]
  (let [[left right] (:content node)]
    (emit-expr env right)
    (emit env "  push %%rax")
    (emit-expr env left)
    (emit env "  pop %%rbx")
    (case (:tag node)
      :sub (emit env "  sub %%rbx, %%rax")
      :add (emit env "  add %%rbx, %%rax")
      :mul (emit env "  imul %%rbx, %%rax")
      :div (do (emit env "  xor %%rdx, %%rdx")
               (emit env "  div %%rbx")))))

(defn emit-assignment [env node]
  (debug :emit-assignment env)
  (emit env "  /* start assignment */")
  ;;[:assign [:ident "a"] [:num "2"]]
  (let [[ident expr] (-> node :content)
        offsets (:local-offsets env)
        offset (get offsets ident)]
    (emit-expr env expr)
    (emit env "  lea %d(%%rbp), %%rdi" (- offset))
    (emit env "  mov %%rax, (%%rdi)"))
  (emit env "  /* end   assignment */"))

(defn emit-var-ref [env node]
  (debug :emit-var-ref)
  (emit env "  /* start var-ref */")
  ;;[:var-ref [:ident "a"]]
  #_(prn :var-ref env node)
  (let [ident (-> node :content first)
        offsets (:local-offsets env)
        offset (get offsets ident)]
    (emit env "  lea %d(%%rbp), %%rdi" (- offset))
    (emit env "  mov (%%rdi), %%rax"))
  (emit env "  /* end   var-ref */"))

(defn emit-return [env node]
  (debug :emit-return)
  ;;[:return expr]
  (emit-expr env (-> node :content first))
  (emit env "  jmp .L.return"))

(comment
  (parser "x=2;return x;13;")
  (compile-and-run "x=5;x+1;13;")
  (:local-offsets (build-env (parser "x=2;return x;13;")))
  (compiler "a=2;b=3;b;")
  (compile-and-run "a=2;b=3;a+b;")
  ;(binding [*debugging* true]
    (compile-and-run "{ if (0) return 9; return 7; }")
   ; )
  ;;
  )

(defn next-label [env]
  (swap! (:counter env) inc))

(defn emit-if [env node]
  ;; [:if test block]
  (debug :emit-if node)
  (emit env "  /* start if */")
  (let [[test then else-branch] (:content node)
        n (next-label env)]
    (emit env "  /* start if-test-expr */")
    (emit-expr env test)
    (emit env "  /* end   if-test-expr */")
    (emit env "  cmp $0, %%rax")
    (emit env "  jz .L.else.%d" n)
    (emit-expr env  then)
    (emit env "  jmp .L.done.%d" n)
    (emit env ".L.else.%d:" n)
    (emit-expr env (if else-branch
                     else-branch
                     {:tag :block :content '()}))
    (emit env ".L.done.%d:" n))
  (emit env "  /* end if */"))

(defn emit-for [env node]
  (let [[init-expr test-expr inc-expr body] (:content node)
        n (next-label env)]
    (emit-expr env init-expr)
    (emit env ".L.For.Top.%d:" n)
    (emit-expr env test-expr)
    (emit env "  cmp $0, %%rax")
    (emit env "  jz .L.For.End.%d" n)
    (emit-expr env body)
    (emit-expr env inc-expr)
    (emit env "  jmp .L.For.Top.%d" n)
    (emit env ".L.For.End.%d:" n)))

(defn emit-while [env node]
  (let [[test-expr body] (:content node)
        n (next-label env)]
    (emit env ".L.While.Top.%d:" n)
    (emit-expr env test-expr)
    (emit env "  cmp $0, %%rax")
    (emit env "  jz .L.While.End.%d" n)
    (emit-expr env body)
    (emit env "  jmp .L.While.Top.%d" n)
    (emit env ".L.While.End.%d:" n)))

(defn emit-block [env node]
  (debug :emit-block node)
  ;; [:block expr ...]
  (doseq [expr (:content node)]
    (emit-expr env expr)))


(defn emit-deref [env node]
  (let [[expr] (:content node)]
    (emit-expr env expr)
    (emit env "  mov %%rax, %%rdi")
    (emit env "  mov (%%rdi), %%rax")))

(defn is-lval? [node]
  (let [expr-type (-> node :content first :tag)]
    (= expr-type :var-ref)))

(defn emit-addr-of [env node]
  (let [expr-type (-> node :content first :tag)]
    (if (not (is-lval? node))
      (throw (Exception. (str "TypeError: Can't take reference to " expr-type))))
    (let [ident (-> node :content first :content first)
          offsets (:local-offsets env)
          offset (get offsets ident)]
      (debug :addr-of ident offsets offset)
      (emit env "  lea %d(%%rbp), %%rax" (- offset)))))


(defn emit-expr [env node]
  (case (:tag node); type
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
    :if (emit-if env node)
    :for (emit-for env node)
    :while (emit-while env node)
    :deref (emit-deref env node)
    :addr-of (emit-addr-of env node)
    :empty-stmt nil
    ))



(defn emit-program [env]
  (debug :emit-program env)
  (emit env "/*")
  (emit env (with-out-str (clojure.pprint/pprint (prn-ast (:block env)))))
  (emit env "*/")
  (emit env "  /* start emit-program */")
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
  (emit env "  ret")
  (emit env "  /* end emit-program */"))


(defn first-type
  "grab the type of the first child node"
  [node]
  (-> node first :type))

(defn type-check
  "recursivly check the AST for type correctness and apply a :type
  field to each node"
  [node]
  node
  #_(if (not (map? node))
    node
    (let [typed-content (map type-check (:content node))
          node-type (case (:tag node)
                      :addr-of (do (if (not= (-> typed-content first :tag) :var-ref)
                                     (throw (ex-info
                                             (str "TypeError: Can't take reference to "
                                                  (-> typed-content first :tag))
                                             {})))
                                 [:ptr (first-type typed-content)])
                      :deref (do (if (not= (-> typed-content first :type first) :ptr)
                                     (throw (ex-info
                                             (str "TypeError: Can't dereference a "
                                                  (-> typed-content first :tag))
                                             {})))
                                 (second (first-type typed-content)))
                      :return (first-type typed-content)
                      ;; else
                      ;; TODO right now everything is an int. This will need to change
                      :i64)]
      (assoc node :type node-type
             :content typed-content))))

(comment 
  (compile-and-run
   "{ i=4; j=&i; return *j + 5; }")

  (compile-and-run
    "{i=&0;}")

  (testing "Cant deref constants"
    (is (thrown? Exception (-> "{i=0;j=*i;}" compiler))))

  (compile-and-run "{i=0;j=*i;}")

  ;; TODO: make the local vars offset happen incrementally so we can't
  ;; see all variables in all scopes right away
  #_(testing "need to declar a variable before its used"
    (is (thrown? Exception (-> "{i; i=0;}" compiler))))
  ;;
  )


(defn find-locals [ast]
  (debug :find-locals ast)
  (case (:tag ast)
    :block (apply clojure.set/union
                  (for [s (:content ast)]
                    (find-locals s)))
    :assign #{(-> ast :content first)}
    #{}))


(defn build-env [stmts]
  (debug :build-env stmts)
  (let [locals (find-locals stmts)
        stack-size (* 8 (count locals))]
    {:block stmts
     :asm (atom "")
     :counter (atom 0)
     :stack-size stack-size
     :local-offsets (into {}
                        (map vector
                             locals
                             (range 0 stack-size 8)))}))

(defn compiler [src]
  ;; TODO parser returns seq at top level
  (let [ast (first (parser src))]
    (if (insta/failure? ast)
      (throw (Exception. (pr-str (insta/get-failure ast))))
      (let [checked-ast (type-check ast)
            env (build-env checked-ast)] 
        (debug :initial-ast ast)
        (debug :initial-env env)
        (emit-program env)
        @(:asm env)))))

(defn -main []
  (let [src (first *command-line-args*)]
    (compiler src)))


(defn compile-and-run [src] 
  (let [asm (compiler src)]
    (spit "tmp.s" asm)
    (let [compile-out (clojure.java.shell/sh "gcc" "-static" "-g" "-o" "tmp" "tmp.s")]
      (if (< 0 (:exit compile-out))
        (throw (Exception. (:err compile-out)))
        (let [out (clojure.java.shell/sh "./tmp")]
          (prn out)
          (:exit out))))))

(defmacro assert-return [ret code]
  `(do (is (not  (insta/failure? (parser ~code))))
       (is (=  ~ret (compile-and-run ~code)))))

(deftest  test-compiler-works
  (testing "Constants"
    (assert-return 0  "{return 0;}")
    (assert-return 42  "{return 42;}"))
  (testing "addition and subtraction"
    (assert-return 21  "{return 5+20-4;}")
    (assert-return 41   " { return 12 + 34 - 5 ; }"))
  (testing "multiplication"
    (assert-return 3   " {return 5*3 -12;  }")
    (assert-return 47   " {return 5+6*7;}"))
  (testing "brackets"
    (assert-return 15  "{return 5*(9-6);}")
    (assert-return 15  "{return 5*(6*2-9);}")
    (assert-return 51  "{return 5* 6*2-9 ;}"))
  (testing "division"
    (assert-return 4  "{return (3+5)/2;}")
    (assert-return 4  "{return (3+10)/3;}"))
  (testing "unary minus"
    (assert-return 10  "{return -10+20;}")
    (assert-return 10  "{return - -10;}")
    (assert-return 10  "{return - - +10;}"))
  (testing "equality"
    (assert-return 0  "{return 0 ==1;}")
    (assert-return 1  "{return 42==42;}")
    (assert-return 1  "{return 0!=1;}")
    (assert-return 0  "{return 42!=42;}"))
  (testing "greater than"
    (assert-return 1  "{return 10>5;}")
    (assert-return 0  "{return 10>50;}")
    (assert-return 0  "{return 10>=50;}")
    (assert-return 1  "{return 10>=10;}"))
  (testing "less than"
    (assert-return 1  "{return 0<1;}")
    (assert-return 0  "{return 1<1;}")
    (assert-return 0  "{return 2<1;}")
    (assert-return 1  "{return 0<=1;}")
    (assert-return 1  "{return 1<=1;}")
    (assert-return 0  "{return 2<=1;}"))
  (testing "multiple expressions"
    (assert-return 1  "{3;2; return 1;}"))
  (testing "variable assignment"
    (assert-return 3  "{x=2; return 3;}")
    (assert-return 7  "{x=2;x=3; return 7;}")
    (assert-return 0  "{x=2;y=3; return 0;}"))
  (testing "variable references"
    (assert-return 3  "{foo=3; return foo;}")
    (assert-return 8  "{foo123=3; bar=5; return foo123+bar;}")
    (assert-return 2  "{x=2; return x;}")
    (assert-return 3  "{a=2;b=3; return b;}")
    (assert-return 5  "{a=2;b=3; return a+b;}")
    (assert-return 12  "{a=2;b=3;c=10; return a+b+c-3;}")
    (assert-return 6  "{a=2;b=3; return a*b;}"))
  (testing "early return"
    (assert-return 1  "{ return 1; 2; }")
    (assert-return 7  "{ a=12; return 5+2; a*2;}"))
  (testing "blocks"
    (assert-return 3  "{ {1; {2;} return 3;} }")
    (assert-return 2  "{ {1; {return 2;} return 3;} }"))
  (testing "unused semi-colons work"
    (assert-return 3  "{ ;;; return 3;}"))
  (testing "if statements"
    (assert-return 3  "{ if (0) return 2; return 3; }")
    (assert-return 3  "{ if (1-1) return 2; return 3; }")
    (assert-return 2  "{ if (1) return 2; return 3; }")
    (assert-return 2  "{ if (2-1) return 2; return 3; }")
    (assert-return 4  "{ if (0) { 1; 2; return 3; } else { return 4; } }")
    (assert-return 3  "{ if (1) { 1; 2; return 3; } else { return 4; } }"))
  (testing "multiple if statements"
    (assert-return 3  "{ if (1) if (0) return 5; return 3; }")
    (assert-return 5  "{ if (1) return 5; if (0) return 3; }"))
  (testing "for statement"
    (assert-return 55 "{ i=0; j=0; for (i=0; i<=10; i=i+1) j=i+j; return j; }")
    (assert-return 3 "{ for (;;) {return 3;} return 5; }"))
  (testing "while statement"
    (assert-return 15 "{ i=0; j=0; while (i < 5) { i = i+1; j=j+3; } return j; }")
    (assert-return 3 "{ while(1) {return 3;} return 5; }"))
  (testing "references"
    (assert-return 4 "{ i=4; return *&i; }")
    (assert-return 9 "{ i=4; j=&i; return *j + 5; }")))


(deftest  test-typechecking
  (testing "Cant take reference of constants"
    (is (thrown? Exception (-> "{i=&0;}"
                               compiler))))
  (testing "Cant deref constants"
    (is (thrown? Exception (-> "{i=0;j=*i;}" compiler))))
  ;; TODO: make the local vars offset happen incrementally so we can't
  ;; see all variables in all scopes right away
  #_(testing "need to declar a variable before its used"
    (is (thrown? Exception (-> "{i; i=0;}" compiler)))))

;TODO Better Error handling some how, maybe return an ast node with :error on it
