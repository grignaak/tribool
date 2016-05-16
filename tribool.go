/*
Package tribool implements a tri-state boolean where the extra state is indeterminate.

Maybe

A tri-state boolean has the values False, True, and Maybe. Maybe represents a
value that is neither true or false but it is inderminate. For example, you
don't know if an http POST was successful if the connection is dropped after
the request was made but before the response came back. This can be modeled
with the indeterminate Maybe value.

Tribool provides tri-state logical operators that act like their boolean
counterparts. The logic tables are documented below and on each method. There
are also methods for mixed Tribool and bool operations.

Parsing

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


Truth Tables

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

*/
package tribool

import (
	"errors"

	"encoding/json"
)

/*
Tribool is a tri-state boolean where the extra state is indeterminate.

The default value for a Tribool is False, just like a boolean.
*/
type Tribool int

const (
	no Tribool = 0

	// No is equivalent to boolean false
	No Tribool = no

	// False is equivalent to boolean false
	False

	// Off is equivalent to boolean false
	Off
)

const (
	maybe Tribool = 1

	// Maybe represents a value either true or false, but don't know which.
	Maybe Tribool = maybe

	// Perhaps is a synonym for Maybe
	Perhaps

	// Indeterminate is a synonym for Maybe
	Indeterminate
)

const (
	yes Tribool = 2

	// Yes is equivalent to boolean true
	Yes Tribool = yes

	// True is equivalent to boolean true
	True

	// On is equivalent to boolean true
	On
)

var values = [3]Tribool{No, Maybe, Yes}
var strings = [3]string{"no", "maybe", "yes"}

/*
FromBool converts a bool to an equivalent Tribool.
*/
func FromBool(b bool) Tribool {
	if b {
		return yes
	}
	return no
}

/*
String converts a Tribool to a string that can be parsed with FromString
*/
func (a Tribool) String() string {
	return strings[a]
}

/*
WithMaybeAsTrue converts the Tribool to a boolean by coercing Maybe to true.

		a | a.WithMaybeAsTrue()
		--+----------------------
		N | N
		? | Y
		Y | Y
*/
func (a Tribool) WithMaybeAsTrue() bool {
	return a != no
}

/*
WithMaybeAsFalse converts the Tribool to a boolean by coercing Maybe to false.

		a | a.WithMaybeAsTrue()
		--+----------------------
		N | N
		? | N
		Y | Y
*/
func (a Tribool) WithMaybeAsFalse() bool {
	return a == yes
}

/*
And implements logical and.

		    | a.And(b)
		a b | b.And(a)
		----+---------
		N N | N
		N ? | N
		N Y | N
		? N | N
		? ? | ?
		? Y | ?
		Y N | N
		Y ? | ?
		Y Y | Y
*/
func (a Tribool) And(b Tribool) Tribool {
	return values[min(a, b)]
}

func min(a, b Tribool) Tribool {
	if a < b {
		return a
	}
	return b
}

/*
AndBool is equivalent to a.And(FromBool(b))
*/
func (a Tribool) AndBool(b bool) Tribool {
	return a.And(FromBool(b))
}

/*
Or implements logical inclusive-or.

		    | a.Or(b)
		a b | b.Or(a)
		----+----------
		N N | N
		N ? | ?
		N Y | Y
		? N | ?
		? ? | ?
		? Y | Y
		Y N | Y
		Y ? | Y
		Y Y | Y
*/
func (a Tribool) Or(b Tribool) Tribool {
	return values[max(a, b)]
}

func max(a, b Tribool) Tribool {
	if a > b {
		return a
	}
	return b
}

/*
OrBool is equivalent to a.Or(FromBool(b))
*/
func (a Tribool) OrBool(b bool) Tribool {
	return a.Or(FromBool(b))
}

/*
Nand implements logical nand.

		    | a.Nand(b)
		a b | b.Nand(a)
		----+----------
		N N | Y
		N ? | Y
		N Y | Y
		? N | Y
		? ? | ?
		? Y | ?
		Y N | Y
		Y ? | ?
		Y Y | N
*/
func (a Tribool) Nand(b Tribool) Tribool {
	return values[2-min(a, b)]
}

/*
NandBool is equivalent to a.Nand(FromBool(b))
*/
func (a Tribool) NandBool(b bool) Tribool {
	return a.Nand(FromBool(b))
}

/*
Not implements logical not.

		 a | a.Not()
		 --+--------
		 N | Y
		 ? | ?
		 Y | N
*/
func (a Tribool) Not() Tribool {
	return values[2-a]
}

/*
Nor implements logical nor.

		    | a.Nor(b)
		a b | b.Nor(a)
		----+------------
		N N | Y
		N ? | ?
		N Y | N
		? N | ?
		? ? | ?
		? Y | N
		Y N | N
		Y ? | N
		Y Y | N
*/
func (a Tribool) Nor(b Tribool) Tribool {
	return values[2-max(a, b)]
}

