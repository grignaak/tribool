/*
Package tribool implements a tri-state boolean where the extra state is indeterminate.
*/
package tribool

/*
Tribool is a tri-state boolean where the extra state is indeterminate.
*/
type Tribool int

const (
	// No is equivalent to false
	No Tribool = iota

	// Maybe true; maybe false
	Maybe

	// Yes is quivalent to true
	Yes

	// Indeterminate is a synonym for Maybe
	Indeterminate = Maybe
)

var values = [3]Tribool{No, Maybe, Yes}
var strings = [3]string{"No", "Maybe", "Yes"}

/*
FromBool converts a boolean to a Tribool.
*/
func FromBool(b bool) Tribool {
	if b {
		return Yes
	}
	return No
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
	return a != No
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
	return a == Yes
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
Xor implements logical exclusive-or.

			| a.Xor(b)
		a b | b.Xor(a)
		----+----
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
	if s == "Yes" {
		return Yes
	}

	switch len(s) {
	case 1:
		switch s[0] {
		case 't', 'T', 'y', 'Y', '1':
			return Yes
		case 'f', 'F', 'n', 'N', '0':
			return No
		}
	case 2:
		ch0, ch1 := s[0], s[1]
		switch {
		case (ch0 == 'o' || ch0 == 'O') &&
			(ch1 == 'n' || ch1 == 'N'):
			return Yes
		case (ch0 == 'n' || ch0 == 'N') &&
			(ch1 == 'o' || ch1 == 'O'):
			return No
		}
	case 3:
		ch0, ch1, ch2 := s[0], s[1], s[2]
		switch {
		case (ch0 == 'y' || ch0 == 'Y') &&
			(ch1 == 'e' || ch1 == 'E') &&
			(ch2 == 's' || ch2 == 'S'):
			return Yes
		case (ch0 == 'o' || ch0 == 'O') &&
			(ch1 == 'f' || ch1 == 'F') &&
			(ch2 == 'f' || ch2 == 'F'):
			return No
		}
	case 4:
		ch0, ch1, ch2, ch3 := s[0], s[1], s[2], s[3]
		if (ch0 == 't' || ch0 == 'T') &&
			(ch1 == 'r' || ch1 == 'R') &&
			(ch2 == 'u' || ch2 == 'U') &&
			(ch3 == 'e' || ch3 == 'E') {
			return Yes
		}
	case 5:
		ch0, ch1, ch2, ch3, ch4 := s[0], s[1], s[2], s[3], s[4]
		if (ch0 == 'f' || ch0 == 'F') &&
			(ch1 == 'a' || ch1 == 'A') &&
			(ch2 == 'l' || ch2 == 'L') &&
			(ch3 == 's' || ch3 == 'S') &&
			(ch4 == 'e' || ch4 == 'E') {
			return No
		}
	}

	return Maybe
}
