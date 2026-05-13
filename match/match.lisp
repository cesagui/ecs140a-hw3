;; helper function walks through the assertion trying to match the rest of the pattern at each suffix, consuming at least one element

(defun helper (rest-pattern assertion)
  (cond
    ((null assertion) nil)
    ((match rest-pattern assertion) t)
    (t (helper rest-pattern (cdr assertion)))))

(defun match (pattern assertion)
  (cond
    ((and (null pattern) (null assertion)) t) ;; exhausted both lists return true
    ((null pattern) nil) ;; exhausted assertion, return nil

    ((eq (car pattern) '!) ;; ! case
     (let ((remain-pattern (cdr pattern))) ;; def remaining pattern as the cdr
       (if (null remain-pattern) ;; if remaining pattern has been exhausted
           (not (null assertion)) ;; if the assertion has been exhausted => return T ; OW : return false
           (helper remain-pattern (cdr assertion))))) ;; OW: call helper function on remain-pattern and the cdr

    ((null assertion) nil) ;; exhausted assertion but remaining pattern, return nil

    ((eq (car pattern) '?) ;; wild card, just move forward
     (match (cdr pattern) (cdr assertion)))

    ((equal (car pattern) (car assertion)) ;; match, just move forward
     (match (cdr pattern) (cdr assertion)))

    (t nil))) ;; none hold ? return nil (remaining non-matching pattern and assertion)