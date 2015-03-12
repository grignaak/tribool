# tribool - a little golang library 
Package tribool implements a tri-state boolean where the extra state is indeterminate.


## Getting started

To get the package, execute:

    go get gopkg.in/grignaak/tribool.v1

To import this package, add the following line to your code:

    import "gopkg.in/grignaak/tribool.v1"

Refer to it as `tribool`.

For more details, see the [API documentation](http://godoc.org/gopkg.in/grignaak/tribool.v1)



## Maybe

A tri-state boolean has the values False, True, and Maybe. Maybe represents a
value that is either true or false but it is inderminate which it is. For
example, you don't know if an http POST was successful if the connection is
dropped after the request was made but before the response came back. This can
be modeled with the indeterminate Maybe value.

Tribool provides tri-state logical operators that act like their boolean
counterparts. The logic tables are documented below.

## Parsing

The tribool package is especially useful for parsing flags that may need a
default value. For example:

	var x string // from somehere
	var flag bool = tribool.FromString(x).WithMaybeAsFalse()
	
Parsing is case sensitive. The following table shows what will be parsed to
true and false values, anything else (including the empty string) results in
the indeterminate value. 

	case insensitive | result
	-----------------+-------
	               t | Yes
	               y | Yes
	               1 | Yes
	              on | Yes
	             yes | Yes
	            true | Yes
	               f | No
	               n | No
	               0 | No
	              no | No
	             off | No
	           false | No
	 <anything else> | Maybe


## Truth Tables

Tribool supports the following binary operations:

		    |  and      or       nand        nor      xor    iff   implies 
		    | a ∧ b   a ∨ b   ¬(a ∧ b)   ¬(a ∨ b)   a ⊕ b  a ⇔ b  a ⇒ b    
		a b | b ∧ a   b ∨ a   ¬(b ∧ a)   ¬(b ∨ a)   b ⊕ a  b ⇔ a    —      
		----+-----------------------------------------------------------------
		N N |   N       N          Y          Y        N       Y      Y    
		N ? |   N       ?          Y          ?        ?       ?      Y    
		N Y |   N       Y          Y          N        Y       N      Y    
		? N |   N       ?          Y          ?        ?       ?      ?    
		? ? |   ?       ?          ?          ?        ?       ?      ?    
		? Y |   ?       Y          ?          N        ?       ?      Y    
		Y N |   N       Y          Y          N        Y       N      N    
		Y ? |   ?       Y          ?          N        ?       ?      ?    
		Y Y |   Y       Y          N          N        N       Y      Y    

Tribool supports the following unary operations:

		a | not   upgrade    downgrade
		--+----------------------------
		N |  Y       N          N
		? |  ?       Y          N
		Y |  N       Y          Y

