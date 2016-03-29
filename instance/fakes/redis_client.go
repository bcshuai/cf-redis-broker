// This file was generated by counterfeiter
package fakes

import (
	"sync"
	"time"

	"github.com/bcshuai/cf-redis-broker/redis/client"
)

type FakeRedisClient struct {
	DisconnectStub        func() error
	disconnectMutex       sync.RWMutex
	disconnectArgsForCall []struct{}
	disconnectReturns     struct {
		result1 error
	}
	WaitUntilRedisNotLoadingStub        func(timeoutMilliseconds int) error
	waitUntilRedisNotLoadingMutex       sync.RWMutex
	waitUntilRedisNotLoadingArgsForCall []struct {
		timeoutMilliseconds int
	}
	waitUntilRedisNotLoadingReturns struct {
		result1 error
	}
	EnableAOFStub        func() error
	enableAOFMutex       sync.RWMutex
	enableAOFArgsForCall []struct{}
	enableAOFReturns     struct {
		result1 error
	}
	LastRDBSaveTimeStub        func() (int64, error)
	lastRDBSaveTimeMutex       sync.RWMutex
	lastRDBSaveTimeArgsForCall []struct{}
	lastRDBSaveTimeReturns     struct {
		result1 int64
		result2 error
	}
	InfoStub        func() (map[string]string, error)
	infoMutex       sync.RWMutex
	infoArgsForCall []struct{}
	infoReturns     struct {
		result1 map[string]string
		result2 error
	}
	InfoFieldStub        func(fieldName string) (string, error)
	infoFieldMutex       sync.RWMutex
	infoFieldArgsForCall []struct {
		fieldName string
	}
	infoFieldReturns struct {
		result1 string
		result2 error
	}
	GetConfigStub        func(key string) (string, error)
	getConfigMutex       sync.RWMutex
	getConfigArgsForCall []struct {
		key string
	}
	getConfigReturns struct {
		result1 string
		result2 error
	}
	RDBPathStub        func() (string, error)
	rDBPathMutex       sync.RWMutex
	rDBPathArgsForCall []struct{}
	rDBPathReturns     struct {
		result1 string
		result2 error
	}
	AddressStub        func() string
	addressMutex       sync.RWMutex
	addressArgsForCall []struct{}
	addressReturns     struct {
		result1 string
	}
	WaitForNewSaveSinceStub        func(lastSaveTime int64, timeout time.Duration) error
	waitForNewSaveSinceMutex       sync.RWMutex
	waitForNewSaveSinceArgsForCall []struct {
		lastSaveTime int64
		timeout      time.Duration
	}
	waitForNewSaveSinceReturns struct {
		result1 error
	}
	RunBGSaveStub        func() error
	runBGSaveMutex       sync.RWMutex
	runBGSaveArgsForCall []struct{}
	runBGSaveReturns     struct {
		result1 error
	}
}

func (fake *FakeRedisClient) Disconnect() error {
	fake.disconnectMutex.Lock()
	fake.disconnectArgsForCall = append(fake.disconnectArgsForCall, struct{}{})
	fake.disconnectMutex.Unlock()
	if fake.DisconnectStub != nil {
		return fake.DisconnectStub()
	} else {
		return fake.disconnectReturns.result1
	}
}

func (fake *FakeRedisClient) DisconnectCallCount() int {
	fake.disconnectMutex.RLock()
	defer fake.disconnectMutex.RUnlock()
	return len(fake.disconnectArgsForCall)
}

func (fake *FakeRedisClient) DisconnectReturns(result1 error) {
	fake.DisconnectStub = nil
	fake.disconnectReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRedisClient) WaitUntilRedisNotLoading(timeoutMilliseconds int) error {
	fake.waitUntilRedisNotLoadingMutex.Lock()
	fake.waitUntilRedisNotLoadingArgsForCall = append(fake.waitUntilRedisNotLoadingArgsForCall, struct {
		timeoutMilliseconds int
	}{timeoutMilliseconds})
	fake.waitUntilRedisNotLoadingMutex.Unlock()
	if fake.WaitUntilRedisNotLoadingStub != nil {
		return fake.WaitUntilRedisNotLoadingStub(timeoutMilliseconds)
	} else {
		return fake.waitUntilRedisNotLoadingReturns.result1
	}
}

func (fake *FakeRedisClient) WaitUntilRedisNotLoadingCallCount() int {
	fake.waitUntilRedisNotLoadingMutex.RLock()
	defer fake.waitUntilRedisNotLoadingMutex.RUnlock()
	return len(fake.waitUntilRedisNotLoadingArgsForCall)
}

