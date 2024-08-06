#!/bin/sh

tests() {
    section "constants"
    assert_return  0  "fun main() { return 0; }"
    assert_return  42  "fun main() { return 42; }"

    section "add_sub"
    assert_return 21 "fun main() {return 5+20-4;}"
    assert_return 41 "fun main() { return 12 + 34 - 5 ; }"

    section "multiplication"
    assert_return 3  " fun main() {return 5*3 -12;  }"
    assert_return 47 " fun main() {return 5+6*7;}"

    section "brackets"
    assert_return 15  "fun main() {return 5*(9-6);}"
    assert_return 39  "fun main() {return (5*9)-6;}"
    assert_return 15  "fun main() {return 5*(6*2-9);}"
    assert_return 51  "fun main() {return 5* 6*2-9 ;}"

    section "division"
    assert_return 4  "fun main() {return (3+5)/2;}"
    assert_return 4  "fun main() {return (3+10)/3;}"

    section "unary minus"
    assert_return 10  "fun main() {return -10+20;}"
    assert_return 10  "fun main() {return - -10;}"
    assert_return 10  "fun main() {return - - +10;}"

    section "equality"
    assert_return 0  "fun main(){return 0 ==1;}"
    assert_return 1  "fun main(){return 42==42;}"
    assert_return 1  "fun main(){return 0!=1;}"
    assert_return 0  "fun main(){return 42!=42;}"

    section "greater than"
    assert_return 1  "fun main(){return 10>5;}"
    assert_return 0  "fun main(){return 10>50;}"
    assert_return 0  "fun main(){return 10>=50;}"
    assert_return 1  "fun main(){return 10>=10;}"

    section "less than"
    assert_return 1  "fun main(){return 0<1;}"
    assert_return 0  "fun main(){return 1<1;}"
    assert_return 0  "fun main(){return 2<1;}"
    assert_return 1  "fun main(){return 0<=1;}"
    assert_return 1  "fun main(){return 1<=1;}"
    assert_return 0  "fun main(){return 2<=1;}"

    section "multiple expressions"
    assert_return 1  "fun main(){3;2; return 1;}"
    assert_return 7  "fun main(){3;2; return 7;}"

    section "variable assignment"
    assert_return 3  "fun main(){thing=2; return 3;}"
    assert_return 7  "fun main(){x=2;x=3; return 7;}"
    assert_return 0  "fun main(){x=2;y=3; return 0;}"

    section "variable references"
    assert_return 3  "fun main(){foo=3; return foo;}"
    assert_return 8  "fun main(){foo123=3; bar=5; return foo123+bar;}"
    assert_return 2  "fun main(){x=2; return x;}"
    assert_return 3  "fun main(){a=2;b=3; return b;}"
    assert_return 5  "fun main(){a=2;b=3; return a+b;}"
    assert_return 12  "fun main(){a=2;b=3;c=10; return a+b+c-3;}"
    assert_return 6  "fun main(){a=2;b=3; return a*b;}"

    section "early return"
    assert_return 1  "fun main(){ return 1; 2; }"
    assert_return 7  "fun main(){ a=12; return 5+2; a*2;}"

    section "blocks"
    assert_return 3  "fun main(){ {1; {2;}; return 3;}; }"
    assert_return 2  "fun main(){ {1; {return 2;}; return 3;}; }"

    section "unused semi-colons work"
    assert_return 3  "fun main(){ ;;; return 3;}"

    section "if statements"
    assert_return 3  "fun main(){ if (0) return 2; return 3; }"
    assert_return 3  "fun main(){ if (1-1) return 2; return 3; }"
    assert_return 2  "fun main(){ if (1) return 2; return 3; }"
    assert_return 2  "fun main(){ if (2-1) return 2; return 3; }"
    assert_return 4  "fun main(){ if (0) { 1; 2; return 3; } else { return 4; } }"
    assert_return 3  "fun main(){ if (1) { 1; 2; return 3; } else { return 4; } }"

    #section "multiple if statements"
    #assert_return 3  "fun main(){ if (1) if (0) return 5; return 3; }"
    #assert_return 5  "fun main(){ if (1) return 5; if (0) return 3; }"

    #section "for statement"
    #assert_return 55 "fun main(){ i=0; j=0; for (i=0; i<=10; i=i+1) j=i+j; return j; }"
    #assert_return 3 "fun main(){ for (;;) {return 3;} return 5; }"

    #section "while statement"
    #assert_return 15 "fun main(){ i=0; j=0; while (i < 5) { i = i+1; j=j+3; } return j; }"
    #assert_return 3 "fun main(){ while(1) {return 3;} return 5; }"

    #section "references"
    #assert_return 4 "fun main(){ i=4; return *&i; }"
    #assert_return 9 "fun main(){ i=4; j=&i; return *j + 5; }"
}


setup() {

    basedir=$(pwd)/_testtmp
    mkdir -p $basedir

    #basedir=$(mktemp -d /tmp/compy.XXXXX)
    echo 'Base Dir: ' $basedir
}

cleanup() {
    rm -rf $basedir
}

section() {
    test_section=$(echo $1 | tr '[ \t]' '_')
    test_number=0
}

assert_return() {
    test_number=$(($test_number+1))
    ret="$1"
    inline_code="$2"
    src=$basedir/$test_section.$test_number.lx
    bin=$src.exe


    echo "$inline_code" > $src
    if ! ./compy-aux compile $src; then
	echo "FAIL: Cant compile"
	return 1
    fi
    $bin
    act="$?"
    if [ $ret -eq $act ]; then
	echo -n '.'
	#echo "OK"
	rm $src $src.*
    else
	echo ''
	echo -n $test_section $test_number " ..."
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
