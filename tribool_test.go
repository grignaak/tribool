package tribool

import "testing"

func TestTribool_Ops1(t *testing.T) {
	N, x, Y := No, Maybe, Yes
	table := []struct {
		op       string
		a        Tribool
		expected Tribool
	}{
		{"not", N, Y},
		{"not", x, x},
		{"not", Y, N},

		{"upgrade", N, N},
		{"upgrade", x, Y},
		{"upgrade", Y, Y},

		{"downgrade", N, N},
		{"downgrade", x, N},
		{"downgrade", Y, Y},
	}

	op1 := map[string]func(a Tribool) Tribool{
		"not": func(a Tribool) Tribool {
			return a.Not()
		},
		"upgrade": func(a Tribool) Tribool {
			return FromBool(a.WithMaybeAsTrue())
		},
		"downgrade": func(a Tribool) Tribool {
			return FromBool(a.WithMaybeAsFalse())
		},
	}

	for _, test := range table {
		actual := op1[test.op](test.a)
		if actual != test.expected {
			t.Errorf("(%s %s) => %s instead of the expected %s", test.op, test.a, actual, test.expected)
		}
	}
}
func TestTribool_Ops2(t *testing.T) {
	N, x, Y := No, Maybe, Yes
	table := []struct {
		a, b     Tribool
		op       string
		expected Tribool
	}{
		{N, N, "and", N},
		{N, x, "and", N},
		{N, Y, "and", N},
		{x, N, "and", N},
		{x, x, "and", x},
		{x, Y, "and", x},
		{Y, N, "and", N},
		{Y, x, "and", x},
		{Y, Y, "and", Y},

		{N, N, "nand", Y},
		{N, x, "nand", Y},
		{N, Y, "nand", Y},
		{x, N, "nand", Y},
		{x, x, "nand", x},
		{x, Y, "nand", x},
		{Y, N, "nand", Y},
		{Y, x, "nand", x},
		{Y, Y, "nand", N},

		{N, N, "or", N},
		{N, x, "or", x},
		{N, Y, "or", Y},
		{x, N, "or", x},
		{x, x, "or", x},
		{x, Y, "or", Y},
		{Y, N, "or", Y},
		{Y, x, "or", Y},
		{Y, Y, "or", Y},

		{N, N, "nor", Y},
		{N, x, "nor", x},
		{N, Y, "nor", N},
		{x, N, "nor", x},
		{x, x, "nor", x},
		{x, Y, "nor", N},
		{Y, N, "nor", N},
		{Y, x, "nor", N},
		{Y, Y, "nor", N},

		{N, N, "xor", N},
		{N, x, "xor", x},
		{N, Y, "xor", Y},
		{x, N, "xor", x},
		{x, x, "xor", x},
		{x, Y, "xor", x},
		{Y, N, "xor", Y},
		{Y, x, "xor", x},
		{Y, Y, "xor", N},

		{N, N, "equiv", Y},
		{N, x, "equiv", x},
		{N, Y, "equiv", N},
		{x, N, "equiv", x},
		{x, x, "equiv", x},
		{x, Y, "equiv", x},
		{Y, N, "equiv", N},
		{Y, x, "equiv", x},
		{Y, Y, "equiv", Y},

		{N, N, "implies", Y},
		{N, x, "implies", Y},
		{N, Y, "implies", Y},
		{x, N, "implies", x},
		{x, x, "implies", x},
		{x, Y, "implies", Y},
		{Y, N, "implies", N},
		{Y, x, "implies", x},
		{Y, Y, "implies", Y},
	}

	op2 := map[string]func(a, b Tribool) Tribool{
		"and": func(a, b Tribool) Tribool {
			return a.And(b)
		},
		"nand": func(a, b Tribool) Tribool {
			return a.Nand(b)
		},
		"or": func(a, b Tribool) Tribool {
			return a.Or(b)
		},
		"nor": func(a, b Tribool) Tribool {
			return a.Nor(b)
		},
		"xor": func(a, b Tribool) Tribool {
			return a.Xor(b)
		},
		"equiv": func(a, b Tribool) Tribool {
			return a.Equiv(b)
		},
		"implies": func(a, b Tribool) Tribool {
			return a.Imply(b)
		},
	}

	for _, test := range table {
		actual := op2[test.op](test.a, test.b)
		if actual != test.expected {
			t.Errorf("(%s %s %s) => %s instead of the expected %s", test.a, test.op, test.b, actual, test.expected)
		}
	}
}

func TestTribool_parse(t *testing.T) {
	table := []struct {
		raw      string
		expected Tribool
	}{
		{"t", Yes}, {"f", No},
		{"T", Yes}, {"F", No},
		{"y", Yes}, {"n", No},
		{"Y", Yes}, {"N", No},
		{"1", Yes}, {"0", No},

		{"on", Yes}, {"no", No},
		{"ON", Yes}, {"NO", No},
		{"On", Yes}, {"No", No},
		{"oN", Yes}, {"nO", No},

		{"yes", Yes}, {"off", No},
		{"yeS", Yes}, {"ofF", No},
		{"yEs", Yes}, {"oFf", No},
		{"yES", Yes}, {"oFF", No},
		{"Yes", Yes}, {"Off", No},
		{"YeS", Yes}, {"OfF", No},
		{"YEs", Yes}, {"OFf", No},
		{"YES", Yes}, {"OFF", No},

		{"true", Yes}, {"True", Yes},
		{"truE", Yes}, {"TruE", Yes},
		{"trUe", Yes}, {"TrUe", Yes},
		{"trUE", Yes}, {"TrUE", Yes},
		{"tRue", Yes}, {"TRue", Yes},
		{"tRuE", Yes}, {"TRuE", Yes},
		{"tRUe", Yes}, {"TRUe", Yes},
		{"tRUE", Yes}, {"TRUE", Yes},

		{"false", No}, {"False", No},
		{"falsE", No}, {"FalsE", No},
		{"falSe", No}, {"FalSe", No},
		{"falSE", No}, {"FalSE", No},
		{"faLse", No}, {"FaLse", No},
		{"faLsE", No}, {"FaLsE", No},
		{"faLSe", No}, {"FaLSe", No},
		{"faLSE", No}, {"FaLSE", No},
		{"fAlse", No}, {"FAlse", No},
		{"fAlsE", No}, {"FAlsE", No},
		{"fAlSe", No}, {"FAlSe", No},
		{"fAlSE", No}, {"FAlSE", No},
		{"fALse", No}, {"FALse", No},
		{"fALsE", No}, {"FALsE", No},
		{"fALSe", No}, {"FALSe", No},
		{"fALSE", No}, {"FALSE", No},
	}
	errf := "FromString(%s) => %s instead of the expected %s"
	for _, test := range table {
		actual := FromString(test.raw)
		if actual != test.expected {
			t.Errorf(errf, test.raw, actual, test.expected)
		}
		// replace each character with an x, it should result in a maybe
		for i := range test.raw {
			bs := []byte(test.raw)
			bs[i] = 'x'
			actual := FromString(string(bs))
			if actual != Maybe {
				t.Errorf(errf, test.raw, actual, Maybe)
			}
		}
	}

	if FromString("") != Maybe {
		t.Errorf(errf, "", FromString(""), Maybe)
	}
}
