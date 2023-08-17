#!/bin/sh

tests() {
    echo "testing Constants"
    assert_return "const1" 0  "fun main() { return 0; }"
    assert_return "const2" 42  "fun main() { return 42; }"

    echo "testing addition and subtraction"
    assert_return "addsub1" 21 "fun main() {return 5+20-4;}"
    assert_return "addsub2" 41 "fun main() { return 12 + 34 - 5 ; }"

    #echo "testing multiplication"
    #assert_return 3   " {return 5*3 -12;  }"
    #assert_return 47   " {return 5+6*7;}"

    #echo "testing brackets"
    #assert_return 15  "{return 5*(9-6);}"
    #assert_return 15  "{return 5*(6*2-9);}"
    #assert_return 51  "{return 5* 6*2-9 ;}"

    #echo "testing division"
    #assert_return 4  "{return (3+5)/2;}"
    #assert_return 4  "{return (3+10)/3;}"

    #echo "testing unary minus"
    #assert_return 10  "{return -10+20;}"
    #assert_return 10  "{return - -10;}"
    #assert_return 10  "{return - - +10;}"

    #echo "testing equality"
    #assert_return 0  "{return 0 ==1;}"
    #assert_return 1  "{return 42==42;}"
    #assert_return 1  "{return 0!=1;}"
    #assert_return 0  "{return 42!=42;}"

    #echo "testing greater than"
    #assert_return 1  "{return 10>5;}"
    #assert_return 0  "{return 10>50;}"
    #assert_return 0  "{return 10>=50;}"
    #assert_return 1  "{return 10>=10;}"

    #echo "testing less than"
    #assert_return 1  "{return 0<1;}"
    #assert_return 0  "{return 1<1;}"
    #assert_return 0  "{return 2<1;}"
    #assert_return 1  "{return 0<=1;}"
    #assert_return 1  "{return 1<=1;}"
    #assert_return 0  "{return 2<=1;}"

    #echo "testing multiple expressions"
    #assert_return 1  "{3;2; return 1;}"

    #echo "testing variable assignment"
    #assert_return 3  "{x=2; return 3;}"
    #assert_return 7  "{x=2;x=3; return 7;}"
    #assert_return 0  "{x=2;y=3; return 0;}"

    #echo "testing variable references"
    #assert_return 3  "{foo=3; return foo;}"
    #assert_return 8  "{foo123=3; bar=5; return foo123+bar;}"
    #assert_return 2  "{x=2; return x;}"
    #assert_return 3  "{a=2;b=3; return b;}"
    #assert_return 5  "{a=2;b=3; return a+b;}"
    #assert_return 12  "{a=2;b=3;c=10; return a+b+c-3;}"
    #assert_return 6  "{a=2;b=3; return a*b;}"

    #echo "testing early return"
    #assert_return 1  "{ return 1; 2; }"
    #assert_return 7  "{ a=12; return 5+2; a*2;}"
    #
    #echo "testing blocks"
    #assert_return 3  "{ {1; {2;} return 3;} }"
    #assert_return 2  "{ {1; {return 2;} return 3;} }"

    #echo "testing unused semi-colons work"
    #assert_return 3  "{ ;;; return 3;}"

    #echo "testing if statements"
    #assert_return 3  "{ if (0) return 2; return 3; }"
    #assert_return 3  "{ if (1-1) return 2; return 3; }"
    #assert_return 2  "{ if (1) return 2; return 3; }"
    #assert_return 2  "{ if (2-1) return 2; return 3; }"
    #assert_return 4  "{ if (0) { 1; 2; return 3; } else { return 4; } }"
    #assert_return 3  "{ if (1) { 1; 2; return 3; } else { return 4; } }"

    #echo "testing multiple if statements"
    #assert_return 3  "{ if (1) if (0) return 5; return 3; }"
    #assert_return 5  "{ if (1) return 5; if (0) return 3; }"

    #echo "testing for statement"
    #assert_return 55 "{ i=0; j=0; for (i=0; i<=10; i=i+1) j=i+j; return j; }"
    #assert_return 3 "{ for (;;) {return 3;} return 5; }"

    #echo "testing while statement"
    #assert_return 15 "{ i=0; j=0; while (i < 5) { i = i+1; j=j+3; } return j; }"
    #assert_return 3 "{ while(1) {return 3;} return 5; }"

    #echo "testing references"
    #assert_return 4 "{ i=4; return *&i; }"
    #assert_return 9 "{ i=4; j=&i; return *j + 5; }"
}


setup() {
    basedir=$(mktemp -d /tmp/compy.XXXXX)
}

cleanup() {
    rm -rf $basedir
}

assert_return() {
    name="$1"
    ret="$2"
    inline_code="$3"
    src=$basedir/test.$name.lx
    bin=$basedir/test.$name.lx.exe
    echo "$inline_code" > $src
    if ! ./compy-aux compile $src; then
	echo "FAIL: Cant compile"
	return 1
    fi
    $bin
    act="$?"
    if [ $ret -eq $act ]; then
	echo "OK"
    else
	echo "FAIL"
	echo "Want: $ret"
	echo " Got: $act"
	echo " src: $2"
    fi
}

if [ "${BASIC_TEST_DEBUG:-0}" -eq 1 ]
then
    set -x
fi
setup && tests
#cleanup
