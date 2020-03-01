package tahwil

import (
	"reflect"
	"testing"
)

type resolverTest struct {
	in             *Value
	out            *Value
	hasUnresolved  bool
	unresolvedRefs []uint64
}

func resolverTests() []resolverTest {
	result := make([]resolverTest, 0)

	result = append(result, resolverTest{
		in: &Value{
			Refid: 1,
			Kind:  "ptr",
			Value: &Value{
				Refid: 2,
				Kind:  "bool",
				Value: true,
			},
		},
		out: &Value{
			Refid: 1,
			Kind:  "ptr",
			Value: &Value{
				Refid: 2,
				Kind:  "bool",
				Value: true,
			},
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	result = append(result, resolverTest{
		in: &Value{
			Refid: 1,
			Kind:  "map",
			Value: map[string]interface{}{
				"test": &Value{
					Refid: 2,
					Kind:  "bool",
					Value: true,
				},
			},
		},
		out: &Value{
			Refid: 1,
			Kind:  "map",
			Value: map[string]interface{}{
				"test": &Value{
					Refid: 2,
					Kind:  "bool",
					Value: true,
				},
			},
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	result = append(result, resolverTest{
		in: &Value{
			Refid: 1,
			Kind:  "slice",
			Value: []interface{}{
				&Value{
					Refid: 2,
					Kind:  "bool",
					Value: true,
				},
			},
		},
		out: &Value{
			Refid: 1,
			Kind:  "slice",
			Value: []interface{}{
				&Value{
					Refid: 2,
					Kind:  "bool",
					Value: true,
				},
			},
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	result = append(result, resolverTest{
		in: &Value{
			Refid: 1,
			Kind:  "ptr",
			Value: nil,
		},
		out: &Value{
			Refid: 1,
			Kind:  "ptr",
			Value: nil,
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	res1 := &Value{
		Refid: 1,
		Kind:  "ptr",
		Value: &Value{
			Refid: 2,
			Kind:  "ref",
			Value: true,
		},
	}
	res1.Value.(*Value).Value = &Reference{
		Refid: 1,
		Value: res1,
	}
	result = append(result, resolverTest{
		in: &Value{
			Refid: 1,
			Kind:  "ptr",
			Value: &Value{
				Refid: 2,
				Kind:  "ref",
				Value: uint64(1),
			},
		},
		out:            res1,
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	res2 := &Value{
		Refid: 1,
		Kind:  "ptr",
		Value: &Value{
			Refid: 2,
			Kind:  "struct",
			Value: map[string]*Value{
				"Name": {
					Refid: 3,
					Kind:  "ptr",
					Value: &Value{
						Refid: 4,
						Kind:  "string",
						Value: "Mike",
					},
				},
				"Children": {
					Refid: 5,
					Kind:  "slice",
					Value: []*Value{
						{
							Refid: 6,
							Kind:  "ref",
							Value: nil,
						},
					},
				},
			},
		},
	}
	res2.Value.(*Value).Value.(map[string]*Value)["Children"].Value.([]*Value)[0].Value = &Reference{
		Refid: 1,
		Value: res2,
	}
	result = append(result, resolverTest{
		in: &Value{
			Refid: 1,
			Kind:  "ptr",
			Value: &Value{
				Refid: 2,
				Kind:  "struct",
				Value: map[string]*Value{
					"Name": {
						Refid: 3,
						Kind:  "ptr",
						Value: &Value{
							Refid: 4,
							Kind:  "string",
							Value: "Mike",
						},
					},
					"Children": {
						Refid: 5,
						Kind:  "slice",
						Value: []*Value{
							{
								Refid: 6,
								Kind:  "ref",
								Value: uint64(1),
							},
						},
					},
				},
			},
		},
		out:            res2,
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	result = append(result, resolverTest{
		in: &Value{
			Refid: 1,
			Kind:  "ptr",
			Value: &Value{
				Refid: 2,
				Kind:  "struct",
				Value: map[string]*Value{
					"Name": {
						Refid: 3,
						Kind:  "ptr",
						Value: &Value{
							Refid: 4,
							Kind:  "string",
							Value: "Mike",
						},
					},
					"Children": {
						Refid: 5,
						Kind:  "slice",
						Value: []*Value{
							{
								Refid: 6,
								Kind:  "ref",
								Value: uint64(9),
							},
						},
					},
				},
			},
		},
		hasUnresolved:  true,
		unresolvedRefs: []uint64{9},
	})

	res3 := &Value{
		Refid: 1,
		Kind:  "ptr",
		Value: &Value{
			Refid: 2,
			Kind:  "struct",
			Value: map[string]*Value{
				"Sibling": {
					Refid: 3,
					Kind:  "ref",
					Value: uint64(10),
				},
				"Name": {
					Refid: 4,
					Kind:  "string",
					Value: "Mike",
				},
				"Parent": {
					Refid: 5,
					Kind:  "ptr",
					Value: &Value{
						Refid: 6,
						Kind:  "struct",
						Value: map[string]*Value{
							"Name": {
								Refid: 7,
								Kind:  "string",
								Value: "Frank",
							},
							"Children": {
								Refid: 8,
								Kind:  "slice",
								Value: []*Value{
									{
										Refid: 9,
										Kind:  "ref",
										Value: uint64(1),
									},
									{
										Refid: 10,
										Kind:  "ptr",
										Value: &Value{
											Refid: 11,
											Kind:  "struct",
											Value: map[string]*Value{
												"Name": {
													Refid: 12,
													Kind:  "string",
													Value: "Zak",
												},
												"Sibling": {
													Refid: 13,
													Kind:  "ref",
													Value: uint64(1),
												},
												"Parent": {
													Refid: 14,
													Kind:  "ref",
													Value: uint64(5),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	res3.Value.(*Value).
		Value.(map[string]*Value)["Parent"].
		Value.(*Value).
		Value.(map[string]*Value)["Children"].
		Value.([]*Value)[0].
		Value = &Reference{
		Refid: 1,
		Value: res3,
	}
	res3.Value.(*Value).
		Value.(map[string]*Value)["Parent"].
		Value.(*Value).
		Value.(map[string]*Value)["Children"].
		Value.([]*Value)[1].
		Value.(*Value).
		Value.(map[string]*Value)["Parent"].Value = &Reference{
		Refid: 5,
		Value: res3.Value.(*Value).
			Value.(map[string]*Value)["Parent"],
	}
	res3.Value.(*Value).Value.(map[string]*Value)["Parent"].
		Value.(*Value).Value.(map[string]*Value)["Children"].Value.([]*Value)[1].
		Value.(*Value).Value.(map[string]*Value)["Sibling"].Value = &Reference{
		Refid: 1,
		Value: res3,
	}
	res3.Value.(*Value).Value.(map[string]*Value)["Sibling"].Value = &Reference{
		Refid: 10,
		Value: res3.Value.(*Value).
			Value.(map[string]*Value)["Parent"].
			Value.(*Value).
			Value.(map[string]*Value)["Children"].
			Value.([]*Value)[1],
	}
	result = append(result, resolverTest{
		in: &Value{
			Refid: 1,
			Kind:  "ptr",
			Value: &Value{
				Refid: 2,
				Kind:  "struct",
				Value: map[string]*Value{
					"Sibling": {
						Refid: 3,
						Kind:  "ref",
						Value: uint64(10),
					},
					"Name": {
						Refid: 4,
						Kind:  "string",
						Value: "Mike",
					},
					"Parent": {
						Refid: 5,
						Kind:  "ptr",
						Value: &Value{
							Refid: 6,
							Kind:  "struct",
							Value: map[string]*Value{
								"Name": {
									Refid: 7,
									Kind:  "string",
									Value: "Frank",
								},
								"Children": {
									Refid: 8,
									Kind:  "slice",
									Value: []*Value{
										{
											Refid: 9,
											Kind:  "ref",
											Value: uint64(1),
										},
										{
											Refid: 10,
											Kind:  "ptr",
											Value: &Value{
												Refid: 11,
												Kind:  "struct",
												Value: map[string]*Value{
													"Name": {
														Refid: 12,
														Kind:  "string",
														Value: "Zak",
													},
													"Sibling": {
														Refid: 13,
														Kind:  "ref",
														Value: uint64(1),
													},
													"Parent": {
														Refid: 14,
														Kind:  "ref",
														Value: uint64(5),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		out:            res3,
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	return result
}

func TestResolver_Resolve(t *testing.T) {
	for i, arg := range resolverTests() {
		r := NewResolver(arg.in)
		err := r.Resolve()
		if err != nil {
			t.Errorf("#%d: Resolver.Resolve() returned an error: %s", i, err.Error())
			continue
		}
		hasUnresolved := r.HasUnresolved()
		unresolvedRefs := r.Unresolved()
		if hasUnresolved != arg.hasUnresolved {
			t.Errorf("#%d: Resolver.HasUnresolved mismatch\nhave: %v\nwant: %v", i, hasUnresolved, arg.hasUnresolved)
		}
		if !reflect.DeepEqual(unresolvedRefs, arg.unresolvedRefs) {
			t.Errorf("#%d: Resolver.Unresolved mismatch\nhave: %#+v\nwant: %#+v", i, unresolvedRefs, arg.unresolvedRefs)
		}
		// don't check if arg.out is nil
		if arg.out != nil && !reflect.DeepEqual(arg.in, arg.out) {
			t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, arg.in, arg.out)
		}
	}
}
