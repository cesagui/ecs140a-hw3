package sexpr
 
import (
	"fmt"
	"errors"
	"math/big"
)
 
// ErrEval is the error value returned by the Evaluator if the contains
// an invalid token.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrEval = errors.New("eval error")

// checks if an sexpr represents nil, either as the empty list or the NIL symbol
func (expr *SExpr) isNilValue() bool {
	if expr == nil {
		return false
	}
	// check for it's the empty list (&SExpr{})
	if expr.isNil() {
		return true
	}
	// check symbol "nil"
	if expr.isSymbol() && expr.atom.literal == "NIL" {
		return true
	}
	return false
}

 
func (expr *SExpr) Eval() (*SExpr, error) {
	// check that expr is valid
	if expr == nil {
		return nil, ErrEval
	}

	if expr.isNilValue() {
		return mkNil(), nil
	}
 
	if expr.isNumber() {
		return expr.Number()
	}
 
	if expr.isConsCell() {
		// check car is def
		car, err := expr.Car()
		if err != nil {
			return nil, ErrEval
		}
 
		if car.isSymbol() {
			switch car.atom.literal {
			case "QUOTE":
				return expr.Quote()
			case "CAR":
				return expr.CarFunc()
			case "CDR":
				return expr.CdrFunc()
			case "CONS":
				return expr.Cons()
			case "LENGTH":
				return expr.Length()
			default:
				return nil, ErrEval
			}
		}
	}
 
	return nil, ErrEval
}
 
func (expr *SExpr) Quote() (*SExpr, error) {
	// get arg by taking the cdr (arg . NIL)
	args, err := expr.Cdr()
	if err != nil || args == nil || args.isNil(){
		return nil, ErrEval
	}
 
	// args must be a cons cell
	if !args.isConsCell() {
		return nil, ErrEval
	}
 
	// verify that only one arg is passed
	argsCdr, err := args.Cdr()
	if err != nil {
		return nil, ErrEval
	}
 
	// check args CDR is nil VAL (empty list or NIL sym)
	if !argsCdr.isNilValue() {
		return nil, ErrEval
	}
 
	// get the first element (the thing to quote)
	quoteExpr, err := args.Car()
	if err != nil {
		return nil, ErrEval
	}
 
	return quoteExpr, nil
}
 
func (expr *SExpr) Number() (*SExpr, error) {
	// we previously checked that we have a number
	return expr, nil
}
 
func (expr *SExpr) CarFunc() (*SExpr, error) {
	// expr is (CAR <arg>)
	args, err := expr.Cdr()
	if err != nil {
		return nil, ErrEval
	}

	// verify exactly one argument
	argsCdr, err := args.Cdr()
	if err != nil {
		return nil, ErrEval
	}

	if argsCdr == nil || !argsCdr.isNil() {
		return nil, ErrEval
	}


	// get the argument and evaluate it
	arg, err := args.Car()
	if err != nil {
		return nil, ErrEval
	}

	argVal, err := arg.Eval()
	if err != nil {
		return nil, ErrEval
	}
	if argVal.isNil() {
		return mkNil(), nil
	}

	// return the car of the evaluated result
	return argVal.Car()
}

func (expr *SExpr) CdrFunc() (*SExpr, error) {
	// expr is (CDR <arg>)
	args, err := expr.Cdr()
	if err != nil {
		return nil, ErrEval
	}

	// verify exactly one argument
	argsCdr, err := args.Cdr()
	if err != nil {
		return nil, ErrEval
	}

	if argsCdr == nil || !argsCdr.isNil() {
		return nil, ErrEval
	}

	// get the argument and evaluate it
	arg, err := args.Car()
	if err != nil {
		return nil, ErrEval
	}

	argVal, err := arg.Eval()
	if err != nil {
		return nil, ErrEval
	}
	if argVal.isNil() {
		return mkNil(), nil
	}

	// return the cdr of the evaluated result
	return argVal.Cdr()
}
 
func (expr *SExpr) Cdr() (*SExpr, error) {
	// return the CDR of this cell
	return expr.cdr, nil
}
 
func (expr *SExpr) Car() (*SExpr, error) {
	// return the CAR of this cell
	return expr.car, nil
}

func (expr *SExpr) Cons() (*SExpr, error) {
	// get the first arg by taking the cdr (arg1 arg2 . NIL)
	arg1, err := expr.Cdr()
	fmt.Println(arg1)
	if err != nil || arg1 == nil || arg1.isNil(){
		return nil, ErrEval
	}
	// get the CAR of cell
	arg1Cell, err := arg1.Car()
	if err != nil || arg1Cell == nil || arg1Cell.isNil(){
		return nil, ErrEval
	}
	// eval arg1
	arg1Eval, _ := arg1Cell.Eval()

	
	arg2, err := arg1.Cdr()
	if err != nil {
		fmt.Println(2)
		return nil, ErrEval
	}
	arg2Cell, err := arg2.Car()
	if err != nil || arg2Cell == nil || arg2Cell.isNil(){
		return nil, ErrEval
	}
	arg2Eval, _ := arg2Cell.Eval()

	// check args2 CDR is nil VAL (empty list or NIL sym)
	argsCdr, err := arg2.Cdr()
	if err != nil {
		fmt.Println(3)
		return nil, ErrEval
	}
	if !argsCdr.isNilValue() {
		fmt.Println(4)
		return nil, ErrEval
	}
	return mkConsCell(arg1Eval, arg2Eval), nil
}
func (expr *SExpr) lengthHelper() (int64, error) {
	if expr.isNilValue() {
		return 0, nil
	}
	// try to grab the cdr of this cell 
	exprCdr, err := expr.Cdr()
	if err != nil {
		return 0, ErrEval
	}
	recursiveCall, err := exprCdr.lengthHelper()
	if err != nil {
		return 0, ErrEval
	}
	return 1 + recursiveCall, nil
}

func (expr *SExpr) Length() (*SExpr, error) {
	// eval CAR
	arg1, err := expr.Cdr()
	if err != nil || arg1 == nil || arg1.isNil(){
		return nil, ErrEval
	}
	// get the CAR of cell
	arg1Cell, err := arg1.Car()
	if err != nil || arg1Cell == nil || arg1Cell.isNil(){
		return nil, ErrEval
	}
	// eval arg1
	arg1Eval, _ := arg1Cell.Eval()

	length, err := arg1Eval.lengthHelper()
	if err != nil {
		return nil, ErrEval
	}

	bigLen := big.NewInt(length)
	lengthCell := mkNumber(bigLen)
	return lengthCell, nil
}
 
/*
length:
	expect a sexpr whose:
		car is the symbol CONS
		arglist of size 1:
			we navigate the cdr this is arg 1, then make sure cdr is nil
	now that we have the cell we want to navigate
		if its a func, evaluate it
		count the num of cells needed to r
 
 
– CAR, CDR, CONS and LENGTH;
– Unary predicates ATOM, LISTP and ZEROP;
– Arithmetic operations + and *. To support arbitrary-precision arithmetic for
integers you should use the package big.
*/
 