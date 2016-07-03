// This file was generated by counterfeiter
package repositoryfakes

import (
	"sync"

	"github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/sets/repository"
)

type FakeSetRepository struct {
	GetByIdStub        func(id int) *datamodel.Set
	getByIdMutex       sync.RWMutex
	getByIdArgsForCall []struct {
		id int
	}
	getByIdReturns struct {
		result1 *datamodel.Set
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSetRepository) GetById(id int) *datamodel.Set {
	fake.getByIdMutex.Lock()
	fake.getByIdArgsForCall = append(fake.getByIdArgsForCall, struct {
		id int
	}{id})
	fake.recordInvocation("GetById", []interface{}{id})
	fake.getByIdMutex.Unlock()
	if fake.GetByIdStub != nil {
		return fake.GetByIdStub(id)
	} else {
		return fake.getByIdReturns.result1
	}
}

func (fake *FakeSetRepository) GetByIdCallCount() int {
	fake.getByIdMutex.RLock()
	defer fake.getByIdMutex.RUnlock()
	return len(fake.getByIdArgsForCall)
}

func (fake *FakeSetRepository) GetByIdArgsForCall(i int) int {
	fake.getByIdMutex.RLock()
	defer fake.getByIdMutex.RUnlock()
	return fake.getByIdArgsForCall[i].id
}

func (fake *FakeSetRepository) GetByIdReturns(result1 *datamodel.Set) {
	fake.GetByIdStub = nil
	fake.getByIdReturns = struct {
		result1 *datamodel.Set
	}{result1}
}

func (fake *FakeSetRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getByIdMutex.RLock()
	defer fake.getByIdMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeSetRepository) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ repository.SetRepository = new(FakeSetRepository)
