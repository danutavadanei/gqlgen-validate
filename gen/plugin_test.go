package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeDependableRule(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "empty", input: "", expect: ""},
		{name: "unchanged-required", input: "required", expect: "required"},
		{name: "unchanged-min", input: "min=10", expect: "min=10"},
		{name: "pipe-mixed", input: "eqfield=confirmPassword|required_with=email", expect: "eqfield=ConfirmPassword|required_with=Email"},
		{name: "comma-mixed", input: "eqfield=confirmPassword,required_with=email", expect: "eqfield=ConfirmPassword,required_with=Email"},
		{name: "eqfield", input: "eqfield=questionId", expect: "eqfield=QuestionID"},
		{name: "eqcsfield", input: "eqcsfield=parent.child", expect: "eqcsfield=Parent.Child"},
		{name: "eqsfield", input: "eqsfield=parent.child", expect: "eqsfield=Parent.Child"},
		{name: "required_with", input: "required_with=email phone", expect: "required_with=Email Phone"},
		{name: "required_if", input: "required_if=otherField foo anotherField bar", expect: "required_if=OtherField foo AnotherField bar"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, toGoRuleParams(tc.input))
		})
	}
}

const pluginSchema = `
    directive @validate(rule: String!, message: String) on INPUT_FIELD_DEFINITION | ARGUMENT_DEFINITION

    input SimpleInput {
        name: String @validate(rule: "min=2,required", message: "name is required")
    }

    input OwnershipInput {
        userOwned: Boolean!
        legalName: String @validate(rule: "required_if=userOwned false")
    }

    input ProxyInput {
        proxy: Boolean @validate(rule: "required")
        target: String @validate(rule: "required_if=proxy true")
    }

    type Mutation {
        registerUser(input: OwnershipInput!): Boolean!
    }
`