/*
NorBool is equivalent to a.Nor(FromBool(b))
*/
func (a Tribool) NorBool(b bool) Tribool {
	return a.Nor(FromBool(b))
}

/*
Xor implements logical exclusive-or.

		    | a.Xor(b)
		a b | b.Xor(a)
		----+---------
		N N | N
		N ? | ?
		N Y | Y
		? N | ?
		? ? | ?
		? Y | ?
		Y N | Y
		Y ? | ?
		Y Y | N
*/
func (a Tribool) Xor(b Tribool) Tribool {
	return a.Or(b).And(a.Nand(b))
}

/*
XorBool is equivalent to a.Xor(FromBool(b))
*/
func (a Tribool) XorBool(b bool) Tribool {
	return a.Xor(FromBool(b))
}

/*
Imply implements logical implication.

Implication is not reflexive. That is a.Imply(b) is not the same as b.Imply(a)

		a b | a.Imply(b)
		----+------------
		N N | Y
		N ? | Y
		N Y | Y
		? N | ?
		? ? | ?
		? Y | Y
		Y N | N
		Y ? | ?
		Y Y | Y
*/
func (a Tribool) Imply(b Tribool) Tribool {
	return b.Or(a.Not())
}

/*
ImplyBool is equivalent to a.Imply(FromBool(b))
*/
func (a Tribool) ImplyBool(b bool) Tribool {
	return a.Imply(FromBool(b))
}

/*
Equiv implements logical equivalence.

		    | a.Equiv(b)
		a b | b.Equiv(a)
		----+------
		N N | Y
		N ? | ?
		N Y | N
		? N | ?
		? ? | ?
		? Y | ?
		Y N | N
		Y ? | ?
		Y Y | Y
*/
func (a Tribool) Equiv(b Tribool) Tribool {
	return a.And(b).Or(a.Nor(b))
}

/*
EquivBool is equivalent to a.Equiv(FromBool(b))
*/
func (a Tribool) EquivBool(b bool) Tribool {
	return a.Equiv(FromBool(b))
}

/*
FromString converts a string to a Tribool.

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
*/
func FromString(s string) Tribool {
	// most flags will be marked as true. This is the fast-path.
	if s == "true" {
		return yes
	}

	switch len(s) {
	case 1:
		switch s[0] {
		case 't', 'T', 'y', 'Y', '1':
			return yes
		case 'f', 'F', 'n', 'N', '0':
			return no
		}
	case 2:
		ch0, ch1 := s[0], s[1]
		switch {
		case (ch0 == 'o' || ch0 == 'O') &&
			(ch1 == 'n' || ch1 == 'N'):
			return yes
		case (ch0 == 'n' || ch0 == 'N') &&
			(ch1 == 'o' || ch1 == 'O'):
			return no
		}
	case 3:
		ch0, ch1, ch2 := s[0], s[1], s[2]
		switch {
		case (ch0 == 'y' || ch0 == 'Y') &&
			(ch1 == 'e' || ch1 == 'E') &&
			(ch2 == 's' || ch2 == 'S'):
			return yes
		case (ch0 == 'o' || ch0 == 'O') &&
			(ch1 == 'f' || ch1 == 'F') &&
			(ch2 == 'f' || ch2 == 'F'):
			return no
		}
	case 4:
		ch0, ch1, ch2, ch3 := s[0], s[1], s[2], s[3]
		if (ch0 == 't' || ch0 == 'T') &&
			(ch1 == 'r' || ch1 == 'R') &&
			(ch2 == 'u' || ch2 == 'U') &&
			(ch3 == 'e' || ch3 == 'E') {
			return yes
		}
	case 5:
		ch0, ch1, ch2, ch3, ch4 := s[0], s[1], s[2], s[3], s[4]
		if (ch0 == 'f' || ch0 == 'F') &&
			(ch1 == 'a' || ch1 == 'A') &&
			(ch2 == 'l' || ch2 == 'L') &&
			(ch3 == 's' || ch3 == 'S') &&
			(ch4 == 'e' || ch4 == 'E') {
			return no
		}
	}

	return maybe
}

// MarshalJSON marshals tribools to strings, using the Tribool.String() method.
func (a Tribool) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON supports unmarshalling from a json string (using `FromString()`), a json
// boolean (using `FromBool()`), and treats anything else as `maybe`.
func (a *Tribool) UnmarshalJSON(data []byte) error {
	if a == nil {
		return errors.New("tribool.TriBool: UnmarshalJSON on nil pointer")
	}
	var s string
	var b bool
	if err := json.Unmarshal(data, &s); err == nil {
		*a = FromString(s)
	} else if err := json.Unmarshal(data, &b); err == nil {
		*a = FromBool(b)
	} else {
		*a = Maybe
	}
	return nil
}
