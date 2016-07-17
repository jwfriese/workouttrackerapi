// This file was generated by counterfeiter
package repositoryfakes

import (
	"sync"

	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
	"github.com/jwfriese/workouttrackerapi/workouts/repository"
)

type FakeWorkoutRepository struct {
	AllStub        func() []*workoutdatamodel.Workout
	allMutex       sync.RWMutex
	allArgsForCall []struct{}
	allReturns     struct {
		result1 []*workoutdatamodel.Workout
	}
	GetByIdStub        func(id int) *workoutdatamodel.Workout
	getByIdMutex       sync.RWMutex
	getByIdArgsForCall []struct {
		id int
	}
	getByIdReturns struct {
		result1 *workoutdatamodel.Workout
	}
	InsertStub        func(workout *workoutdatamodel.Workout) (int, error)
	insertMutex       sync.RWMutex
	insertArgsForCall []struct {
		workout *workoutdatamodel.Workout
	}
	insertReturns struct {
		result1 int
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeWorkoutRepository) All() []*workoutdatamodel.Workout {
	fake.allMutex.Lock()
	fake.allArgsForCall = append(fake.allArgsForCall, struct{}{})
	fake.recordInvocation("All", []interface{}{})
	fake.allMutex.Unlock()
	if fake.AllStub != nil {
		return fake.AllStub()
	} else {
		return fake.allReturns.result1
	}
}

func (fake *FakeWorkoutRepository) AllCallCount() int {
	fake.allMutex.RLock()
	defer fake.allMutex.RUnlock()
	return len(fake.allArgsForCall)
}

func (fake *FakeWorkoutRepository) AllReturns(result1 []*workoutdatamodel.Workout) {
	fake.AllStub = nil
	fake.allReturns = struct {
		result1 []*workoutdatamodel.Workout
	}{result1}
}

func (fake *FakeWorkoutRepository) GetById(id int) *workoutdatamodel.Workout {
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

func (fake *FakeWorkoutRepository) GetByIdCallCount() int {
	fake.getByIdMutex.RLock()
	defer fake.getByIdMutex.RUnlock()
	return len(fake.getByIdArgsForCall)
}

func (fake *FakeWorkoutRepository) GetByIdArgsForCall(i int) int {
	fake.getByIdMutex.RLock()
	defer fake.getByIdMutex.RUnlock()
	return fake.getByIdArgsForCall[i].id
}

func (fake *FakeWorkoutRepository) GetByIdReturns(result1 *workoutdatamodel.Workout) {
	fake.GetByIdStub = nil
	fake.getByIdReturns = struct {
		result1 *workoutdatamodel.Workout
	}{result1}
}

func (fake *FakeWorkoutRepository) Insert(workout *workoutdatamodel.Workout) (int, error) {
	fake.insertMutex.Lock()
	fake.insertArgsForCall = append(fake.insertArgsForCall, struct {
		workout *workoutdatamodel.Workout
	}{workout})
	fake.recordInvocation("Insert", []interface{}{workout})
	fake.insertMutex.Unlock()
	if fake.InsertStub != nil {
		return fake.InsertStub(workout)
	} else {
		return fake.insertReturns.result1, fake.insertReturns.result2
	}
}

func (fake *FakeWorkoutRepository) InsertCallCount() int {
	fake.insertMutex.RLock()
	defer fake.insertMutex.RUnlock()
	return len(fake.insertArgsForCall)
}

func (fake *FakeWorkoutRepository) InsertArgsForCall(i int) *workoutdatamodel.Workout {
	fake.insertMutex.RLock()
	defer fake.insertMutex.RUnlock()
	return fake.insertArgsForCall[i].workout
}

func (fake *FakeWorkoutRepository) InsertReturns(result1 int, result2 error) {
	fake.InsertStub = nil
	fake.insertReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeWorkoutRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.allMutex.RLock()
	defer fake.allMutex.RUnlock()
	fake.getByIdMutex.RLock()
	defer fake.getByIdMutex.RUnlock()
	fake.insertMutex.RLock()
	defer fake.insertMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeWorkoutRepository) recordInvocation(key string, args []interface{}) {
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

var _ repository.WorkoutRepository = new(FakeWorkoutRepository)
