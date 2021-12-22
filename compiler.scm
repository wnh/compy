
(import (scheme small)
	(srfi 115)
	(srfi 166))

(define (read-token token)
  (call-with-port (open-input-string token)
		  (lambda (p) (read p))))

(define (parse)
  (let* ((l (read-line))
	 (tokens (regexp-partition '(or "+" "-") l)))
    (map read-token tokens)))

(define (emit input)
  (show #t "  .globl main\n")
  (show #t "main:\n")
  (show #t "  mov $" (car input) ", %rax\n")
  (let loop ((ins (cdr input)))
    (cond 
     ((null? ins) #f)
     ((equal? (car ins) '+)
      (begin (show #t "  add $" (cadr ins) ", %rax\n")
	     (loop (cddr ins))))
      ((equal? (car ins) '-)
       (begin (show #t "  sub $" (cadr ins) ", %rax\n")
	      (loop (cddr ins))))))
  (show #t "ret\n"))

(define (main)
  (let ((input (parse)))
    (emit input)))

(main)
