(import (scheme small))

(define input (read))

;(display input)
;(display "\n")
(write-string (string-append "
    .globl main\n
main:
    mov $" (number->string input) ", %rax
    ret
"))

