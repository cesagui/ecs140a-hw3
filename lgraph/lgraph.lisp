;; You may define helper functions here


(defun check_loop (g start_edges term sequence)
    (if (null start_edges)
        nil
        (let* ((curr_edge       (car start_edges))
               (curr_edge_dest  (car (cdr curr_edge))) ;; dest is second element
               (curr_edge_label (car curr_edge))
               (sequence_start  (car sequence)))
            (if (and (equal curr_edge_label sequence_start)
                     (check g curr_edge_dest term (cdr sequence))) ;; check that the label and sequence start match. call
                t
                (check_loop g (cdr start_edges) term sequence)))))
                
;; create a helper function: that follows a sequence of runes through graph g to see if the path exists
(defun check (g start term sequence)
    (let ((start_edges (funcall g start))) ;; retrieve outgoing edges from start node
        (cond
            ((null sequence)    (and start_edges (equal start term))) ;; empty path is only valid when the node exists
            ((null start_edges) nil) ;; if no outgoing edges, return nil
            ((null (car start_edges)) nil) ;; handle nil  for no edges
            ((check_loop g start_edges term sequence) t) ;; loop thru outgoing edges, check to see if one matches first label
            (t nil)))) ;; base case

(defun find_loop (edges g1 g2 c_node t_node k s prefix)
    (if (null edges)
        nil
        (let* ((curr_edge       (car edges))
               (curr_edge_dest  (car (cdr curr_edge))) ;; dest is second element
               (curr_edge_label (car curr_edge)))
            (let ((recursive_call ;; check at find at current level
                   (find-path g1 g2 curr_edge_dest t_node (- k 1) s (append prefix (list curr_edge_label)))))
                (if recursive_call
                    recursive_call
                    (find_loop (cdr edges) g1 g2 c_node t_node k s prefix)))))) ;; iterate to the next edge

(defun find-path (g1 g2 c_node t_node k s prefix)
    (let ((edges (funcall g1 c_node))) ;; get outgoing edges from current node c
        (cond
            ((zerop k) ;; base case: check that current node matches end node
                (if (and (equal c_node t_node)
                         (not (check g2 s t_node prefix)))
                    (cons prefix t) ;; return the sequence paired with t
                    nil))
            ((null edges) nil) ;; check that the outgoing edges from C exist
            ((null (car edges)) nil) ;; handle nil for no edges
            ;; loop through the edges and recursively call find on them, building up the prefix as we go
            ((find_loop edges g1 g2 c_node t_node k s prefix))
            (t nil))))

(defun find-sequence (g1 g2 start target k)
    (find-path g1 g2 start target k start '()))