func (fake *FakeRedisClient) WaitUntilRedisNotLoadingArgsForCall(i int) int {
	fake.waitUntilRedisNotLoadingMutex.RLock()
	defer fake.waitUntilRedisNotLoadingMutex.RUnlock()
	return fake.waitUntilRedisNotLoadingArgsForCall[i].timeoutMilliseconds
}

func (fake *FakeRedisClient) WaitUntilRedisNotLoadingReturns(result1 error) {
	fake.WaitUntilRedisNotLoadingStub = nil
	fake.waitUntilRedisNotLoadingReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRedisClient) EnableAOF() error {
	fake.enableAOFMutex.Lock()
	fake.enableAOFArgsForCall = append(fake.enableAOFArgsForCall, struct{}{})
	fake.enableAOFMutex.Unlock()
	if fake.EnableAOFStub != nil {
		return fake.EnableAOFStub()
	} else {
		return fake.enableAOFReturns.result1
	}
}

func (fake *FakeRedisClient) EnableAOFCallCount() int {
	fake.enableAOFMutex.RLock()
	defer fake.enableAOFMutex.RUnlock()
	return len(fake.enableAOFArgsForCall)
}

func (fake *FakeRedisClient) EnableAOFReturns(result1 error) {
	fake.EnableAOFStub = nil
	fake.enableAOFReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRedisClient) LastRDBSaveTime() (int64, error) {
	fake.lastRDBSaveTimeMutex.Lock()
	fake.lastRDBSaveTimeArgsForCall = append(fake.lastRDBSaveTimeArgsForCall, struct{}{})
	fake.lastRDBSaveTimeMutex.Unlock()
	if fake.LastRDBSaveTimeStub != nil {
		return fake.LastRDBSaveTimeStub()
	} else {
		return fake.lastRDBSaveTimeReturns.result1, fake.lastRDBSaveTimeReturns.result2
	}
}

func (fake *FakeRedisClient) LastRDBSaveTimeCallCount() int {
	fake.lastRDBSaveTimeMutex.RLock()
	defer fake.lastRDBSaveTimeMutex.RUnlock()
	return len(fake.lastRDBSaveTimeArgsForCall)
}

func (fake *FakeRedisClient) LastRDBSaveTimeReturns(result1 int64, result2 error) {
	fake.LastRDBSaveTimeStub = nil
	fake.lastRDBSaveTimeReturns = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *FakeRedisClient) Info() (map[string]string, error) {
	fake.infoMutex.Lock()
	fake.infoArgsForCall = append(fake.infoArgsForCall, struct{}{})
	fake.infoMutex.Unlock()
	if fake.InfoStub != nil {
		return fake.InfoStub()
	} else {
		return fake.infoReturns.result1, fake.infoReturns.result2
	}
}

func (fake *FakeRedisClient) InfoCallCount() int {
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	return len(fake.infoArgsForCall)
}

