package symbolextraction

import (
	"fmt"
	"slices"
)

func (f *FunctionSymbol) ToUnlinkedCallNode() CallNode {
	return CallNode{
		Name:                f.Name,
		MaxStackSizeCallees: f.MaxStackSizeCallees,
		StackSize:           f.StackSize,
		Address:             f.Address,
		// Calls: - no recursion :)
	}
}

func (f *FunctionSymbol) ToUnlinkedCallRoot() []CallNode {
	return []CallNode{f.ToUnlinkedCallNode()}
}

func (f *FunctionSymbol) ToUnlinkedCallBranchRoot() CallBranch {
	return CallBranch{
		CallList:  f.ToUnlinkedCallRoot(),
		StackSize: f.StackSize,
	}
}

func (t *CallTree) PushCallerToCallBranch(f *FunctionSymbol) {
	node := f.ToUnlinkedCallNode()
	last := t.CurrentBranch.CallList[len(t.CurrentBranch.CallList)-1]
	node.Caller = &last
	t.CurrentBranch.CallList = append(t.CurrentBranch.CallList, node)
	t.CurrentBranch.StackSize += f.StackSize
	// for _, v := range t.CurrentBranch.CallList {
	// 	fmt.Print(fmt.Printf("%X ", v.Address))
	// }
	// fmt.Println("\n\tl push", len(t.CurrentBranch.CallList), f.Name)
}

func (t *CallTree) PopCallerToCallBranch() {
	l := len(t.CurrentBranch.CallList) - 1
	f := t.CurrentBranch.CallList[l]
	// fmt.Println("l pop", len(t.CurrentBranch.CallList), f.Name)
	t.CurrentBranch.CallList = t.CurrentBranch.CallList[:l]
	t.CurrentBranch.StackSize -= f.StackSize
}

func (t *CallTree) hasIndirectRecursion(callAddr uint64) bool {
	return slices.ContainsFunc(t.CurrentBranch.CallList, func(n CallNode) bool {
		return n.Address == callAddr
	})
}

func (f *FunctionSymbol) traverseCallSubGraph(parent *CallNode, t *CallTree, s *SElfReport) {

	// if f.Visited {
	// 	// fmt.Println(strings.Repeat("\t", int(calldepth)), f.Name, "revisited stacksize:", f.StackSize, "biggest calletreesize:", f.MaxStackSizeCallees)
	// 	return f.MaxStackSizeCallees
	// }
	f.Visited = true
	if len(f.Callees) == 0 {
		call_branch := CallBranch{}
		for _, c := range t.CurrentBranch.CallList {
			cn := CallNode{Name: c.Name, StackSize: c.StackSize, Recursion: c.Recursion}
			call_branch.CallList = append(call_branch.CallList, cn)
		}
		call_branch.StackSize = t.CurrentBranch.StackSize + f.StackSize
		// fmt.Println("fincalstr", len(call_branch.CallList), call_branch.StackSize)
		t.Branches = append(t.Branches, call_branch)
	} else {
		for i := 0; i < len(f.Callees); i++ {
			call := &f.Callees[i]
			callAddr := call.CallTo
			if callAddr == nil {
				t.UnresolvedCalls = append(t.UnresolvedCalls, *call)
				continue
			}
			fnCallee, ok := s.Addr2FnMap[*callAddr]
			if !ok {
				t.UnresolvedCalls = append(t.UnresolvedCalls, *call)
				continue
			}
			directRecursion := fnCallee.Address == f.Address
			if directRecursion {
				// TODO log user
				rec := &t.CurrentBranch.CallList[len(t.CurrentBranch.CallList)-1]
				rec.Recursion = true
				msg := fmt.Sprintf("warning direct recursion in '%s' not supported yet, TODO needs a recursion limit", f.Name)
				s.Info = append(s.Info, msg)
			} else if t.hasIndirectRecursion(fnCallee.Address) {
				found := false
				for i, _ := range t.CurrentBranch.CallList {
					c := &t.CurrentBranch.CallList[i]
					if found || c.Address == fnCallee.Address {
						c.Recursion = true
						found = true
					}
				}
				msg := fmt.Sprintf("warning indirect recursion in '%s' not supported yet, TODO needs a recursion limit", f.Name)
				s.Info = append(s.Info, msg)
			} else {
				n := fnCallee.ToUnlinkedCallNode()
				n.Calls = make([]CallNode, 0)

				t.PushCallerToCallBranch(fnCallee)
				fnCallee.traverseCallSubGraph(&n, t, s)
				parent.Calls = append(parent.Calls, n)
				t.PopCallerToCallBranch()
			}
		}
	}
}

func (f *FunctionSymbol) GetCallTreeJson(s *SElfReport) *CallTree {
	root := f.ToUnlinkedCallNode()
	var t *CallTree = &CallTree{
		Tree:            root,
		Branches:        []CallBranch{},
		UnresolvedCalls: []FunctionCall{},
		CurrentBranch:   CallBranch{CallList: f.ToUnlinkedCallRoot(), StackSize: f.StackSize},
	}
	// reset visited to avoid recursion...
	for i := 0; i < len(s.Functions); i++ {
		function := &s.Functions[i]
		function.Visited = false
	}
	f.traverseCallSubGraph(&t.Tree, t, s)
	// f.traverseCallSubGraph(t, s)
	// sort biggest callpath first
	slices.SortFunc(t.Branches, func(i, j CallBranch) int {
		return int(j.StackSize) - int(i.StackSize)
	})
	return t
}
