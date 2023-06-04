package tahwil_test

import (
	"reflect"
	"testing"

	"github.com/go-extras/tahwil"
)

type resolverTest struct {
	in             *tahwil.Value
	out            *tahwil.Value
	hasUnresolved  bool
	unresolvedRefs []uint64
}

func resolverTests() []resolverTest {
	result := make([]resolverTest, 0)

	result = append(result, resolverTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Bool,
				Value: true,
			},
		},
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Bool,
				Value: true,
			},
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	result = append(result, resolverTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Map,
			Value: map[string]interface{}{
				"test": &tahwil.Value{
					Refid: 2,
					Kind:  tahwil.Bool,
					Value: true,
				},
			},
		},
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Map,
			Value: map[string]interface{}{
				"test": &tahwil.Value{
					Refid: 2,
					Kind:  tahwil.Bool,
					Value: true,
				},
			},
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	result = append(result, resolverTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Slice,
			Value: []interface{}{
				&tahwil.Value{
					Refid: 2,
					Kind:  tahwil.Bool,
					Value: true,
				},
			},
		},
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Slice,
			Value: []interface{}{
				&tahwil.Value{
					Refid: 2,
					Kind:  tahwil.Bool,
					Value: true,
				},
			},
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	result = append(result, resolverTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: nil,
		},
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: nil,
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	res1 := &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Ptr,
		Value: &tahwil.Value{
			Refid: 2,
			Kind:  tahwil.Ref,
			Value: true,
		},
	}
	res1.Value.(*tahwil.Value).Value = &tahwil.Reference{
		Refid: 1,
		Value: res1,
	}
	result = append(result, resolverTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Ref,
				Value: uint64(1),
			},
		},
		out:            res1,
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	res2 := &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Ptr,
		Value: &tahwil.Value{
			Refid: 2,
			Kind:  tahwil.Struct,
			Value: map[string]*tahwil.Value{
				"Name": {
					Refid: 3,
					Kind:  tahwil.Ptr,
					Value: &tahwil.Value{
						Refid: 4,
						Kind:  tahwil.String,
						Value: "Mike",
					},
				},
				"Children": {
					Refid: 5,
					Kind:  tahwil.Slice,
					Value: []*tahwil.Value{
						{
							Refid: 6,
							Kind:  tahwil.Ref,
							Value: nil,
						},
					},
				},
			},
		},
	}
	res2.Value.(*tahwil.Value).Value.(map[string]*tahwil.Value)["Children"].Value.([]*tahwil.Value)[0].Value = &tahwil.Reference{
		Refid: 1,
		Value: res2,
	}
	result = append(result, resolverTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"Name": {
						Refid: 3,
						Kind:  tahwil.Ptr,
						Value: &tahwil.Value{
							Refid: 4,
							Kind:  tahwil.String,
							Value: "Mike",
						},
					},
					"Children": {
						Refid: 5,
						Kind:  tahwil.Slice,
						Value: []*tahwil.Value{
							{
								Refid: 6,
								Kind:  tahwil.Ref,
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
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"Name": {
						Refid: 3,
						Kind:  tahwil.Ptr,
						Value: &tahwil.Value{
							Refid: 4,
							Kind:  tahwil.String,
							Value: "Mike",
						},
					},
					"Children": {
						Refid: 5,
						Kind:  tahwil.Slice,
						Value: []*tahwil.Value{
							{
								Refid: 6,
								Kind:  tahwil.Ref,
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

	res3 := &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Ptr,
		Value: &tahwil.Value{
			Refid: 2,
			Kind:  tahwil.Struct,
			Value: map[string]*tahwil.Value{
				"Sibling": {
					Refid: 3,
					Kind:  tahwil.Ref,
					Value: uint64(10),
				},
				"Name": {
					Refid: 4,
					Kind:  tahwil.String,
					Value: "Mike",
				},
				"Parent": {
					Refid: 5,
					Kind:  tahwil.Ptr,
					Value: &tahwil.Value{
						Refid: 6,
						Kind:  tahwil.Struct,
						Value: map[string]*tahwil.Value{
							"Name": {
								Refid: 7,
								Kind:  tahwil.String,
								Value: "Frank",
							},
							"Children": {
								Refid: 8,
								Kind:  tahwil.Slice,
								Value: []*tahwil.Value{
									{
										Refid: 9,
										Kind:  tahwil.Ref,
										Value: uint64(1),
									},
									{
										Refid: 10,
										Kind:  tahwil.Ptr,
										Value: &tahwil.Value{
											Refid: 11,
											Kind:  tahwil.Struct,
											Value: map[string]*tahwil.Value{
												"Name": {
													Refid: 12,
													Kind:  tahwil.String,
													Value: "Zak",
												},
												"Sibling": {
													Refid: 13,
													Kind:  tahwil.Ref,
													Value: uint64(1),
												},
												"Parent": {
													Refid: 14,
													Kind:  tahwil.Ref,
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
	res3.Value.(*tahwil.Value).
		Value.(map[string]*tahwil.Value)["Parent"].
		Value.(*tahwil.Value).
		Value.(map[string]*tahwil.Value)["Children"].
		Value.([]*tahwil.Value)[0].
		Value = &tahwil.Reference{
		Refid: 1,
		Value: res3,
	}
	res3.Value.(*tahwil.Value).
		Value.(map[string]*tahwil.Value)["Parent"].
		Value.(*tahwil.Value).
		Value.(map[string]*tahwil.Value)["Children"].
		Value.([]*tahwil.Value)[1].
		Value.(*tahwil.Value).
		Value.(map[string]*tahwil.Value)["Parent"].Value = &tahwil.Reference{
		Refid: 5,
		Value: res3.Value.(*tahwil.Value).
			Value.(map[string]*tahwil.Value)["Parent"],
	}
	res3.Value.(*tahwil.Value).Value.(map[string]*tahwil.Value)["Parent"].
		Value.(*tahwil.Value).Value.(map[string]*tahwil.Value)["Children"].Value.([]*tahwil.Value)[1].
		Value.(*tahwil.Value).Value.(map[string]*tahwil.Value)["Sibling"].Value = &tahwil.Reference{
		Refid: 1,
		Value: res3,
	}
	res3.Value.(*tahwil.Value).Value.(map[string]*tahwil.Value)["Sibling"].Value = &tahwil.Reference{
		Refid: 10,
		Value: res3.Value.(*tahwil.Value).
			Value.(map[string]*tahwil.Value)["Parent"].
			Value.(*tahwil.Value).
			Value.(map[string]*tahwil.Value)["Children"].
			Value.([]*tahwil.Value)[1],
	}
	result = append(result, resolverTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"Sibling": {
						Refid: 3,
						Kind:  tahwil.Ref,
						Value: uint64(10),
					},
					"Name": {
						Refid: 4,
						Kind:  tahwil.String,
						Value: "Mike",
					},
					"Parent": {
						Refid: 5,
						Kind:  tahwil.Ptr,
						Value: &tahwil.Value{
							Refid: 6,
							Kind:  tahwil.Struct,
							Value: map[string]*tahwil.Value{
								"Name": {
									Refid: 7,
									Kind:  tahwil.String,
									Value: "Frank",
								},
								"Children": {
									Refid: 8,
									Kind:  tahwil.Slice,
									Value: []*tahwil.Value{
										{
											Refid: 9,
											Kind:  tahwil.Ref,
											Value: uint64(1),
										},
										{
											Refid: 10,
											Kind:  tahwil.Ptr,
											Value: &tahwil.Value{
												Refid: 11,
												Kind:  tahwil.Struct,
												Value: map[string]*tahwil.Value{
													"Name": {
														Refid: 12,
														Kind:  tahwil.String,
														Value: "Zak",
													},
													"Sibling": {
														Refid: 13,
														Kind:  tahwil.Ref,
														Value: uint64(1),
													},
													"Parent": {
														Refid: 14,
														Kind:  tahwil.Ref,
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
		r := tahwil.NewResolver(arg.in)
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