func (fake *FakeRedisClient) InfoReturns(result1 map[string]string, result2 error) {
	fake.InfoStub = nil
	fake.infoReturns = struct {
		result1 map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeRedisClient) InfoField(fieldName string) (string, error) {
	fake.infoFieldMutex.Lock()
	fake.infoFieldArgsForCall = append(fake.infoFieldArgsForCall, struct {
		fieldName string
	}{fieldName})
	fake.infoFieldMutex.Unlock()
	if fake.InfoFieldStub != nil {
		return fake.InfoFieldStub(fieldName)
	} else {
		return fake.infoFieldReturns.result1, fake.infoFieldReturns.result2
	}
}

func (fake *FakeRedisClient) InfoFieldCallCount() int {
	fake.infoFieldMutex.RLock()
	defer fake.infoFieldMutex.RUnlock()
	return len(fake.infoFieldArgsForCall)
}

func (fake *FakeRedisClient) InfoFieldArgsForCall(i int) string {
	fake.infoFieldMutex.RLock()
	defer fake.infoFieldMutex.RUnlock()
	return fake.infoFieldArgsForCall[i].fieldName
}

func (fake *FakeRedisClient) InfoFieldReturns(result1 string, result2 error) {
	fake.InfoFieldStub = nil
	fake.infoFieldReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRedisClient) GetConfig(key string) (string, error) {
	fake.getConfigMutex.Lock()
	fake.getConfigArgsForCall = append(fake.getConfigArgsForCall, struct {
		key string
	}{key})
	fake.getConfigMutex.Unlock()
	if fake.GetConfigStub != nil {
		return fake.GetConfigStub(key)
	} else {
		return fake.getConfigReturns.result1, fake.getConfigReturns.result2
	}
}

func (fake *FakeRedisClient) GetConfigCallCount() int {
	fake.getConfigMutex.RLock()
	defer fake.getConfigMutex.RUnlock()
	return len(fake.getConfigArgsForCall)
}

func (fake *FakeRedisClient) GetConfigArgsForCall(i int) string {
	fake.getConfigMutex.RLock()
	defer fake.getConfigMutex.RUnlock()
	return fake.getConfigArgsForCall[i].key
}

func (fake *FakeRedisClient) GetConfigReturns(result1 string, result2 error) {
	fake.GetConfigStub = nil
	fake.getConfigReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRedisClient) RDBPath() (string, error) {
	fake.rDBPathMutex.Lock()
	fake.rDBPathArgsForCall = append(fake.rDBPathArgsForCall, struct{}{})
	fake.rDBPathMutex.Unlock()
	if fake.RDBPathStub != nil {
		return fake.RDBPathStub()
	} else {
		return fake.rDBPathReturns.result1, fake.rDBPathReturns.result2
	}
}

func (fake *FakeRedisClient) RDBPathCallCount() int {
	fake.rDBPathMutex.RLock()
	defer fake.rDBPathMutex.RUnlock()
	return len(fake.rDBPathArgsForCall)
}

func (fake *FakeRedisClient) RDBPathReturns(result1 string, result2 error) {
	fake.RDBPathStub = nil
	fake.rDBPathReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRedisClient) Address() string {
	fake.addressMutex.Lock()
	fake.addressArgsForCall = append(fake.addressArgsForCall, struct{}{})
	fake.addressMutex.Unlock()
	if fake.AddressStub != nil {
		return fake.AddressStub()
	} else {
		return fake.addressReturns.result1
	}
}

func (fake *FakeRedisClient) AddressCallCount() int {
	fake.addressMutex.RLock()
	defer fake.addressMutex.RUnlock()
	return len(fake.addressArgsForCall)
}

func (fake *FakeRedisClient) AddressReturns(result1 string) {
	fake.AddressStub = nil
	fake.addressReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeRedisClient) WaitForNewSaveSince(lastSaveTime int64, timeout time.Duration) error {
	fake.waitForNewSaveSinceMutex.Lock()
	fake.waitForNewSaveSinceArgsForCall = append(fake.waitForNewSaveSinceArgsForCall, struct {
		lastSaveTime int64
		timeout      time.Duration
	}{lastSaveTime, timeout})
	fake.waitForNewSaveSinceMutex.Unlock()
	if fake.WaitForNewSaveSinceStub != nil {
		return fake.WaitForNewSaveSinceStub(lastSaveTime, timeout)
	} else {
		return fake.waitForNewSaveSinceReturns.result1
	}
}

func (fake *FakeRedisClient) WaitForNewSaveSinceCallCount() int {
	fake.waitForNewSaveSinceMutex.RLock()
	defer fake.waitForNewSaveSinceMutex.RUnlock()
	return len(fake.waitForNewSaveSinceArgsForCall)
}

func (fake *FakeRedisClient) WaitForNewSaveSinceArgsForCall(i int) (int64, time.Duration) {
	fake.waitForNewSaveSinceMutex.RLock()
	defer fake.waitForNewSaveSinceMutex.RUnlock()
	return fake.waitForNewSaveSinceArgsForCall[i].lastSaveTime, fake.waitForNewSaveSinceArgsForCall[i].timeout
}

func (fake *FakeRedisClient) WaitForNewSaveSinceReturns(result1 error) {
	fake.WaitForNewSaveSinceStub = nil
	fake.waitForNewSaveSinceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRedisClient) RunBGSave() error {
	fake.runBGSaveMutex.Lock()
	fake.runBGSaveArgsForCall = append(fake.runBGSaveArgsForCall, struct{}{})
	fake.runBGSaveMutex.Unlock()
	if fake.RunBGSaveStub != nil {
		return fake.RunBGSaveStub()
	} else {
		return fake.runBGSaveReturns.result1
	}
}

func (fake *FakeRedisClient) RunBGSaveCallCount() int {
	fake.runBGSaveMutex.RLock()
	defer fake.runBGSaveMutex.RUnlock()
	return len(fake.runBGSaveArgsForCall)
}

func (fake *FakeRedisClient) RunBGSaveReturns(result1 error) {
	fake.RunBGSaveStub = nil
	fake.runBGSaveReturns = struct {
		result1 error
	}{result1}
}

var _ client.Client = new(FakeRedisClient)
