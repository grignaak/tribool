package tribool

import "fmt"

func ExampleTribool() {
	var u User // initialized elsewhere

	// The user may not have given their age.
	canVote1 := u.isOlderThan(18).
		OrBool(u.isDev()).
		AndBool(!u.isContestJudge()).
		WithMaybeAsFalse()

	canVote2 :=
		(u.isOlderThan(18).WithMaybeAsFalse() || u.isDev()) &&
			!u.isContestJudge()

	fmt.Println(canVote1, canVote2)
}

type User interface {
	isOlderThan(age int) Tribool
	isContestJudge() bool
	isDev() bool
}

func ExampleTribool_other() {
	var n Node // initialized elsewhere

	// Don't know if a node is activated or not
	ok1 := Maybe.
		And(n.IsActive.ImplyBool(n.ID >= 0)).
		And(n.IsActive.Not().ImplyBool(n.ID == -1)).
		WithMaybeAsTrue()

	ok2 := Maybe.
		Or(n.IsActive.AndBool(n.ID >= 0)).
		Or(n.IsActive.Not().AndBool(n.ID == -1)).
		WithMaybeAsTrue()

	fmt.Println(ok1, ok2)
}

type Node struct {
	IsActive Tribool
	ID       int
}

func ExampleTribool_parsing() {
	fmt.Println(
		FromString("true"),
		FromString("false"),

		FromString("yes"),
		FromString("no"),

		FromString("on"),
		FromString("off"),

		FromString("1"),
		FromString("0"),

		FromString("y"),
		FromString("n"),

		FromString(""),
		FromString("huh?"),
	)

	fmt.Println(
		FromString("true").WithMaybeAsTrue(),
		FromString("false").WithMaybeAsTrue(),
		FromString("").WithMaybeAsTrue(),
	)

	// Output:
	// yes no yes no yes no yes no yes no maybe maybe
	// true false true
}
